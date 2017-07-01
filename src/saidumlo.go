package main

import (
	"bufio"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
)

func (saidumlo *SaiDumLo) getConfigDir(configFilePath string) {
	var dirPrefix = "./"
	var result string
	for {
		if _, err := os.Stat(dirPrefix + configFilePath); os.IsNotExist(err) {
			dirPrefix += "../"
			var curDir, _ = filepath.Abs(filepath.Dir(dirPrefix))
			if curDir == "/" {
				logError("Could not find config file '%v'", configFilePath)
				os.Exit(1)
			}
		} else {
			result, err = filepath.Abs(filepath.Dir(dirPrefix + configFilePath))
			checkErr(err)
			break
		}
	}

	saidumlo.ConfigDir = result + "/"
}

// TODO: make less confusing -> proper var names like relativeConfigFilePath..
func (saidumlo *SaiDumLo) parseConfig(configFile string) {
	saidumlo.getConfigDir(configFile)
	saidumlo.ConfigFileName = filepath.Base(configFile)
	saidumlo.Config = Config{}
	configFilePath := saidumlo.ConfigDir + saidumlo.ConfigFileName
	logInfo("Using config %v", configFilePath)
	s, e := ioutil.ReadFile(configFilePath)
	checkErr(e)
	e = yaml.Unmarshal(s, &saidumlo.Config)
	checkErr(e)

	logDebug("%#v", saidumlo.Config)
}

func saidumlo(configFile string) SaiDumLo {
	saidumlo := SaiDumLo{}
	saidumlo.parseConfig(configFile)
	return saidumlo
}

func (saidumlo *SaiDumLo) readSecretFromVault(secretMapping SecretMapping) {
	logInfo("%s read -field=value %s > %s/%s", saidumlo.Config.Vault.Bin, secretMapping.Vault, saidumlo.ConfigDir, secretMapping.Local)

	// TODO: This always overwrites existing file. File should still exist if vault error occurs
	outfile, fileErr := os.Create(fmt.Sprintf("%s/%s", saidumlo.ConfigDir, secretMapping.Local))
	checkErr(fileErr)
	defer outfile.Close()

	env := os.Environ()
	env = append(env, fmt.Sprintf("VAULT_ADDR=%s", saidumlo.Config.Vault.Address))

	cmd := exec.Command(saidumlo.Config.Vault.Bin, "read", "-field=value", secretMapping.Vault)
	cmd.Env = env
	cmd.Dir = saidumlo.ConfigDir
	cmd.Stdin = os.Stdin
	cmd.Stdout = outfile
	cmd.Stderr = os.Stderr
	commandErr := cmd.Run()
	checkErr(commandErr)
}

func (saidumlo *SaiDumLo) writeSecretToVault(secretMapping SecretMapping) {
	logInfo("%s write %s value=@%s/%s", saidumlo.Config.Vault.Bin, secretMapping.Vault, saidumlo.ConfigDir, secretMapping.Local)

	env := os.Environ()
	env = append(env, fmt.Sprintf("VAULT_ADDR=%s", saidumlo.Config.Vault.Address))

	cmd := exec.Command(saidumlo.Config.Vault.Bin, "write", secretMapping.Vault, fmt.Sprintf("value=@%s/%s", saidumlo.ConfigDir, secretMapping.Local))
	cmd.Env = env
	cmd.Dir = saidumlo.ConfigDir
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	commandErr := cmd.Run()
	checkErr(commandErr)
}

func (saidumlo *SaiDumLo) authenticate() {
	if saidumlo.Config.Vault.VaultAuth.Method != "" {

		// read credentials from credential file
		args := []string{"auth", fmt.Sprintf("-method=%s", saidumlo.Config.Vault.VaultAuth.Method)}
		if file, err := os.Open(fmt.Sprintf("%s/%s", saidumlo.ConfigDir, saidumlo.Config.Vault.VaultAuth.CredentialFilePath)); err == nil {

			defer file.Close()

			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				args = append(args, scanner.Text())
			}

			if err = scanner.Err(); err != nil {
				logFatal("%v", err)
			}

		} else {
			logFatal("%v", err)
		}

		// set environment
		env := os.Environ()
		env = append(env, fmt.Sprintf("VAULT_ADDR=%s", saidumlo.Config.Vault.Address))

		// execute command
		cmd := exec.Command(saidumlo.Config.Vault.Bin, args...)
		cmd.Env = env
		cmd.Dir = saidumlo.ConfigDir
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		commandErr := cmd.Run()
		checkErr(commandErr)
	} else {
		logInfo("No auth method defined - skip auth")
	}
}

func processSecretGroups(method string, secretGroups []string) {
	sdl := saidumlo(configFile)
	logDebug("%+v\n", sdl)

	sdl.authenticate()

	var groupsToProcess = secretGroups
	if len(secretGroups) == 0 {
		groupsToProcess = getMapKeys(sdl.Config.Groups)
	}

	for _, secretGroupName := range groupsToProcess {
		for _, secretMapping := range sdl.Config.Groups[secretGroupName].Mappings {
			if method == writeOperationID {
				sdl.writeSecretToVault(secretMapping)
			} else if method == readOperationID {
				sdl.readSecretFromVault(secretMapping)
			} else {
				logError("Unknown operation %s", method)
			}
		}
	}
}
