/*
 * Variant 4.1.1 - DL
 * version 1.0
 *
 * --Additional task 1
 * --Additional task 2
 * --Additional task 3
 */

#include <iostream>
#include <string>
#include <utility>
#include <unordered_map>

using namespace std;

#define error() throw invalid_argument("Incorrect input data")

class Expression {
public:
    virtual Expression *eval() = 0;

    virtual int getValue() {
        throw invalid_argument("Invalid type of expression for 'getValue'");
    }

    virtual void checkFunction() {
        throw invalid_argument("It is not function");
    }

    [[nodiscard]] virtual bool isVar() const {
        return false;
    }

    virtual string getId() {
        throw invalid_argument("Invalid class for 'getId'");
    }

    virtual string toString() {
        return "<expression>";
    }

    [[nodiscard]] virtual bool isFunction() const {
        return false;
    }
};

unordered_map<string, Expression *> env;

Expression *fromEnv(const string &id) {
    if (env.find(id) == env.end()) throw invalid_argument("Invalid id for 'fromEnv'");
    return env[id];
}

int getValue(Expression *exp) {
    return exp->getValue();
}

unordered_map<string, Expression *> existing_fs;

class Val : public Expression {
    int n;

public:
    explicit Val(int n) : n(n) {}

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

    Expression *eval() override {
        return fromEnv(id);
    }

    [[nodiscard]] bool isVar() const override {
        return true;
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
        if (e_value->isFunction()) existing_fs[id] = e_value;
    }

    Expression *eval() override {
        env[id] = e_value;
        return e_body->eval();
    }
};

class Function : public Expression {
    string id;
    Expression *e;

public:
    explicit Function(string &id, Expression *e) : id(id), e(e) {}

    Expression *eval() override {
        return this;
    }

    void checkFunction() override {}

    string getId() override {
        return id;
    }

    [[nodiscard]] Expression *getE() const {
        return e;
    }

    [[nodiscard]] bool isFunction() const override {
        return true;
    }
};

class Call : public Expression {
    Expression *f_e;
    Expression *arg_e;

public:
    explicit Call(Expression *f_e, Expression *arg_e) : f_e(f_e), arg_e(arg_e) {}

    Expression *eval() override {
        if (f_e->isVar()) {
            string id = f_e->getId();
            if (existing_fs.find(id) == existing_fs.end()) throw invalid_argument("There is not a function");
            Expression *result = (new Call(existing_fs[id], arg_e))->eval();
            return result;
        } else {
            f_e->checkFunction();
            Expression *t = ((Function *) f_e)->getE();
            string id = f_e->getId();
            Expression *old = nullptr;
            if (env.find(id) != env.end())
                old = env[id];
            env[id] = arg_e->eval();
            Expression *result = t->eval();
            env[id] = old;
            return result;
        }
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

Expression *extract() {
    char c;
    do cin >> c; while (c == ' ' || c == '\n' || c == '\t');
    if (c != '(') error();
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
    } else {
        error();
    }
}

int main() {
    try {
        Expression *exp = extract();
        cout << exp->eval()->toString() << endl;
    } catch (exception &e) {
        cout << "ERROR" << endl;
    }
    return 0;
}
