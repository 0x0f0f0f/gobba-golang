package ast

// Unique Identifier to be used in Identifier Expresions,
// types and algorithmic contexts
type UniqueIdentifier struct {
	Value string
	Id    int
}

// To increment every time an UID is generated
var global_counter *int
