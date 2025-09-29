#|
	Model of the Go language
|#


;; Comments
(obj "line comment" :at "text" string)  ; // single-line comment
(obj "general comment" :at "text" string)  ; /* multi-line comment */


;; Identifiers
(obj "identifier" :union (string))  ; <doc> letter { letter | unicode_digit }
(obj "qualified identifier type" :enum ("identifier" "qualified identifier"))
(obj "qualified identifier" :at "package" "qualified identified type" :at "name" "identifier")  ; <PackageName>.<identifier>


;; Constants
(obj "bool constant" :enum ("true" "false"))
(obj "complex constant" :at "re" real :at "im" real)
(obj "constant" :union ("bool constant" int real string))  ; ??? Go allows the use of complex constants ???


;; Variables
(obj "variable" :union ("identifier"))


;; Types
(obj "boolean type" :enum ("bool"))
(obj "signed int type" :enum ("int" "int8" "int16" "int32" "int64"))
(obj "unsigned int type" :enum ("uint" "uint8" "uint16" " uint32" "uint64" "uintptr"))
(obj "int type" :union ("signed int type" "unsigned int type"))
(obj "float type" :enum ("float32" "float64"))
(obj "real-valued type" :enum ("int type" "float type"))
(obj "complex type" :enum ("complex32" "complex64"))
(obj "byte type" :enum ("byte"))  ; <doc> alias for uint8
(obj "rune type" :enum ("rune"))  ; <doc> alias for int32
(obj "numeric type" :union ("real-valued type" "complex type" "byte type" "rune type"))
(obj "string type" :enum ("string"))
(obj "base type" :union ("boolean type" "numeric type" "string type"))
;composite types
(obj "array type" :at "len" nat :at "element type" "type")
(obj "slice type" :at "underlying array" "array type")  ; for decl :at "element type" "type" :at "start index" "nat" :at "len" nat :at "cap" nat
(obj "named field" :at "name" "identifier" :at "type" "type")
(obj "embedded field" :at "type" "type")
(obj "field decl" :union ("named field" "embedded field"))
(obj "struct type" :at "fields" (listt "field"))
(obj "pointer type" :at "type" "type")
(obj "function type")  ; todo
(obj "interface type")  ; todo
(obj "map type")  ; todo
(obj "channel type")  ; todo
(obj "composite type" :union ("array type" "slice type" "struct type" "pointer type" "function type" "interface type" "map type" "channel type"))
(obj "type" :union ("base type" "composite type"))


;; Operators
;unary operators
(obj "+." :at "arg" "expression") ; +x == 5
(obj "-." :at "arg" "expression") ; -x == -1 * x
(obj "logical not operator" :at "arg" "expression") ; !true == false
(obj "bin not unary operator" :at "arg" "expression") ; ^10(2) == 01(2)
(obj "*." :at "arg" "expression")  ; *p
(obj "&" :at "arg" "expression")  ; &x
;<- operator ---
(obj "unary operator" :union ("add unary operator" "sub unary operator" "logical not operator" "bin not unary operator" "dereference unary operator" "address unary operator"))
;binary multiplication operators
(obj "mul mul operator" :at "arg1" "expression" :at "arg2" "expression") ; 5 * 10 == 50
(obj "div mul operator" :at "arg1" "expression" :at "arg2" "expression") ; 4 / 2 = 2.0
(obj "remainder of div mul operator" :at "arg1" "expression" :at "arg2" "expression") ; 5 % 2 == 1
(obj "bin left shift mul operator" :at "arg1" "expression" :at "arg2" "expression") ; 1 << 12 == 2 ** 12
(obj "bin right shift mul operator" :at "arg1" "expression" :at "arg2" "expression") ; 32 >> 2 == 2 ** 3
(obj "bin and mul operator" :at "arg1" "expression" :at "arg2" "expression") ; 1100(2) & 1010(2) == 1000(2)
(obj "bin and not mul operator" :at "arg1" "expression" :at "arg2" "expression") ; 11001010(2) &^ 10101100(2) == 01000010(2) ; 1 <=> 1 &^ 0
(obj "mul operator" :union ("mul mul operator" "div mul operator" "remainder of div mul operator" "bin left shift mul operator" "bin right shift mul operator" "bin and mul operator" "bin and not mul operator"))
;binary addition operators
(obj "add add operator" :at "arg1" "expression" :at "arg2" "expression") ; 2 + 3 == 5
(obj "sub add operator" :at "arg1" "expression" :at "arg2" "expression") ; 2 - 3 == -1
(obj "bin or add oeprator" :at "arg1" "expression" :at "arg2" "expression") ; 1100(2) | 1010(2) == 1110(2)
(obj "bin xor add operator" :at "arg1" "expression" :at "arg2" "expression") ; 1100(2) ^ 1010(2) == 0110(2)
(obj "add operator" :union ("add add operator" "sub add operator" "bin or add operator" "bin xor add operator"))
;relation operators
(obj "equality operator" :at "arg1" "expression" :at "arg2" "expression")  ; x == 3 | 4 == 5 | x == y
(obj "inequality operator" :at "arg1" "expression" :at "arg2" "expression")  ; x != 5
(obj "less than operator" :at "arg1" "expression" :at "arg2" "expression")  ; x < 3 | 3 < 5 | x < y
(obj "no greater than operator" :at "arg1" "expression" :at "arg2" "expression")  ; x <= 3 | 3 <= 4 | x <= y
(obj "greater than operator" :at "arg1" "expression" :at "arg2" "expression")  ; x > 3 | 4 > 5 | x > y
(obj "no less than operator" :at "arg1" "expression" :at "arg2" "expression")  ; x >= 3 | 3 >= 4 | x >= y
(obj "unequality operator" :union ("less than expression" "greater than expression" "no less than expression" "no greater than expression"))
(obj "relation operator" :union ("equality expression" "inequality operator" "unequality operator"))
(obj "logical or operator" :at "arg1" "expression" :at "arg2" "expression")  ; true || false == false
(obj "logical and operator" :at "arg1" "expression" :at "arg2" "expression")  ; true && false == false
(obj "binary operator" :union ("logical or operator" "logical and operator" "rel operator" "add operator" "mul operator"))
(obj "expression" :union ("constant" "variable" "unary expression" "binary expression"))


;; Expressions
(obj "operand" :union ("constant" "variable"))  ;...
(obj "primary expression" :union ("operand"))  ;...
(obj "function literal" :at "args" (listt "variable") :at "return type" "type" :at "body" "expression")  ; function literals alias for anonymous functions
(obj "unary expression" :union ("primary expression" "unary operator")) ; <doc> PrimaryExpr | unary_op UnaryExpr
(obj "expression" :union ("unary expression" "binary operator")) ; <doc> UnaryExpr | Expression binary_op Expression


;; Statements
(obj "empty statement")  ; <doc> The empty statement does nothing
(obj "label statement" :at "label" string :at "statement" "statement")  ; <example> Error: log.Panic("error encountered") ???
(obj "expression statement" :union ("expression"))
(obj "send statement" :at "channel" "expression" :at "message" "expression")
(obj "inc statement" :at "expression" "expression")  ; x++
(obj "dec statement" :at "expression" "expression")  ; x--
(obj "assignment statement" :at "location" "expression" :at "expression" "expression")  ; check names of atrributes
(obj "simple statement" :union ("empty statement" "expression statement" "inc statement" "dec statement" "assignment statement"))  ;...
(obj "if else type" :union ("if statement" "block"))
(obj "if statement" :at "init statements" (listt "simple statement") :at "condition" :union ("expression") :at "then" "statement" :at "else" "if else type")  ; if i := 1; i + 3; i < 5 { <then> }
;(obj "switch statement")  ; ---
(obj "for statement" :at "init statement" "simple statement" :at "condition" "expression" :at "post statement" "simple statement") ; var declaration
(obj "statement" :union ("label statement" "simple statement" "if statement" "for statement"))


;; Other
(obj "package" :union ("identifier"))  ; package main
(obj "library" :union (string))  ; import "fmt"


;; For the future
#|
(obj "variable declaration" :at "name" "identifier" :at "type" "type" :at "value" "expression")  ; Need to check: (is-instance (aget variable "value") (aget variable "type"))
(obj "array" :at "type" "array type" :at "value" (listt any))
(obj "slice" :at "type" "slice type" :at "start index" nat :at "length" nat :at "capacity" nat)
|#
