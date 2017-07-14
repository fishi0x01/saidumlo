package main

import (
	"github.com/fatih/color"
)

/******
* Logging
******/
func logDebug(s string, args ...interface{}) {
	if verbose {
		color.Blue("[DEBUG] "+s, args...)
	}
}

func logInfo(s string, args ...interface{}) {
	color.Green("[INFO] "+s, args...)
}

func logError(s string, args ...interface{}) {
	color.Red("[ERROR] "+s, args...)
}

func logFatal(s string, args ...interface{}) {
	color.Red("[FATAL] "+s, args...)
}

/*****
* Helpers
******/
func checkErr(err error) {
	if err != nil {
		logFatal("%v", err)
	}
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// TODO: generic return type
func getMapKeys(m map[string]SecretGroup) []string {
	var keys []string

	for key := range m {
		keys = append(keys, key)
	}

	return keys
}
