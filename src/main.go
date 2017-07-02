package main

import (
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
	"strconv"
)

var (
	// Ops Ids
	readOperationID  = "read"
	writeOperationID = "write"

	// cmd config values
	configFile = ".secrets.yml"
	verbose    = false
	vaultID    = ""
)

// CommandWithGroups represents a command which addresses specific secret groups.
type CommandWithSecretGroups struct {
	SecretGroups []string
	VaultID      string
}

func (arg *CommandWithSecretGroups) read(c *kingpin.ParseContext) error {
	arg.VaultID = vaultID
	arg.processCommandWithSecretGroups(readOperationID)
	return nil
}

func (arg *CommandWithSecretGroups) write(c *kingpin.ParseContext) error {
	arg.VaultID = vaultID
	arg.processCommandWithSecretGroups(writeOperationID)
	return nil
}

func configureCommands(app *kingpin.Application) {
	// read
	readArg := &CommandWithSecretGroups{}
	readCmd := app.Command("read", "Read secret groups.").Action(readArg.read)
	readCmd.Arg("secretGroup", "The secret groups to read.").StringsVar(&readArg.SecretGroups)

	// write
	writeArg := &CommandWithSecretGroups{}
	writeCmd := app.Command("write", "Write secret groups.").Action(writeArg.write)
	writeCmd.Arg("secretGroup", "The secret groups to write.").StringsVar(&writeArg.SecretGroups)

	// flags
	app.Flag("vault", "Vault").Short('b').StringVar(&vaultID)
	app.Flag("config-file", "Config file").Default(configFile).Short('f').StringVar(&configFile)
	app.Flag("verbose", "Verbose").Default(strconv.FormatBool(verbose)).Short('v').BoolVar(&verbose)
}

func main() {
	app := kingpin.New("saidumlo", "SaiDumLo client secret management tool.")
	configureCommands(app)
	kingpin.MustParse(app.Parse(os.Args[1:]))
}
