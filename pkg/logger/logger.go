package logger

import (
	"fmt"
	"log"
)

var logFatal = log.Fatal
var logFatalf = log.Fatalf

var Debug bool

// StatusTagWidth is the fixed column width reserved for a bracketed tag (e.g. "[IMAGES]"), so the
// symbol and message columns built by StatusFmt always line up regardless of tag length.
const StatusTagWidth = 9

// StatusFmt builds a "tag\tsymbol\tformat" line with tag padded to StatusTagWidth. The result is
// itself a format string: any verbs inside format are left untouched for the caller's Printf call.
func StatusFmt(tag, symbol, format string) string {
	return fmt.Sprintf("%-*s\t%s\t%s", StatusTagWidth, tag, symbol, format)
}

// syncFlags keeps the timestamp in sync with Debug: only debug output needs it, plain logs read better without it.
func syncFlags() {
	if Debug {
		log.SetFlags(log.LstdFlags)
	} else {
		log.SetFlags(0)
	}
}

func Debugf(format string, v ...interface{}) {
	syncFlags()
	if Debug {
		log.Printf("[DEBUG] "+format, v...)
	}
}

func Printf(format string, v ...interface{}) {
	syncFlags()
	log.Printf(format, v...)
}

func Println(v ...interface{}) {
	syncFlags()
	log.Println(v...)
}

func Fatal(v ...interface{}) {
	syncFlags()
	logFatal(v...)
}

func Fatalf(format string, v ...interface{}) {
	syncFlags()
	logFatalf(format, v...)
}
