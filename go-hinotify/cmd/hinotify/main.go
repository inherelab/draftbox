package main

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/gookit/color"
	"github.com/gookit/goutil/cliutil"
	"github.com/gookit/goutil/jsonutil"
	"github.com/inhere/hinotify"
)

var opts = struct {
	fontName   string
	visualMode bool
	list       bool
	sample     bool
}{}

// task lists
var tasks map[string]interface{}

// test run:
//  go build ./_examples/alone && ./alone -h
//  go run ./cmd/hinotify
func main() {
	cmd := hinotify.FileWatcher(dispatchTasks)

	// load config
	err := jsonutil.ReadFile(".hinotify.json", hinotify.Options())
	if err != nil {
		color.Error.Println("Read config error:", err.Error())
		return
	}

	fmt.Println(hinotify.Options())

	// Alone Running
	cmd.MustRun(nil)
	// cmd.Run(os.Args[1:])
}

func dispatchTasks(ev fsnotify.Event) {
	opts := hinotify.Options()

	for name, items := range opts.Tasks {
		switch itemTyp := items.(type) {
		case string:
			cliutil.QuickExec(itemTyp)
		case []string:

		}

		fmt.Println(name, items)
	}
}

func executeTask()  {

}
