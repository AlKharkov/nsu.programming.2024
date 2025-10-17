#|
	Model of the Go language

	Last edit: 17/10/2025
|#

;;; Lexical elements
(obj "line comment" :at "text" string)     ; // single-line comment
(obj "general comment" :at "text" string)  ; /* multi-line comment */
(obj "identifier" :union (string))         ; a | _x9 | ThisVariableIsExported | αβ


;;; Constants
(obj "constant" :union ("bool constant" int real "complex constant" string))
(obj "bool constant" :enum ("true" "false"))
(obj "complex constant" :at "re" real :at "im" real)


;;; Variables
(obj "variable" :union ("identifier"))


;;; Types
(obj "type" :union ("base type" "composite type"))
(obj "base type" :union ("boolean type" "numeric type" "string type"))
(obj "boolean type" :enum ("bool"))
(obj "numeric type" :union ("real-valued type" "complex type" "byte type" "rune type"))
(obj "real-valued type" :enum ("int type" "float type"))
(obj "int type" :union ("signed int type" "unsigned int type"))
(obj "signed int type" :enum ("int" "int8" "int16" "int32" "int64"))
(obj "unsigned int type" :enum ("uint" "uint8" "uint16" " uint32" "uint64" "uintptr"))
(obj "float type" :enum ("float32" "float64"))
(obj "complex type" :enum ("complex32" "complex64"))
(obj "byte type" :enum ("byte"))  ; <doc> alias for uint8
(obj "rune type" :enum ("rune"))  ; 'a' | 'ä' | '本' | '\t' | '\377' | '\U00101234'
(obj "string type" :enum ("string"))
;;composite types
(obj "composite type" :union ("array type" "slice type" "struct type" "pointer type" "function type" "interface type" "map type" "channel type"))
(obj "array type" :at "len" nat :at "element type" "type")  ; [32]byte | [2][2][2]float64 ~ same as [2]([2]([2]float64))
(obj "slice type" :at "underlying array" "array type")
; a := [5]int{1, 2, 3, 4, 5}    // a = [1, 2, 3, 4, 5] - underlying array
; s := a[1:4]                   // s = [2, 3, 4] - slice
(obj "struct type" :at "fields" (listt "field decl"))
(obj "field decl" :union ("named field" "embedded field"))
(obj "named field" :at "name" "identifier" :at "type" "type")
(obj "embedded field" :at "type" "type")
; struct {
; 	x, y int    // named field 
; 	A *[]int
;   T           // embedded field
; }
(obj "pointer type" :at "type" "type")
;;function types
(obj "function type" :at "signature" "signature")
(obj "signature" :at "parameters" "parameters" :at "result" "function result")
(obj "parameters" :union((listt "parameter decl")))
(obj "paramater decl" :at "name" "identifier" :at "type" "type | variadic type")
(obj "type | variadic type" :union ("type" "variadic type"))
(obj "variadic type" :at "type" "type")  ; func f(numbers ...int) -> f(1, 2, 3, 4)
(obj "function result" :union ("parameters" "type"))
; func("parameters") "result"
; func(a, b int, z float32) bool
; func(int, int, float64) (bool, int) ~ same as func(_, _ int, _ float64) (bool, int)
; func(n int) func(p *T)
;;interface types
(obj "interface type" :at "elements" (listt "interface elem"))
(obj "interface elem" :union ("method elem" "type elem"))
(obj "type elem" :union ((listt "type term")))
(obj "method elem" :at "name" "identifier" :at "signature" "signature")
(obj "type term" :union ("type" "underlying type"))
; type A interface {
;   f(n int) (q bool)     // method elem: "name" = `f`, "signature" = `(n int) (q bool)`
; }
; type B interface {
;   A                     // type elem: includes methods of A in B's method set
;   g(p *int) (a int[]) 
; }
; interface {             // An interface representing all types with underlying type int that implement the String method.
; 	~int
; 	String() string
; }
(obj "underlying type" :at "type" "type")  ; ~Type
; type (    
;   // the underlying type of string, A1 and A2 is string
; 	A1 = string
; 	A2 = A1
;
;   // the underlying type of B1 and B2 is string; the underlying type of []B1, B3 and B4 is []B1
; 	B1 string
; 	B2 B1
; 	B3 []B1
; 	B4 B3
; )
(obj "map type" :at "key type" "type" :at "element type" "type")  ; map[*T]struct{ x, y float64 }
(obj "channel type" :at "chan" :at "type" "type" :at "direction" "direction")
(obj "direction" :enum ("send" "receive" "bidirectional"))
; chan T            // can be used to send and receive values of type T
; chan<- float64    // can only be used to send float64s
; <-chan int        // can only be used to receive ints


;;; Blocks
(obj "block" :at "statements" (listt "statement"))  ; { "statement" ... "statement" }


;;; Declarations
(obj "declaration" :union ("const decl" "type decl" "var decl"))
(obj "top level decl" :union ("declaration" "function decl" "method decl"))
;;const declarations
(obj "const decl" :at "specs" (listt "const spec"))
(obj "const spec" :at "names" "identifier list" :at "type" "type" :at "initializers" "expr list")
; const (
; 	size int64 = 1024
; 	eof        = -1       // untyped integer constant
; )
(obj "identifier list" :union ((listt "identifier")))
(obj "expr list" :union ((listt "expression")))
;;type declarations
(obj "type decl" :union ((listt "alias decl | type def")))
(obj "alias decl | type def" :union ("alias decl" "type def"))
(obj "alias decl" :at "name" "identifier" :at "type parameters" "type parameters" :at "type" "type")
; type (
; 	nodeList = []*Node    // nodeList and []*Node are identical types
; 	Polar    = polar      // Polar and polar denote identical types
; )
; type set[P comparable] = map[P]bool
(obj "type parameters" :at "param declarations" (listt "type param decl"))
(obj "type param decl" :at "identifier list" "identifier list" :at "type constraint" "type constraint")
(obj "type constraint" :at "constraints" "type elem")  ; [T1 comparable, T2 any]
(obj "type def" :at "identifier" "identifier" :at "type parameters" "type parameters" :at "type" "type")
; type (
; 	Point struct{ x, y float64 }    // Point and struct{ x, y float64 } are different types
; 	polar Point                     // polar and Point denote different types
; )
;;var declarations
(obj "var decl" :at "specifiers" (listt "var spec"))
(obj "var spec" :at "identifier list" "identifier list" :at "type" "type" :at "expr list" "expr list")
; var U, V, W float64
; var k = 0
; var x, y float32 = -1, -2
(obj "short var decl" :at "identifier list" "identifier list" :at "expr list" "expr list")
; f := func() int { return 7 }
;;function declarations
(obj "function decl" :at "identifier" "identifier" :at "type parameters" "type parameters" :at "signature" "signature" :at "body" "function body")
(obj "function body" :at "block" "block")
; func min[T ~int|~float64](x, y T) T {
; 	if x < y {
; 		return x
; 	}
; 	return y
; }
;;method declarations
(obj "method decl" :at "reciever" "parameters" :at "name" "identifier" :at "signature" "signature" :at "body" "function body")
; func (p *Point) Scale(factor float64) {
; 	p.x *= factor
; 	p.y *= factor
; }


;;; Expressions
;;operands
(obj "operand" :union ("literal" "parenthesized expression"))
(obj "literal" :union ("basic literal" "composite literal" "function literal"))
(obj "basic literal" :union("constant"))
(obj "composite literal" :at "type" "literal type" :at "value" (listt "keyed element"))
(obj "literal type" :union ("struct type" "array type" "*array type" "slice type" "map type" "short arr type"))
(obj "... type" :at "type" "type")  ; ... ElementType
(obj "short arr type" :at "name" "identifier | qualified identifier" :at "type args" (listt "type"))
(obj "identifier | qualified identifier" :union ("identifier" "qualified identifier"))
(obj "qualified identifier" :at "package name" "identifier" :at "identifier" "identifier")  ; math.Sin
(obj "keyed element" :at "key" "key" :at "element" "element")
(obj "key" :union ("field name" "expression" (listt "keyed element")))
(obj "element" :union ("expression" (listt "keyed element")))
(obj "function literal" :at "signature" "signature" :at "body" "block")  ; func(a, b int, z float64) bool { return a*b < int(z) }
(obj "parenthesized expression" :at "expression" "expression")
; import (
;   "fmt"
;   "math"
; )
;;primary expressions
(obj "primary expression" :union ("operand" "conversion" "method expr" "primary expr & selector" "primary expr & index expr" "primary expr & slice expr" "primary expr & type assertion" "primary expr & argument"))
(obj "conversion" :at "type" "type" :at "expression" "expression")  ; (*Point)(p)    // p is converted to *Point
(obj "method expression" :at "receiver type" "type" :at "name" "identifier")
; type T struct {
; 	a int
; }
; func (tv T) Mv(b int) int { return tv.a + b }    // value receiver
(obj "primary expr & selector" :at "primary expr" "primary expression" :at "selector" "selector")
(obj "selector" :at "name" "identifier")  ; x.id
(obj "primary expr & index expr" :at "primary expr" "primary expression" :at "index" "index")
(obj "index expr" :at "expression" "expression")  ; x[e]
(obj "primary expr & slice expr" :at "primary expr" "primary expression" :at "slice" "slice")
(obj "slice expr" :at "low" "expression" :at "high" "expression" :at "max" "expression")  ; x[e1 : e2] | x[e1 : e2 : e3]
(obj "primary expr & type assertion" :at "primary expr" "primary expression" :at "type assertion" "type assertion")
(obj "type assertion" :at "type" "type")  ; .(Type)
; var x interface{} = 7    // x has dynamic type int and value 7
; i := x.(int)             // i has type int and value 7
(obj "primary expr & arguments" :at "primary expr" "primary expression" :at "arguments" (listt "expression"))  ; math.Sin(2)
;;expressions
(obj "expression" :union ("unary expression" "1 & binary operator & 2"))
(obj "1 & binary operator & 2" :at 1 "expression" :at "binary operator" "binary operator" :at 2 "expression")
(obj "unary expression" :union ("primary expression" "unary operator & unary expression"))
(obj "unary operator & unary expression" :at "operator" "unary operator" :at "expression" "unary expression")
;;operators
;;unary operators
(obj "unary operator" :union ("+1" "-1" "^1" "!1" "*1" "&1" "<-"))
(obj "+1" :at 1 "expression")  ; +x == x
(obj "-1" :at 1 "expression")  ; -x == -1 * x
(obj "^1" :at 1 "expression")  ; ^10(2) == 01(2)
(obj "!1" :at 1 "expression")  ; !true == false
(obj "*1" :at 1 "expression")
(obj "&1" :at 1 "expression")
(obj "<-" :at "channel" "channel type")
;;binary operators
(obj "binary operator" :union ("1||2" "1&&2" "rel operator" "add operator" "mul operator"))
(obj "1||2" :at 1 "expression" :at 2 "expression")
(obj "1&&2" :at 1 "expression" :at 2 "expression")
;;binary multiplication operators
(obj "mul operator" :union ("1*2" "1/2" "1%2" "1<<2" "1>>2" "1&2" "1&^2"))
(obj "1*2" :at 1 "expression" :at 2 "expression")
(obj "1/2" :at 1 "expression" :at 2 "expression")
(obj "1%2" :at 1 "expression" :at 2 "expression")
(obj "1<<2" :at 1 "expression" :at 2 "expression")
(obj "1>>2" :at 1 "expression" :at 2 "expression")
(obj "1&2" :at 1 "expression" :at 2 "expression")
(obj "1&^2" :at 1 "expression" :at 2 "expression") ; 11001010(2) &^ 10101100(2) == 01000010(2)    // 1 <=> 1 &^ 0
;;binary addition operators
(obj "add operator" :union ("1+2" "1-2" "1|2" "1^2"))
(obj "1+2" :at 1 "expression" :at 2 "expression")
(obj "1-2" :at 1 "expression" :at 2 "expression")
(obj "1|2" :at 1 "expression" :at 2 "expression")
(obj "1^2" :at 1 "expression" :at 2 "expression")
;;relation operators
(obj "rel operator" :union ("1=2" "1!2" "unequality operator"))
(obj "1==2" :at 1 "expression" :at 2 "expression")
(obj "1!=2" :at 1 "expression" :at 2 "expression")
(obj "unequality operator" :union ("1<2" "1<2" "1>2" "1>2"))
(obj "1<2" :at 1 "expression" :at 2 "expression")
(obj "1<=2" :at 1 "expression" :at 2 "expression")
(obj "1>2" :at 1 "expression" :at 2 "expression")
(obj "1>=2" :at 1 "expression" :at 2 "expression")

;;; Statements
(obj "statement" :union ("declaration" "label stmt" "simple stmt" "go stmt" "return stmt" "break stmt" "continue stmt" "goto stmt" "fallthrough stmt" "block" "if stmt" "switch stmt" "select stmt" "for stmt" "defer stmt"))
(obj "simple statement" :union ("empty stmt" "expression stmt" "send stmt" "1++" "1--" "assignment stmt" "short var decl"))
(obj "empty stmt")  ; <doc> The empty statement does nothing
(obj "label stmt" :at "label" "label" :at "statement" "statement")
(obj "label" :at "name" "identifier")
; "label": "statement"
;   Error: log.Panic("error encountered")
(obj "expression stmt" :at "expression" "expression")
; h(x+y)
; f.Close()
(obj "send stmt" :at "channel" "expression" :at "message" "expression")
(obj "1++ stmt" :at "1" "expression")  ; x += 1
(obj "1-- stmt" :at "1" "expression")  ; x -= 1
;;assignment statements
(obj "assignment stmt" :union ("1=2" "1+=2" "1-=2" "1|=2" "1^=2" "1*=2" "1/=2" "1%=2" "1<<=2" "1>>=2" "1&=2" "1&^=2"))
(obj "1=2" :at 1 (listt "expression") :at 2 (listt "expression"))
(obj "1+=2" :at 1 (listt "expression") :at 2 (listt "expression"))
(obj "1-=2" :at 1 (listt "expression") :at 2 (listt "expression"))
(obj "1|=2" :at 1 (listt "expression") :at 2 (listt "expression"))
(obj "1^=2" :at 1 (listt "expression") :at 2 (listt "expression"))
(obj "1*=2" :at 1 (listt "expression") :at 2 (listt "expression"))
(obj "1/=2" :at 1 (listt "expression") :at 2 (listt "expression"))
(obj "1%=2" :at 1 (listt "expression") :at 2 (listt "expression"))
(obj "1<<=2" :at 1 (listt "expression") :at 2 (listt "expression"))
(obj "1>>=2" :at 1 (listt "expression") :at 2 (listt "expression"))
(obj "1&=2" :at 1 (listt "expression") :at 2 (listt "expression"))
(obj "1&^=2" :at 1 (listt "expression") :at 2 (listt "expression"))
;;if statement
(obj "if stmt" :at "init" "simple stmt" :at "condition" "expression" :at "then" "block" :at "else" "if stmt | statements")
(obj "if stmt | statements" :union ("if statement" (listt "statement")))
; if x := f(); x < y {
; 	return x
; } else if x > z {
; 	return z
; }
(obj "switch stmt" :union ("expr switch stmt" "type switch stmt"))
(obj "expr switch stmt" :at "init" "simple stmt" :at "controlling expression" "expression" :at "cases" (listt "expr case clause"))
(obj "expr case clause" :at "cases" (listt "expression") :at "statements" (listt "statement"))
; switch tag {
; case 0, 1, 2, 3: s1()
; case 4, 5, 6, 7: s2()
; default: s3()
; }
(obj "type switch stmt" :at "init" "simple stmt" :at "guard" "type switch guard" :at "cases" (listt "type case clause"))
(obj "type switch guard" :at "type variable" "identifier" :at "variable" "primary expression")
(obj "type case clause" :at "cases" (listt "type") :at "statements" (listt "statement"))
; switch t := x.(type) {              // type checker
; case nil:
; 	fmt.Println("x is nil")  
; default:
;   fmt.Println("don't know the type")
(obj "for statement" :union ("for stmt with condition clause" "for stmt with range clause"))
(obj "for stmt with condition clause" :at "init" "simple stmt" :at "condition" "expression" :at "post" "simple stmt" :at "body" "block")
(obj "for stmt with range clause" :at "location" "expr list | identifier list" :at "expression" "expression" :at "body" "block")
(obj "expr list | identifier list" :union ("expr list" "identifier list"))
; for x < y { a *= 2 }                // the usual while loop
; for i := 0; i < 10; i++ { f(i) }    // the usual for loop
; a := [5]int{1, 2, 3, 4, 5}
; for i, v := range a {
;   fmt.Println(i, v)                 // i a[i]
; }
(obj "go stmt" :at "thread body" "expression")  ; <doc> It starts the execution of a function call as an independent concurrent thread of control
; go Server()
; go func(ch chan<- bool) { for { sleep(10); ch <- true }} (c)
(obj "select stmt" :at "statement" (listt "common clause"))
(obj "common clause" :at "case" "common case" :at "stmt list" (listt "statement"))
(obj "common case" :union ("send stmt" "recive stmt"))
(obj "recive stmt" :at "list" "expr list | identifier list" :at "expression" "expression")
; func f1(c1 chan int) {
; 	a := <-c1
; 	c1 <- a
; }
; func f2(c2 chan int) {
; 	b := <-c2
; 	c2 <- b
; }
; c1, c2 := make(chan int), make(chan int)
; go f1(c1)
; go f2(c2)
; c1 <- 3
; c2 <- 4
; select {
; case a := <-c1:
;   fmt.Println(1, a)
; case b := <-c2:
;   fmt.Println(2, b)
; }
(obj "return" :at "expressions" "expr list")
(obj "break " :at "label" "label")
; OuterLoop:
;   for {
;     for {
;       break OuterLoop
;     }
;   }
(obj "continue" :at "label" "label")
(obj "goto" :at "label" "label")
(obj "fallthrough" :enum ("fallthrough"))
; switch {                     // prints "true & false"
;   case true:
;     fmt.Print("true & ")
;     fallthrough              // It transfers control to the next case clause in a switch
;   case false:
;     fmt.Print("false")
; }
(obj "defer stmt" :at "expression" "expression")  ; A "defer" statement invokes a function whose execution is deferred to the moment the surrounding function returns
; for i := 0; i <= 3; i++ {    // prints 3 2 1 0
; 	defer fmt.Print(i)
; }


;;; Source file organization
(obj "source file" :at "package" "identifier" :at "import decl" (listt "import package") :at "declarations" "top level decl")
;;import declaration
(obj "import package" :at "name" "package name" :at "path" "path")
(obj "path" :union(string))  ; import m "lib/math"
