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

func TestCreateTemplate(t *testing.T) {
	name := ""
	template := template2.NewTemplate()
	isSuccess, err := template.CreateTemplate(name)
	if err != nil {
		t.Fatalf("让我看看到底是什么是垃圾报错： %v", err)
	}
	t.Logf("创建模板成功 你太酷了, %v", isSuccess)
}
