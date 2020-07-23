package dao

import (
	"fmt"
	"testing"
)

func TestSetName(t *testing.T) {

	err := SetName("222")

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(111)
}
