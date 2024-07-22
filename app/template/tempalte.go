package template

import (
	"auto-record/config"
	"fmt"
	"os"
)

type Template struct {
	Name string
}

func NewTemplate() *Template {
	return &Template{}
}

func (t *Template) GetAll() []*Template {
	settings := config.Settings
	directory := settings.FilePath.Record
	ans := make([]*Template, 0)

	rd, err := os.ReadDir(directory)
	if err != nil {
		fmt.Printf("read dir fail: %v \n", err)
	}
	for _, dir := range rd {
		if dir.IsDir() {
			ans = append(ans, &Template{
				Name: dir.Name(),
			})
		}
	}
	return ans
}
