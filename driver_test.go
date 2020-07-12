package compclean

import (
	"fmt"
	"testing"
)

func TestDriver(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping testing in short mode")
	}
	err := Drive()
	fmt.Println("here")
	if err != nil {
		t.Fatal(err) 
	}
}
