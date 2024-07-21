package test

import (
	"auto-record/config"
	"fmt"
	"testing"
)

func TestViperConfig(t *testing.T) {
	settings, err := config.NewViperConfig()
	if err != nil {
		return
	}
	application := settings.Application.Name
	fmt.Println(application)
	if application != "go-auto" {
		t.Errorf("unread APPLICATION")
	}
}
