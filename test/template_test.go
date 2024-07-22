package test

import (
	template2 "auto-record/app/template"
	"fmt"
	"testing"
)

func TestGetAllTemplates(t *testing.T) {
	template := template2.NewTemplate()
	ans := template.GetAll()
	for _, item := range ans {
		fmt.Println(item.Name)
	}
}
