package pages

import (
	"html/template"
	"io"
)


func RenderPage(w io.Writer, absolutePagePath string, selfEntity interface{}) {
	t, err := template.ParseFiles(absolutePagePath)
	if err == nil {
		t.Execute(w, selfEntity)
	}
}
