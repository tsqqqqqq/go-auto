package test

import (
	"auto-record/utils"
	"fmt"
	"testing"
)

func TestFileUtils(f *testing.T) {
	dir, err := utils.Dirname()
	if err != nil {
		f.Fatal(err)
	}
	filename, err := utils.Filename()
	rootname, err := utils.Rootname()

	fmt.Println(rootname)
	fmt.Println(filename)
	fmt.Println(dir)
}
