#|
	Model of the Go language
|#

;; Lexical elements
;comments
(obj "line comment" :at "text" string)  ; // single-line comment
(obj "general comment" :at "text" string)  ; /* multi-line comment */
;identifiers
(obj "identifier" :union (string))  ; <doc> letter { letter | unicode_digit }


;; Constants
(obj "constant" :union ("bool constant" int real "complex constant" string))
(obj "bool constant" :enum ("true" "false"))
(obj "complex constant" :union("complex number"))
(obj "complex number" :at "re" real :at "im" real)


;; Variables
(obj "variable" :union ("identifier"))


;; Types
(obj "type" :union ("base type" "composite type"))
;base types
(obj "base type" :union ("boolean type" "numeric type" "string type"))
(obj "boolean type" :enum ("bool"))
;numeric types
(obj "numeric type" :union ("real-valued type" "complex type" "byte type" "rune type"))
(obj "real-valued type" :enum ("int type" "float type"))
(obj "int type" :union ("signed int type" "unsigned int type"))
(obj "signed int type" :enum ("int" "int8" "int16" "int32" "int64"))
(obj "unsigned int type" :enum ("uint" "uint8" "uint16" " uint32" "uint64" "uintptr"))
(obj "float type" :enum ("float32" "float64"))
(obj "complex type" :enum ("complex32" "complex64"))
(obj "byte type" :enum ("byte"))  ; <doc> alias for uint8
(obj "rune type" :enum ("rune"))  ; <doc> alias for int32
;string types
(obj "string type" :enum ("string"))
;composite types
(obj "composite type" :union ("array type" "slice type" "struct type" "pointer type" "function type" "interface type" "map type" "channel type"))
(obj "array type" :at "len" nat :at "element type" "type")
(obj "slice type" :at "underlying array" "array type")  ; for decl :at "element type" "type" :at "start index" "nat" :at "len" nat :at "cap" nat
(obj "struct type" :at "fields" (listt "field decl"))
(obj "field decl" :union ("named field" "embedded field"))
(obj "named field" :at "name" "identifier" :at "type" "type")
(obj "embedded field" :at "type" "type")
(obj "pointer type" :at "type" "type")
;function types
(obj "function type" :at "signature" "signature")
(obj "signature" :at "parameters" "parameters" :at "result" "result")
(obj "parameters" :union((listt "parameter decl")))
(obj "paramater decl" :at "identifier" "identifier" :at "type" :at "type")
(obj "result" :union ("parameters" "type"))
;interface types
(obj "interface type" :at "elements" (listt "interface elem"))
(obj "interface eleme" :union ("method elem" "type elem"))
(obj "type elem" :union ((listt "type term")))
(obj "method elem" :at "name" "identifier" :at "signature" "signature")
(obj "type term" :union ("type" "underlying type"))
(obj "underlying type" :at "type" "type")  ; ~Type
;map types
(obj "map type" :at "key type" :at "element type" "type")
;channel types
(obj "channel type" :at "chan" :at "type" "type" :at "direction" "direction")
(obj "direction" :enum ("send" "recieve" "bidirectional"))
(obj "send")  ; chan<- Type
(obj "receive")  ; <-chan Type
(obj "bidirectional")  ; chan Type


;; Blocks
(obj "block" :union ("statement list"))
(obj "statement list" :union ((listt "statement")))


;; Declarations
(obj "declaration" :union ("const decl" "type decl" "var decl"))
(obj "top level decl" :union ("declaration" "function decl" "method decl"))
;const decl
(obj "const decl" :union((listt "const spec")))
(obj "const spec" :at "id list" "id list" :at "type" "type" :at "expr list" "expr list")
(obj "id list" :union ((listt "identifier")))
(obj "expr list" :union ((listt "expression")))
;type decl
(obj "type decl" :union ((listt "type spec")))
(obj "type spec" :union ("alias decl" "type def"))
(obj "alias decl" :at "identifier" "identifier" :at "type parameters" "type parameters" :at "type" "type")
(obj "type parameters" :union ("type param list"))
(obj "type param list" :union ((listt "type param decl")))
(obj "type param decl" :at "id list" "id list" :at "type constraint" "type constraint")
(obj "type constraint" :union ("type elem"))
(obj "type def" :at "identifier" "identifier" :at "type parameters" "type parameters" :at "type" "type")
;var decl
(obj "var decl" :union ((listt "var spec")))
(obj "var spec" :at "id list" "id list" :at "type" "type" :at "expr list" "expr list")
(obj "short var decl" :at "id list" "id list" :at "expr list" "expr list")
;function decl
(obj "function decl" :at "identifier" "identifier" :at "type parameters" "type parameters" :at "signature" "signature" :at "body" "function body")
(obj "function body" :union("block"))
;method decl
(obj "method decl" :at "reciever" "parameters" :at "name" "identifier" :at "signature" "signature" :at "body" "function body")


;; Expressions
;operands
(obj "operand" :union ("literal"))
(obj "literal" :union ("basic literal" "composite literal" "function literal"))
(obj "basic literal" :union("constant"))
(obj "composite literal" :at "type" "literal type" :at "value" "literal value")
(obj "literal type" :union ("struct type" "array type" "*array type" "slice type" "map type" "short arr type"))
(obj "*array type" :at "type" "type")  ; ... ElementType
(obj "short arr type" :at "type name" "type name" :at "type args" (listt "type"))
(obj "type name" :union ("identifier" "qualified identifier"))
(obj "operand name" :union ("identifier" "qualified identifier"))
(obj "qualified identifier" :at "package name" "identifier" :at "identifier" "identifier")  ; <PackageName>.<identifier>
(obj "literal value" :union((listt "keyed element")))
(obj "keyed element" :at "key" "key" :at "element" "element")
(obj "key" :union ("field name" "expression" "literal value"))
(obj "element" :union ("expression" "literal value"))
(obj "function literal" :at "signature" "signature" :at "body" "block")  ; function literals alias for anonymous functions
;primary expressions
(obj "primary expression" :union ("operand" "conversion" "method expr" "primary expr & selector" "primary expr & index" "primary expr & slice" "primary expr & type assertion" "primary expr & argument"))
(obj "conversion" :at "type" "type" :at "expression" "expression")  ; (*Point)(p) - p is converted to *Point
(obj "method expression" :at "receiver type" "type" :at "name" "identifier")
(obj "primary expr & selector" :at "primary expr" "primary expression" :at "selector" "selector")
(obj "selector" :union ("identifier"))  ; x.id
(obj "primary expr & index" :at "primary expr" "primary expression" :at "index" "index")
(obj "index" :union ("expression"))  ; x[e]
(obj "primary expr & slice" :at "primary expr" "primary expression" :at "slice" "slice")
(obj "slice" :at "low" "expression" :at "high" "expression" :at "max" "expression")  ; x[e1 : e2] | x[e1 : e2 : e3]
(obj "primary expr & type assertion" :at "primary expr" "primary expression" :at "type assertion" "type assertion")
(obj "type assertion" :at "type" "type")  ; .(Type)
(obj "primary expr & arguments" :at "primary expr" "primary expression" :at "arguments" "arguments")
(obj "arguments" :union ((listt "expression")))
;expressions
(obj "expression" :union ("unary expression" "e1 & binary operator & e2"))
(obj "e1 & binary operator & e2" :at "arg1" "expression" :at "binary operator" "binary operator" :at "arg2" "expression")
(obj "unary expression" :union ("primary expression" "unary operator & unary expression"))
(obj "unary operator & unary expression" :at "operator" "unary operator" :at "expression" "unary expression")
;operators
;unary operators
(obj "unary operator" :union ("+x" "-x" "^x" "!b" "*p" "&x" "<-"))
(obj "+x" :at "x" "expression") ; +x == 5
(obj "-x" :at "x" "expression") ; -x == -1 * x
(obj "^x" :at "x" "expression") ; ^10(2) == 01(2)
(obj "!b" :at "b" "expression") ; !true == false
(obj "*p" :at "p" "expression")  ; *p
(obj "&x" :at "x" "expression")  ; &x
(obj "<-" :at "channel" "channel type")
;binary operators
(obj "binary operator" :union ("||" "&&" "rel operator" "add operator" "mul operator"))
(obj "||" :at "arg1" "expression" :at "arg2" "expression")  ; true || false == false
(obj "&&" :at "arg1" "expression" :at "arg2" "expression")  ; true && false == false
;binary multiplication operators
(obj "mul operator" :union ("x*y" "x/y" "x%y" "x<<y" "x>>y" "x&y" "x&^y"))
(obj "x*y" :at "x" "expression" :at "y" "expression") ; 5 * 10 == 50
(obj "x/y" :at "x" "expression" :at "y" "expression") ; 4 / 2 = 2.0
(obj "x%y" :at "x" "expression" :at "y" "expression") ; 5 % 2 == 1
(obj "x<<y" :at "x" "expression" :at "y" "expression") ; 1 << 12 == 2 ** 12
(obj "x>>y" :at "x" "expression" :at "y" "expression") ; 32 >> 2 == 2 ** 3
(obj "x&y" :at "x" "expression" :at "y" "expression") ; 1100(2) & 1010(2) == 1000(2)
(obj "x&^y" :at "x" "expression" :at "y" "expression") ; 11001010(2) &^ 10101100(2) == 01000010(2) ; 1 <=> 1 &^ 0
;binary addition operators
(obj "add operator" :union ("x+y" "x-y" "x|y" "x^y"))
(obj "x+y" :at "x" "expression" :at "y" "expression") ; 2 + 3 == 5
(obj "x-y" :at "x" "expression" :at "y" "expression") ; 2 - 3 == -1
(obj "x|y" :at "x" "expression" :at "y" "expression") ; 1100(2) | 1010(2) == 1110(2)
(obj "x^y" :at "x" "expression" :at "y" "expression") ; 1100(2) ^ 1010(2) == 0110(2)
;relation operators
(obj "rel operator" :union ("x==y" "x!=y" "unequality operator"))
(obj "x==y" :at "x" "expression" :at "y" "expression")  ; x == 3 | 4 == 5 | x == y
(obj "x!=y" :at "x" "expression" :at "y" "expression")  ; x != 5
(obj "unequality operator" :union ("x<y" "x<=y" "x>y" "x>=y"))
(obj "x<y" :at "x" "expression" :at "y" "expression")  ; x < 3 | 3 < 5 | x < y
(obj "x<=y" :at "x" "expression" :at "y" "expression")  ; x <= 3 | 3 <= 4 | x <= y
(obj "x>y" :at "x" "expression" :at "y" "expression")  ; x > 3 | 4 > 5 | x > y
(obj "x>=y" :at "x" "expression" :at "y" "expression")  ; x >= 3 | 3 >= 4 | x >= y

;; Statements
(obj "statement" :union ("declaration" "label stmt" "simple stmt" "go stmt" "return stmt" "break stmt" "continue stmt" "goto stmt" "fallthrough stmt" "block" "if stmt" "switch stmt" "select stmt" "for stmt" "defer stmt"))
(obj "simple statement" :union ("empty stmt" "expression stmt" "send stmt" "x++ stmt | x-- stmt" "assignment" "short var decl"))
(obj "empty stmt")  ; <doc> The empty statement does nothing
(obj "label stmt" :at "label" "label" :at "statement" "statement")  ; <example> Error: log.Panic("error encountered") ???
(obj "label" :union ("identifier"))
(obj "expression stmt" :union ("expression"))
(obj "send stmt" :at "channel" "expression" :at "message" "expression")
(obj "x++ stmt | x-- stmt" :union ("x++ stmt" "x-- stmt"))
(obj "x++ stmt" :at "x" "expression")
(obj "x-- stmt" :at "x" "expression")
(obj "assignment" :at "location" (listt "expression") :at "operator" "assignment operators" :at "expression" (listt "expression"))
(obj "assignment operators" :union ("x=y" "x+=y" "x-=y" "x|=y" "x^=y" "x*=y" "x/=y" "x%=y" "x<<=y" "x>>=y" "x&=y" "x&^=y"))
(obj "x=y" :at "x" "expression" :at "y" "expression")
(obj "x+=y" :at "x" "expression" :at "y" "expression")
(obj "x-=y" :at "x" "expression" :at "y" "expression")
(obj "x|=y" :at "x" "expression" :at "y" "expression")
(obj "x^=y" :at "x" "expression" :at "y" "expression")
(obj "x*=y" :at "x" "expression" :at "y" "expression")
(obj "x/=y" :at "x" "expression" :at "y" "expression")
(obj "x%=y" :at "x" "expression" :at "y" "expression")
(obj "x<<=y" :at "x" "expression" :at "y" "expression")
(obj "x>>=y" :at "x" "expression" :at "y" "expression")
(obj "x&=y" :at "x" "expression" :at "y" "expression")
(obj "x&^=y" :at "x" "expression" :at "y" "expression")
(obj "if stmt" :at "init" "simple stmt" :at "expression" "expression" :at "block" "block" :at "else" "if stmt | block")
(obj "if stmt | block" :union ("if statement" "block"))
(obj "switch stmt" :union ("expr switch stmt" "type switch stmt"))
(obj "expr switch stmt" :at "stmt" "simple stmt" :at "expression" "expression" :at "expr case clause")
(obj "expr case clause" :at "case" "expr switch case" :at "stmt" (listt "statement"))
(obj "expr switch case" :union ((listt "expression")))
(obj "type switch stmt" :at "stmt" "simple stmt" :at "type switch guard" "type switch guard" :at "type case clause" "type case clause")
(obj "type switch guard" :at "identifier" "identifier" :at "expression" "primary expression")
(obj "type case clause" :at "case" "type switch case" :at "stmt" (listt "statement"))
(obj "type switch case" :union ((listt "type")))
(obj "for stmt" :at "variety" "condition | for clause | range clause" :at "block" "block")
(obj "condition | for clause | range clause" :union ("condition" "for clause" "range clause"))
(obj "condition" :union(("expression")))
(obj "for clause" :at "init" "simple stmt" :at "condition" "condition" :at "post" "simple stmt")
(obj "range clause" :at "location" "expr list | id list" :at "expression" "expression")
(obj "expr list | id list" :union ("expr list" "id list"))
(obj "go stmt" :union (("expression")))
(obj "select stmt" :union("common slause"))
(obj "common clause" :at "case" "common case" :at "stmt list" "stmt list")
(obj "common case" :union ("send stmt" "recv stmt"))
(obj "recv stmt" :at "list" "expr list | id list" :at "expression" "expression")
(obj "expr list | id list" :union ("expr list" "id list"))
(obj "stmt list" :union ((listt "statement")))
(obj "return stmt" :at "expr list" "expr list")
(obj "break stmt" :at "label" "label")
(obj "continue stmt" :at "label" "label")
(obj "goto stmt" :at "label" "label")
(obj "fallthrough stmt")  ; A "fallthrough" statement transfers control to the next case clause in a switch
(obj "defer stmt" :at "expression" "expression")


;; Source file organization
(obj "source file" :at "package clause" "package clause" :at "import decl" "import decl" :at "top level decl" "top level decl")
;package clause
(obj "package clause" :union ("package name"))  ; package math
(obj "package name" :union ("identifier"))
;import declaration
(obj "import decl" :at "import spec list" (listt "import spec"))
(obj "import spec" :at "name" "package name" :at "import path" "import path")
(obj "import path" :union(string))  ; import m "lib/math"


;; For the future
#|
(obj "variable declaration" :at "name" "identifier" :at "type" "type" :at "value" "expression")  ; Need to check: (is-instance (aget variable "value") (aget variable "type"))
(obj "array" :at "type" "array type" :at "value" (listt any))
(obj "slice" :at "type" "slice type" :at "start index" nat :at "length" nat :at "capacity" nat)
|#
