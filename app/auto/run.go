package auto

import (
	"auto-record/app/event"
	"bufio"
	"fmt"
	"os"
	"strings"
)

func (ar *AutoRecord) Run() {
	fmt.Println("Run")
	filename := "text.log"
	ioReader, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	scanner := bufio.NewScanner(ioReader)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
		if strings.Contains(scanner.Text(), "HookEnabled") {
			continue
		}
		event.MouseEventFormat(scanner.Text())
		event.KeyboardEventFormat(scanner.Text())
	}

}
