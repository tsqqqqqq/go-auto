package test

import (
	window2 "auto-record/app/window"
	"fmt"
	"testing"
)

func TestGetAllWindows(t *testing.T) {
	template := window2.NewWindow()
	ans := template.GetWindows()
	for _, item := range ans {
		fmt.Println(item)
	}
}
