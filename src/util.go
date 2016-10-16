package main

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"github.com/fatih/color"
	"os"
	"path/filepath"
	"text/template"
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
func hash(s string) string {
	hasher := md5.New()
	hasher.Write([]byte(s))
	return hex.EncodeToString(hasher.Sum(nil))
}

func checkErr(err error) {
	if err != nil {
		logFatal("%v", err)
	}
}

func fillTemplate(s string, mapping Mapping) string {
	var r bytes.Buffer
	tmpl, err := template.New("MyTemplate").Parse(s)
	checkErr(err)
	err = tmpl.Execute(&r, mapping)
	checkErr(err)
	return r.String()
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func (saidumlo *SaiDumLo) getSecretMeta(filePath string) (SecretGroup, Secret) {
	pwd, e := os.Getwd()
	checkErr(e)
	rel, e := filepath.Rel(saidumlo.ConfigDir, pwd+"/"+filePath)
	checkErr(e)

	for _, secretGroup := range saidumlo.Config.SecretGroups {
		for _, secret := range secretGroup.Secrets {
			if secret.Encrypted == rel {
				return secretGroup, secret
			}
		}
	}

	// TODO: better
	logFatal("Could not find declaration of %s", filePath)
	os.Exit(1)
	return SecretGroup{}, Secret{}
}
