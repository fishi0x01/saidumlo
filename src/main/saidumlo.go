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

var (
	configDir = ""
)

/*****************************
 * Config management
 *****************************/
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

func (saidumlo *SaiDumLo) getDefaultVault() Vault {
	var vault Vault
	var vaultID string
	var alreadyFoundDefault bool

	for k, v := range saidumlo.Config.Vaults {
		if v.Default {
			if alreadyFoundDefault {
				logError("Multiple vaults set as default ('%s' and '%s'), but only one allowed.", vaultID, k)
			} else {
				vault = v
				vaultID = k
				alreadyFoundDefault = true
			}
		}
	}

	return vault
}

func saidumlo(configFile string) SaiDumLo {
	saidumlo := SaiDumLo{}
	saidumlo.parseConfig(configFile)
	return saidumlo
}

/***************************
 * Vault interaction
 ***************************/
func (vault *Vault) readSecretMapping(secretMapping SecretMapping) {
	logInfo("%s read -field=value %s > %s/%s", vault.Bin, secretMapping.Vault, configDir, secretMapping.Local)

	// TODO: This always overwrites existing file. File should still exist if vault error occurs
	outfile, fileErr := os.Create(fmt.Sprintf("%s/%s", configDir, secretMapping.Local))
	checkErr(fileErr)
	defer outfile.Close()

	env := os.Environ()
	env = append(env, fmt.Sprintf("VAULT_ADDR=%s", vault.Address))

	cmd := exec.Command(vault.Bin, "read", "-field=value", secretMapping.Vault)
	cmd.Env = env
	cmd.Dir = configDir
	cmd.Stdin = os.Stdin
	cmd.Stdout = outfile
	cmd.Stderr = os.Stderr
	commandErr := cmd.Run()
	checkErr(commandErr)
}

func (vault *Vault) writeSecretMapping(secretMapping SecretMapping, leaseTtl string) {
	if leaseTtl != "" {
		leaseTtl = fmt.Sprintf("ttl=%s", leaseTtl)
	}

	logInfo("%s write %s value=@%s/%s %s", vault.Bin, secretMapping.Vault, configDir, secretMapping.Local, leaseTtl)

	env := os.Environ()
	env = append(env, fmt.Sprintf("VAULT_ADDR=%s", vault.Address))

	cmd := exec.Command(vault.Bin, "write", secretMapping.Vault, fmt.Sprintf("value=@%s/%s", configDir, secretMapping.Local), leaseTtl)
	cmd.Env = env
	cmd.Dir = configDir
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	commandErr := cmd.Run()
	checkErr(commandErr)
}

func (vault *Vault) authenticate() {
	if vault.VaultAuth.Method != "" {

		// read credentials from credential file
		args := []string{"auth", fmt.Sprintf("-method=%s", vault.VaultAuth.Method)}
		if file, err := os.Open(fmt.Sprintf("%s/%s", configDir, vault.VaultAuth.CredentialFilePath)); err == nil {

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
		env = append(env, fmt.Sprintf("VAULT_ADDR=%s", vault.Address))

		// execute command
		cmd := exec.Command(vault.Bin, args...)
		cmd.Env = env
		cmd.Dir = configDir
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		commandErr := cmd.Run()
		checkErr(commandErr)
	} else {
		logInfo("No auth method defined - skip auth")
	}
}

/*****************************
 * Command processors
 *****************************/
func (cwsg *CommandWithSecretGroups) processCommandWithSecretGroups(method string) {
	sdl := saidumlo(configFile)
	vault := sdl.getDefaultVault()
	if cwsg.VaultID != "" {
		vault = sdl.Config.Vaults[cwsg.VaultID]
	}
	configDir = sdl.ConfigDir
	logDebug("%+v\n", sdl)

	vault.authenticate()

	var groupsToProcess = cwsg.SecretGroups
	if len(cwsg.SecretGroups) == 0 {
		groupsToProcess = getMapKeys(sdl.Config.SecretGroups)
	}

	for _, secretGroupName := range groupsToProcess {
		var leaseTtl = sdl.Config.SecretGroups[secretGroupName].LeaseTtl
		for _, secretMapping := range sdl.Config.SecretGroups[secretGroupName].Mappings {
			if method == writeOperationID {
				vault.writeSecretMapping(secretMapping, leaseTtl)
			} else if method == readOperationID {
				vault.readSecretMapping(secretMapping)
			} else {
				logError("Unknown operation %s", method)
			}
		}
	}
}
