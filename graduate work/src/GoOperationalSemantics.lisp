#|
	Operational semantics for the Go language

    Last edit: 17/10/2025
|#


(obj "env" :at "agents" (listt "agent"))


(obj "agent"
    :at "variable location" (cobject "variable" "location")
    :at "location value" (cobject "location" "Go value")
    :at "location type" (cobject "location" "type")
)


;;; Auxiliary generic objects
(obj "location")
(obj "Go value" :union ("constant" "array" "slice" "struct" "function" "interface" "map" "expression" "statement"))
