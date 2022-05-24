package main

import (
	"fmt"

	"github.com/charliemcelfresh/event_worker/cmd"
)

func init() {
	fmt.Println("Running main config")
}

func main() {
	cmd.Execute()
}
