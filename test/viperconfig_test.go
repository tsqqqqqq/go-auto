package test

import (
	"auto-record/config"
	"fmt"
	"testing"
)

func TestViperConfig(t *testing.T) {
	settings := config.Settings
	application := settings.Application.Name
	recordfile := settings.FilePath.Record
	fmt.Println(application)
	fmt.Println(recordfile)
	if application != "go-auto" {
		t.Errorf("unread APPLICATION")
	}
}
