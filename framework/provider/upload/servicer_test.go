package upload

import (
	"io"
	"math"
	"strings"
	"testing"
)

func TestAdd(t *testing.T) {
	var strReader interface{} = strings.NewReader("hello world")
	reader := strReader.(io.Reader)

	if err := upload("/Users/baby/go/src/github.com/hade/files", "test.txt", math.MaxInt64, &reader); err != nil {
		t.Error(err)
	}
}
