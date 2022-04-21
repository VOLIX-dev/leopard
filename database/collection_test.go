package database

import (
	"fmt"
	"testing"
)

func TestGetValues(t *testing.T) {

	collection := New(map[int]string{
		0: "zero",
		1: "one",
	})

	fmt.Println(collection.GetValues())
}
