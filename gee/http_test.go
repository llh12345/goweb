package gee

import (
	"fmt"
	"strings"
	"testing"
)

func Test(t *testing.T) {
	s := "/test/"
	parts := strings.Split(s, "/")
	fmt.Println(len(parts))
}
