package ast

import (
	"fmt"
)

// Unique Identifier to be used in Identifier Expresions,
// types and algorithmic contexts
type UniqueIdentifier struct {
	Value string
	Id    int
}

// Display identifier only
func (u UniqueIdentifier) String() string {
	return fmt.Sprintf("%s", u.Value)
}

// Also display numeric ID
func (u UniqueIdentifier) FullString() string {
	return fmt.Sprintf("(%s,%d)", u.Value, u.Id)
}

// To increment every time an UID is generated
var uid_global_counter int = 1

// Generate a new UID incrementing the global counter
func GenUID(name string) UniqueIdentifier {
	uid := UniqueIdentifier{name, uid_global_counter}
	uid_global_counter++
	return uid
}

func ResetUIDCounter() {
	uid_global_counter = 1
}
