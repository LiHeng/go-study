package tag

import (
	"fmt"
	"reflect"
	"testing"
)

type User struct {
	Name     string `master:"username" json:"Name"`
	Age      int    `master:"age"`
	Password string `master:"password"`
}

func getTag(user User) {
	t := reflect.TypeOf(user)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get("master")
		fmt.Println("get tag is ", tag)
		value, ok := field.Tag.Lookup("json")
		if ok {
			fmt.Println("json tag is ", value)
		}
	}
}

func TestTag(t *testing.T) {
	u := User{
		Name:     "masterliheng",
		Age:      5,
		Password: "123456",
	}
	getTag(u)
}
