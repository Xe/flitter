// Package output defines some simple heroku-like output filters for text.
package output

import "fmt"

// WriteHeader writes a heroku style header.
func WriteHeader(text interface{}) {
	fmt.Printf("-----> %s", text)
}

// WriteData spaces data out and writes it.
func WriteData(text interface{}) {
	fmt.Printf("       %s", text)
}

// WriteError signifies an error condition.
func WriteError(text interface{}) {
	fmt.Printf("     ! %s", text)
}
