// Package output defines some simple output filters for text.
//
// This theme is based on various "hack scripts" that have been published on the
// internet.
package output

import "fmt"

// WriteHeader writes a section header.
//
//    [-] Intializing Low Orbit Cannon
func WriteHeader(text interface{}) {
	fmt.Printf("\n[-] %s\n", text)
}

// WriteData spaces data out and writes it.
//
//        done
func WriteData(text interface{}) {
	fmt.Printf("    %s\n", text)
}

// WriteError signifies an error condition.
//
//    [!] Error: no lazors!!!
func WriteError(text interface{}) {
	fmt.Printf("[!] %s\n", text)
}

// WriteEnd signifies the end of the build.
//
//    [=] Build finished
func WriteEnd(text interface{}) {
	fmt.Printf("[=] %s\n", text)
}
