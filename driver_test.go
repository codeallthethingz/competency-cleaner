package compclean

import (
	"testing"
)

func TestDriver(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping testing in short mode")
	}
	err := Drive()
	if err != nil {
		t.Fatal(err)
	}
}
