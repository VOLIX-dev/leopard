package database

import (
	"fmt"
	"testing"
)

func TestQueryBuilder(t *testing.T) {

	fmt.Println(
		NewQueryBuilder("test").
			Select("test.hi", "test.bye").
			Where("hi", "=", "hello").
			OrWhere("bye", "=", "goodbye").
			Where("hi", "=", "hello").
			OrderBy("hi", "desc").
			Build(),
	)
}
