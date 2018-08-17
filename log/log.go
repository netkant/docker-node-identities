package log

import (
	"fmt"
)

// LogSettings ...
type LogSettings struct {
	Debug bool
}

// ...
var (
	Settings *LogSettings
)

// Debug ...
func Debug(format string, a ...interface{}) {
	if Settings.Debug {
		fmt.Printf("DEBUG: "+format+"\n", a...)
	}
}

// Info ...
func Info(format string, a ...interface{}) {
	fmt.Printf("INFO: "+format+"\n", a...)
}
