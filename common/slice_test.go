package common

import (
	"fmt"
	"testing"
)

func TestSlice(t *testing.T) {
	a := []string{"1", "2", "3", "3", "2"}
	result, err := RemoveDup(a)
	if err != nil {
		t.Fatal("error")
	}
	fmt.Println(result)
}
