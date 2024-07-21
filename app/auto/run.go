package auto

import (
	"auto-record/app/event"
	"auto-record/config"
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func (ar *AutoRecord) Run() {
	fmt.Println("Run")
	filename := filepath.Join(config.Settings.FilePath.Record, "text.log")
	ioReader, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	scanner := bufio.NewScanner(ioReader)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), "HookEnabled") {
			continue
		}
		event.MouseEventFormat(scanner.Text())
		event.KeyboardEventFormat(scanner.Text())
	}

}
