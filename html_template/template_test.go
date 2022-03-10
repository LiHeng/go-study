package html_template

import (
	"html/template"
	"os"
	"testing"
)

func TestHtmlTemplate(t *testing.T) {
	tplStr := "{{.Name}}</br>{{.Age}}\n"
	tpl := template.Must(template.New("tplName").Parse(tplStr))
	err := tpl.Execute(os.Stdout, map[string]interface{}{"Name": "big_cat", "Age": 29})
	if err != nil {
		t.FailNow()
	}
}
