package template

import (
	"auto-record/config"
	"errors"
	"fmt"
	"os"
	"path"
)

type Template struct {
	Name string
}

func NewTemplate() *Template {
	return &Template{}
}

var CurrentTemplate string

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

func (t *Template) ChangeCurrentTemplate(template string) {
	CurrentTemplate = template
}

func (t *Template) CreateTemplate(name string) (bool, error) {
	settings := config.Settings
	directory := settings.FilePath.Record

	filePath := path.Join(directory, name)
	// if err == nil , then file / explorer is already existed
	_, err := os.Stat(filePath)
	if err == nil {
		return false, errors.New(fmt.Sprintf("template is exist, please check file / dir %s is already", filePath))
	}
	if os.IsNotExist(err) {
		err := os.Mkdir(filePath, os.ModePerm)
		if err != nil {
			return false, err
		}
	}
	return true, nil
}
