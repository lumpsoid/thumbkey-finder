package keyboard

import "fmt"

// Define a custom error type
type AlreadyExist struct {
	Message string
}

// Implement the Error() method for the custom error type
func (e AlreadyExist) Error() string {
	return fmt.Sprintf("Error %s", e.Message)
}

// Function to check if the error is of type AlreadyExist
func IsAlreadyExist(err error) bool {
	_, ok := err.(AlreadyExist)
	return ok
}
