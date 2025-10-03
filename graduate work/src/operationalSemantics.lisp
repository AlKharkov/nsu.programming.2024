#|
Operational semantics of Go
|#

(cpath ("Go" "comments"))
(obj "line comment") ; //
(obj "general comment") ; /* multi-line comment */

(cpath ("Go" "names"))
(obj "name" :union (symbol))  ; letter { letter | unicode_digit } ???

(cpath ("Go" "packages"))
(obj "package" :union ("name"))  ; package main

(cpath ("Go" "libraries"))
(obj "library" :union (string))  ; import "fmt" ???

(cpath ("Go" "variables"))
(obj "variable" :union ("name"))  ; x | y | p  ; :at "type" "type" ???

(cpath ("Go" "constants"))
(obj "constant" :union (int bool real))  ; 0 | 2  ; :at "type" "type" ???

(cpath ("Go" "types"))
(obj "bool type" :enum ("bool"))
(obj "signed int type" :enum ("int" "int8" "int16" "int32" "int64"))
(obj "unsigned int type" :enum ("uint" "uint8" "uint16" " uint32" "uint64" "uinttpr"))
(obj "int type" :union ("signed int type" "unsigned int type"))
(obj "float type" :enum ("float32" "float64"))
(obj "real-valued type" :union ("int type" "float type"))
(obj "complex type" :enum ("complex32" "complex64"))
(obj "numeric type" :union ("real-valued type" "complex type"))
(obj "string type" :union ("string"))
(obj "byte type" :enum ("byte"))  ; alias for uint8
(obj "rune type" :enum ("rune"))  ; alias for int32
(obj "base type" :union ("bool type" "numeric type" "string type" "byte type" "rune type"))
(obj "pointer type" :at "type" "type")
(obj "array type" :at "type" "type" :at "len" int)  ; ???
(obj "type" :union ("base type" "pointer type" "array type"))

(cpath ("Go" "operators"))
; unary operators
(obj "add unary operator" :at "arg" "expression") ; +x == 5
(obj "sub unary operator" :at "arg" "expression") ; -x == -1 * x
(obj "factorial unary operator" :at "arg" "expression") ; 5! == 120
(obj "address unary operator" :at "arg" "expression")  ; &x
(obj "dereference unary operator" :at "arg" "expression")  ; *p
(obj "bin not unary operator" :at "arg" "expression") ; ^10(2) == 01(2)
(obj "unary operator" :union ("add unary operator" "sub unary operator" "factorial unary operator" "address unary operator" "dereference unary operator" "bin not unary operator"))
; binary multiplication operators
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
(obj "less than operator" :at "arg1" "expression" :at "arg2" "expression")  ; x < 3 | 3 < 5 | x < y
(obj "greater than operator" :at "arg1" "expression" :at "arg2" "expression")  ; x > 3 | 4 > 5 | x > y
(obj "no less than operator" :at "arg1" "expression" :at "arg2" "expression")  ; x >= 3 | 3 >= 4 | x >= y
(obj "no greater than operator" :at "arg1" "expression" :at "arg2" "expression")  ; x <= 3 | 3 <= 4 | x <= y
(obj "unequality operator" :union ("less than expression" "greater than expression" "no less than expression" "no greater than expression"))
(obj "rel operator" :union ("equality expression" "unequality expression"))
(obj "binary operator" :union ("mul operator" "add operator" "rel operator"))
(obj "operator" :union ("unary operator" "binary operator"))

(cpath ("Go" "expressions"))
(obj "assignment expression" :at "location" "variable" :at "expression" "expression")  ; x = 3 | x := 3 | p = &x
(obj "expression" :union ("variable" "constant" "operator" "assignment expression"))

(cpath ("Go" "statements"))
(obj "expression statement" :at "expression" "expression")
(obj "block statement" :at "statements" (listt "expression"))  ; "statement" ???
(obj "if statement" :at "condition" "expression" :at "then" "expression" :at "else" "expression")  ; if <condition> { <then> } else { <else> }
(obj "variable declaration" :at "type" "type" :at "name" "name" :at "initiliazer" "expression")  ; var x int = 3
(obj "for statement" :at "init statement" "expression" :at "condition" "expression" :at "post statement" "expression" :at "body" "expression")
(obj "statement" :union ("expression statement" "block statement" "if statement" "variable declaration"))

