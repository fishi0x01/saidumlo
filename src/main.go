package main

import (
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
	"strconv"
)

var (
	// Ops Ids
	readOperationId  = "read"
	writeOperationId = "write"

	// cmd config values
	configFile = ".secrets.yml"
	verbose    = false
)

// A CommandWithGroups represents a command which addresses specific secret groups.
type CommandWithGroups struct {
	secretGroups []string
}

func (arg *CommandWithGroups) read(c *kingpin.ParseContext) error {
	processSecretGroups(readOperationId, arg.secretGroups)
	return nil
}

func (arg *CommandWithGroups) write(c *kingpin.ParseContext) error {
	processSecretGroups(writeOperationId, arg.secretGroups)
	return nil
}

func configureCommands(app *kingpin.Application) {
	// read
	readArg := &CommandWithGroups{}
	readCmd := app.Command("read", "Read secret groups.").Action(readArg.read)
	readCmd.Arg("secretGroup", "The secret groups to read.").StringsVar(&readArg.secretGroups)

	// write
	writeArg := &CommandWithGroups{}
	writeCmd := app.Command("write", "Write secret groups.").Action(writeArg.write)
	writeCmd.Arg("secretGroup", "The secret groups to write.").StringsVar(&writeArg.secretGroups)

	// flags
	app.Flag("config-file", "Config file").Default(configFile).Short('f').StringVar(&configFile)
	app.Flag("verbose", "Verbose").Default(strconv.FormatBool(verbose)).Short('v').BoolVar(&verbose)
}

func main() {
	app := kingpin.New("saidumlo", "SaiDumLo client secret management tool.")
	configureCommands(app)
	kingpin.MustParse(app.Parse(os.Args[1:]))
}
