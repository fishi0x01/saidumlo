package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
)

/*****
* Config
******/
// TODO: prettify
func (saidumlo *SaiDumLo) parseConfig(configFileName string) {
	saidumlo.ConfigFileName = configFileName
	saidumlo.getConfigDir()
	saidumlo.Config = Config{}
	configFilePath := saidumlo.ConfigDir + saidumlo.ConfigFileName
	logInfo("Using config %v", configFilePath)
	s, e := ioutil.ReadFile(configFilePath)
	checkErr(e)
	e = yaml.Unmarshal(s, &saidumlo.Config)
	checkErr(e)

	logDebug("%#v", saidumlo.Config)
}

// TODO: prettify
// TODO: maybe return relative path instead of total
func (saidumlo *SaiDumLo) getConfigDir() {
	var dirPrefix = "./"
	var result string
	for {
		if _, err := os.Stat(dirPrefix + saidumlo.ConfigFileName); os.IsNotExist(err) {
			dirPrefix += "../"
			var curDir, _ = filepath.Abs(filepath.Dir(dirPrefix))
			if curDir == "/" {
				logFatal("Could not find %v", saidumlo.ConfigFileName)
				os.Exit(1)
			}
		} else {
			result, err = filepath.Abs(filepath.Dir(dirPrefix + saidumlo.ConfigFileName))
			checkErr(err)
			break
		}
	}

	saidumlo.ConfigDir = result + "/"
}

// SaidumloCacheDir directory for tmp decrypted files
const SaidumloCacheDir = ".saidumlo_cache/"

func (saidumlo *SaiDumLo) cleanCache() {
	os.RemoveAll(saidumlo.ConfigDir + SaidumloCacheDir)
}

func (saidumlo *SaiDumLo) backup(filePath string) string {
	dstDir := saidumlo.ConfigDir + SaidumloCacheDir + string(hash(filePath)) + "/"
	_, file := path.Split(filePath)
	logDebug("Backup %v to %v", saidumlo.ConfigDir+filePath, dstDir+file)
	os.MkdirAll(dstDir, 0700)
	data, err := ioutil.ReadFile(saidumlo.ConfigDir + filePath)
	checkErr(err)
	err = ioutil.WriteFile(dstDir+file, data, 0700)
	checkErr(err)
	return dstDir + file
}

func (saidumlo *SaiDumLo) restore(filePath string) string {
	logDebug("Restore %v", saidumlo.ConfigDir+filePath)
	_, file := path.Split(filePath)
	srcFile := saidumlo.ConfigDir + SaidumloCacheDir + string(hash(filePath)) + "/" + file
	data, err := ioutil.ReadFile(srcFile)
	checkErr(err)
	err = ioutil.WriteFile(filePath, data, 0754)
	checkErr(err)
	return filePath
}

/*****
* Editing
*****/
func (saidumlo *SaiDumLo) editSecret(filePath string) {
	secretGroup, secret := saidumlo.getSecretMeta(filePath)

	encBu := saidumlo.backup(secret.Encrypted)
	decBu := encBu + ".dec"
	if secret.Encrypted == secret.Decrypted {
		decBu = encBu
	}

	saidumlo.decryptFile(secretGroup, encBu, decBu)
	saidumlo.startEditor(decBu)
	saidumlo.encryptFile(secretGroup, encBu, decBu)
	e := os.Rename(encBu, saidumlo.ConfigDir+secret.Encrypted)
	checkErr(e)
}

// DefaultEditor editor to use if $EDITOR not set
const DefaultEditor = "vim"

func (saidumlo *SaiDumLo) startEditor(filePath string) {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		logInfo("Variable '$EDITOR' not set - using vim by default")
		editor = DefaultEditor
	}

	cmd := exec.Command(editor, filePath)
	cmd.Env = os.Environ()
	cmd.Dir = saidumlo.ConfigDir
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	e := cmd.Run()
	checkErr(e)
}

/*****
* Decryption
*****/
func (saidumlo *SaiDumLo) decryptSecretGroup(secretGroup SecretGroup) {
	for _, secret := range secretGroup.Secrets {
		logInfo("Processing '%s'", secret.Encrypted)
		saidumlo.backup(secret.Encrypted)
		txt, err := saidumlo.decryptFile(secretGroup, saidumlo.ConfigDir+secret.Encrypted, saidumlo.ConfigDir+secret.Decrypted)

		if err != nil {
			logError("%s: %s", err, string(txt))
			saidumlo.restore(secret.Encrypted)
		}
	}
}

func (saidumlo *SaiDumLo) decryptFile(secretGroup SecretGroup, encFilePath string, decFilePath string) ([]byte, error) {
	mapping := Mapping{secretGroup.Key, encFilePath, decFilePath}
	executable := fillTemplate(saidumlo.Config.Encryptions[secretGroup.Encryption].Decrypt.Executable, mapping)
	var args []string
	for _, arg := range saidumlo.Config.Encryptions[secretGroup.Encryption].Decrypt.Args {
		args = append(args, fillTemplate(string(arg), mapping))
	}

	envs := os.Environ()
	for _, env := range saidumlo.Config.Encryptions[secretGroup.Encryption].Decrypt.Env {
		envs = append(envs, env.Name+"="+fillTemplate(string(env.Value), mapping))
	}

	logDebug("Environment: %s", strings.Join(envs[:], " "))
	logDebug("Command: %s", executable+" "+strings.Join(args[:], " "))

	cmd := exec.Command(executable, args...)
	cmd.Env = envs
	cmd.Dir = saidumlo.ConfigDir
	return cmd.CombinedOutput()
}

/****
* Encryption
****/
func (saidumlo *SaiDumLo) encryptSecretGroup(secretGroup SecretGroup) {
	for _, secret := range secretGroup.Secrets {
		logInfo("Processing '%s'", secret.Decrypted)
		saidumlo.backup(secret.Decrypted)
		txt, err := saidumlo.encryptFile(secretGroup, saidumlo.ConfigDir+secret.Encrypted, saidumlo.ConfigDir+secret.Decrypted)

		if err != nil {
			logError("%s: %s", err, string(txt))
			saidumlo.restore(secret.Decrypted)
		}
	}
}

func (saidumlo *SaiDumLo) encryptFile(secretGroup SecretGroup, encFilePath string, decFilePath string) ([]byte, error) {
	mapping := Mapping{secretGroup.Key, decFilePath, encFilePath}
	executable := fillTemplate(saidumlo.Config.Encryptions[secretGroup.Encryption].Encrypt.Executable, mapping)
	var args []string
	for _, arg := range saidumlo.Config.Encryptions[secretGroup.Encryption].Encrypt.Args {
		args = append(args, fillTemplate(string(arg), mapping))
	}

	envs := os.Environ()
	for _, env := range saidumlo.Config.Encryptions[secretGroup.Encryption].Encrypt.Env {
		envs = append(envs, env.Name+"="+fillTemplate(string(env.Value), mapping))
	}

	logDebug("Environment: %s", strings.Join(envs[:], " "))
	logDebug("Command: %s", executable+" "+strings.Join(args[:], " "))

	cmd := exec.Command(executable, args...)
	cmd.Env = envs
	cmd.Dir = saidumlo.ConfigDir
	return cmd.CombinedOutput()
}

/******
* Command Controller
******/
var (
	configFileName = ".secrets.yml"
	verbose        = false
)

func saidumlo(configFileName string) SaiDumLo {
	saidumlo := SaiDumLo{}
	saidumlo.parseConfig(configFileName)
	return saidumlo
}

// TODO: check if file exists
// TODO: better error messaging
func processSecretGroups(method string, secretGroups []string) {

	saidumlo := saidumlo(configFileName)

	if len(secretGroups) == 0 {
		secretGroups = make([]string, len(saidumlo.Config.SecretGroups))

		i := 0
		for k := range saidumlo.Config.SecretGroups {
			secretGroups[i] = k
			i++
		}
	}

	for _, k := range secretGroups {
		if secretGroup, ok := saidumlo.Config.SecretGroups[k]; !ok {
			logError("Secret group '%v' is not declared - skipping it", k)
			continue
		} else {
			if method == "encrypt" {
				saidumlo.encryptSecretGroup(secretGroup)
			} else {
				saidumlo.decryptSecretGroup(secretGroup)
			}
		}
	}

	saidumlo.cleanCache()
}
