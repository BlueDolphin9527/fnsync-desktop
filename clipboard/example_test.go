package clipboard_test

import (
	"fmt"
	"log"

	"github.com/cxfksword/fnsync-desktop/clipboard"
)

// The purpose of this example is to demonstrate the implementation of get and set.
func Example() {
	if err := clipboard.Set("gopher"); err != nil {
		log.Fatal(err)
	}

	text, err := clipboard.Get()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(text)

	// Output:
	// gopher
}
