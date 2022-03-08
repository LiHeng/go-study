package filetype

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/h2non/filetype"
)

func Test_filetype(t *testing.T) {
	buf, _ := ioutil.ReadFile("testdata.csv")
	kind, _ := filetype.Match(buf)
	fmt.Println(kind)

	buf, _ = ioutil.ReadFile("testdata.html")
	kind, _ = filetype.Match(buf)
	fmt.Println(kind)

	buf, _ = ioutil.ReadFile("lots_sheets.xlsx")
	kind, _ = filetype.Match(buf)
	fmt.Println(kind)
}
