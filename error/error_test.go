package error

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrorObj(t *testing.T) {
	err := errors.New("")
	assert.NotNil(t, err)

	err = errors.New("not found")
	err1 := fmt.Errorf("some context: %v", err)
	err2 := fmt.Errorf("some context: %w", err)

	assert.NotEqual(t, reflect.TypeOf(err1), reflect.TypeOf(err2))
	assert.Equal(t, err1.Error(), err2.Error())

	if _, ok := err2.(interface{ Unwrap() error }); ok {
		fmt.Printf("err2 is wrap error\n")
	}

	unwrapErr := errors.Unwrap(err2)
	assert.Equal(t, err, unwrapErr)

	assert.True(t, errors.Is(err2, err))
}

func AssertError(err error) {
	if e, ok := err.(*os.PathError); ok {
		fmt.Printf("it's an os.PathError, operation:%s, path:%s, msg: %v\n", e.Op, e.Path, e.Err)
	}
}

func TestAssertError(t *testing.T) {
	err1 := &os.PathError{
		Op:   "write",
		Path: "/root/demo.txt",
		Err:  os.ErrPermission,
	}
	AssertError(err1)

}
