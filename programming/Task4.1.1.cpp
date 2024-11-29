/*
 * Variant 4.1.1 - DL
 * version 1.4
 *
 * --Additional task 1
 * Additional task 2
 * Additional task 3
 */

#include <iostream>
#include <string>
#include <utility>
#include <unordered_map>
#include <vector>

using namespace std;

#define error() throw invalid_argument("Incorrect input data")

enum Types {
    expression, val, var, add, if_t, let, function, call, set, block, arr, gen, at
};

class Expression {
public:
    virtual Types getType() {
        return expression;
    }

    virtual Expression *eval() = 0;

    virtual int getValue() {
        throw invalid_argument("Invalid type of expression for 'getValue'");
    }

    virtual string getId() {
        throw invalid_argument("Invalid class for 'getId'");
    }

    virtual string toString() {
        return "<expression>";
    }
};

unordered_map<string, Expression *> env;

Expression *fromEnv(const string &id) {
    if (env.find(id) == env.end()) throw invalid_argument("Invalid id for 'fromEnv'");
    return env[id];
}

class Val : public Expression {
    int n;

public:
    explicit Val(int n) : n(n) {}

    Types getType() override {
        return val;
    }

    int getValue() override {
        return n;
    }

    Expression *eval() override {
        return this;
    }

    string toString() override {
        return "(val " + to_string(n) + ")";
    }
};

class Var : public Expression {
    string id;

public:
    explicit Var(string s) : id(std::move(s)) {}

    Types getType() override {
        return var;
    }

    Expression *eval() override {
        return fromEnv(id);
    }

    string getId() override {
        return id;
    }

    string toString() override {
        return "(var " + id + ")";
    }
};

class Add : public Expression {
    Expression *e1;
    Expression *e2;

public:
    explicit Add(Expression *a, Expression *b) : e1(a), e2(b) {}

    Types getType() override {
        return add;
    }

    Expression *eval() override {
        return new Val(e1->eval()->getValue() + e2->eval()->getValue());
    }
};

class If : public Expression {
    Expression *e1;
    Expression *e2;
    Expression *e_then;
    Expression *e_else;

public:
    explicit If(Expression *e1, Expression *e2, Expression *e_then, Expression *e_else) :
            e1(e1), e2(e2), e_then(e_then), e_else(e_else) {}

    Types getType() override {
        return if_t;
    }

    Expression *eval() override {
        if (e1->eval()->getValue() > e2->eval()->getValue()) return e_then->eval();
        return e_else->eval();
    }
};

class Let : public Expression {
    string id;
    Expression *e_value;
    Expression *e_body;

public:
    explicit Let(const string &id, Expression *e_value, Expression *e_body) : id(id), e_value(e_value),
                                                                              e_body(e_body) {
        if (e_value->getType() == function) env[id] = e_value;
    }

    Types getType() override {
        return let;
    }

    Expression *eval() override {
        env[id] = e_value->eval();
        return e_body->eval();
    }
};

class Function : public Expression {
    string id;
    Expression *e;

public:
    explicit Function(string id, Expression *e) : id(std::move(id)), e(e) {}

    Types getType() override {
        return function;
    }

    Expression *eval() override {
        return this;
    }

    string getId() override {
        return id;
    }

    [[nodiscard]] Expression *getE() const {
        return e;
    }

    string toString() override {
        return "(function " + id + " " + e->toString() + ")";
    }
};

class Call : public Expression {
    Expression *f_e;
    Expression *arg_e;

public:
    explicit Call(Expression *f_e, Expression *arg_e) : f_e(f_e), arg_e(arg_e) {}

    Types getType() override {
        return call;
    }

    Expression *eval() override {
        if (f_e->getType() != var) {
            if (f_e->eval()->getType() != function) error();
        }
        Expression *t = ((Function *) f_e->eval())->getE();
        string id = f_e->eval()->getId();
        Expression *old = nullptr;
        if (env.find(id) != env.end())
            old = env[id];
        env[id] = arg_e->eval();
        Expression *result = t->eval();
        env[id] = old;
        return result;
    }
};

class Set : public Expression {  // Additional task 2
    string id;
    Expression *e;

public:
    explicit Set(string id, Expression *e) : id(std::move(id)), e(e) {}

    Types getType() override {
        return set;
    }

    Expression *eval() override {
        env[id] = e->eval();
        return this;
    }
};

class Block : public Expression {  // Additional task 2
    vector<Expression *> e;

public:
    explicit Block(vector<Expression *> e) : e(std::move(e)) {}

    Types getType() override {
        return block;
    }

    Expression *eval() override {
        for (int i = 0; i + 1 < e.size(); ++i) e[i]->eval();
        return e.back()->eval();
    }
};

class Arr : public Expression {  // Additional task 3
    vector<Expression *> e;

public:
    explicit Arr(vector<Expression *> e) : e(std::move(e)) {}

    Types getType() override {
        return arr;
    }

    Expression *eval() override {
        vector<Expression *> e2;
        for (Expression *x: e) e2.emplace_back(x->eval());
        return new Arr(e2);
    }

    Expression *getElem(int ind) {
        if (0 <= ind && ind < e.size()) return e[ind];
        else
            error();
    }
};

class Gen : public Expression {  // Additional task 3
    Expression *e1;
    Expression *e2;

public:
    explicit Gen(Expression *e1, Expression *e2) : e1(e1), e2(e2) {}

    Types getType() override {
        return gen;
    }

    Expression *eval() override {
        Expression *e_length = e1->eval();
        if (e_length->getType() != val) error();
        int n = e_length->getValue();
        vector<Expression *> e;
        for (int i = 0; i < n; ++i) {
            e.emplace_back(new Call(e2, new Val(i)));
        }
        return new Arr(e);
    }
};

class At : public Expression {  // Additional task 3
    Expression *e1;
    Expression *e2;

public:
    explicit At(Expression *e1, Expression *e2) : e1(e1), e2(e2) {}

    Types getType() override {
        return at;
    }

    Expression *eval() override {
        Expression *e_array = e1->eval();
        if (e_array->getType() != arr) error();
        int ind = e2->eval()->getValue();
        return ((Arr *) e_array)->getElem(ind);
    }
};

string input_before_close_bracket() {
    string s;
    char c;
    do cin >> c; while (c == ' ' || c == '\n' || c == '\t');
    while (c != ')') {
        s += c;
        cin >> c;
        if (c == ' ' || c == '\t' || c == '\n') {
            do cin >> c; while (c == ' ' || c == '\n' || c == '\t');
            if (c != ')') error();
        }
    }
    return s + c;
}

void isId(const string &id) {
    if (id.empty() || ('0' <= id[0] && id[0] <= '9')) error();
}

Expression *extract(bool need_check_first = true) {
    char c;
    if (need_check_first) {  // Non - for Additional task 2
        do cin >> c; while (c == ' ' || c == '\n' || c == '\t');
        if (c != '(') error();
    }
    string s;
    cin >> s;
    if (s == "val") {
        s = input_before_close_bracket();
        if (*s.rbegin() != ')') {
            do cin >> c; while (c == ' ' || c == '\n' || c == '\t');
            if (c != ')') error();
        }
        return new Val(stoi(s.substr(0, s.size() - 1)));
    } else if (s == "var") {
        s = input_before_close_bracket();
        if (*s.rbegin() != ')') {
            do cin >> c; while (c == ' ' || c == '\n' || c == '\t');
            if (c != ')') error();
        }
        return new Var(s.substr(0, s.size() - 1));
    } else if (s == "add") {
        Expression *e1 = extract();
        Expression *e2 = extract();
        do cin >> c; while (c == ' ' || c == '\n' || c == '\t');
        if (c != ')') error();
        return new Add(e1, e2);
    } else if (s == "if") {
        Expression *e1 = extract();
        Expression *e2 = extract();
        cin >> s;
        if (s != "then") error();
        Expression *e_then = extract();
        cin >> s;
        if (s != "else") error();
        Expression *e_else = extract();
        do cin >> c; while (c == ' ' || c == '\n' || c == '\t');
        if (c != ')') error();
        return new If(e1, e2, e_then, e_else);
    } else if (s == "let") {
        string id;
        cin >> id;
        isId(id);
        cin >> s;
        if (s != "=") error();
        Expression *e1 = extract();
        cin >> s;
        if (s != "in") error();
        Expression *e2 = extract();
        do cin >> c; while (c == ' ' || c == '\n' || c == '\t');
        if (c != ')') error();
        return new Let(id, e1, e2);
    } else if (s == "function") {
        string id;
        cin >> id;
        isId(id);
        Expression *e = extract();
        do cin >> c; while (c == ' ' || c == '\n' || c == '\t');
        if (c != ')') error();
        return new Function(id, e);
    } else if (s == "call") {
        Expression *e1 = extract();
        Expression *e2 = extract();
        do cin >> c; while (c == ' ' || c == '\n' || c == '\t');
        if (c != ')') error();
        return new Call(e1, e2);
    } else if (s == "set") {  // Additional task 2
        string id;
        cin >> id;
        isId(id);
        Expression *e = extract();
        do cin >> c; while (c == ' ' || c == '\n' || c == '\t');
        if (c != ')') error();
        return new Set(id, e);
    } else if (s == "block") {  // Additional task 2
        vector<Expression *> e;
        do cin >> c; while (c == ' ' || c == '\n' || c == '\t');
        if (c == ')') return new Block(e);
        else if (c == '(') e.emplace_back(extract(false));
        else
            error();
    } else if (s == "arr") {  // Additional task 3
        vector<Expression *> e;
        do cin >> c; while (c == ' ' || c == '\n' || c == '\t');
        if (c == ')') return new Arr(e);
        else if (c == '(') e.emplace_back(extract(false));
        else
            error();
    } else if (s == "gen") {  // Additional task 3
        Expression *e1 = extract();
        Expression *e2 = extract();
        do cin >> c; while (c == ' ' || c == '\n' || c == '\t');
        if (c != ')') error();
        return new Gen(e1, e2);
    } else if (s == "at") {  // Additional task 3
        Expression *e1 = extract();
        Expression *e2 = extract();
        do cin >> c; while (c == ' ' || c == '\n' || c == '\t');
        if (c != ')') error();
        return new At(e1, e2);
    } else {
        error();
    }
}

int main() {
    try {
        Expression *exp = extract();
        cout << exp->eval()->toString() << endl;
    }
    catch (exception &e) {
        cout << "ERROR" << endl;
    }
    return 0;
}
