#|
Operational semantics of Go
|#

(cpath ("Go" "names"))
(obj "name" :union (symbol))  ; letter { letter | unicode_digit } ???

(cpath ("Go" "packages"))
(obj "package" :union ("name"))  ; package main

(cpath ("Go" "libraries"))
(obj "library" :union (string))  ; import "fmt" ???

(cpath ("Go" "variables"))
(obj "variable" :union ("name"))  ; x | y | p  ; :at "type" "type" ???

(cpath ("Go" "constants"))
(obj "constant" :union (int))  ; 0 | 2  ; :at "type" "type" ???

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
(obj "pointer type" :at "type" "type")
(obj "type" :union ("bool type" "numeric type" "string type" "pointer type"))

(cpath ("Go" "expressions"))
(obj "address expression" :at "variable" "variable")  ; &x
(obj "dereference expression" :at "variable" "variable")  ; *p
(obj "assignment expression" :at "location" "variable" :at "expression" "expression")  ; x = 3 | x := 3 | p = &x
(obj "equality expression" :at "arg1" "expression" :at "arg2" "expression")  ; x == 3
(obj "less than expression" :at "arg1" "expression" :at "arg2" "expression")  ; x < 3
(obj "greater than expression" :at "arg1" "expression" :at "arg2" "expression")  ; x > 3
(obj "no less than expression" :at "arg1" "expression" :at "arg2" "expression")  ; x >= 3
(obj "no greater than expression" :ar "arg1" "expression" :at "arg2" "expression")  ; x <= 3
(obj "unequality expression" :union ("equality expression" "less than expression" "greater than expression" "no less than expression" "no greater than expression"))
(obj "expression" :union ("variable" "constant" "address expression" "dereference expression" "assignment expression" "unequality expression"))

(cpath ("Go" "statements"))
(obj "expression statement" :at "expression" "expression")
(obj "block statement" :at "statements" (listt "statement"))
(obj "if statement simple" :at "condition" "expression" :at "then" "statement")  ; if <condition> { <then> }
(obj "if-else statement" :at "condition" "expression" :at "then" "statement" :at "else" "expression")  ; if <condition> { <then> } else { <else> }
(obj "if statement" :union ("is statement simple" "if-else statement"))
(obj "variable declaration simple" :at "type" "type" :at "name" "name")  ; var x int
(obj "variable declaration full" :at "type" "type" :at "name" "name" :at "initiliazer" "expression")  ; var x int = 3
(obj "variable declaration short" :at "name" "name" :at "initiliazer" "expression")  ; x := 0 ???
(obj "variable declaration" :union ("variable declaration simple" "variable declaration full" "variable declaration short"))
(obj "statement" :union ("expression statement" "block statement" "if statement" "variable declaration"))

