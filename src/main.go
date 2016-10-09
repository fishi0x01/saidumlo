package main

import (
    "gopkg.in/alecthomas/kingpin.v2"
    "os"
)

// A CommandWithGroups represents a command which addresses specific secret groups.
type CommandWithGroups struct {
    secretGroups []string
}

// A CommandWithFile represents a command which addresses a specific file (like 'edit').
type CommandWithFile struct {
    file string
}

func (arg *CommandWithGroups) encrypt(c *kingpin.ParseContext) error {
    processSecretGroups("encrypt", arg.secretGroups)
    return nil
}

func (arg *CommandWithGroups) decrypt(c *kingpin.ParseContext) error {
    processSecretGroups("decrypt", arg.secretGroups)
    return nil
}

func (arg *CommandWithFile) edit(c *kingpin.ParseContext) error {
    saidumlo := saidumlo(configFileName)
    saidumlo.editSecret(arg.file)
    return nil
}

func configureCommands(app *kingpin.Application) {
    // encrypt
    encArg := &CommandWithGroups{}
    encCmd := app.Command("encrypt", "Encrypt secret groups.").Action(encArg.encrypt)
    encCmd.Arg("secretGroup", "The secret groups to encrypt.").StringsVar(&encArg.secretGroups)

    // decrypt
    decArg := &CommandWithGroups{}
    decCmd := app.Command("decrypt", "Decrypt secret groups.").Action(decArg.decrypt)
    decCmd.Arg("secretGroup", "The secret groups to encrypt.").StringsVar(&decArg.secretGroups)

    // edit
    editArg := &CommandWithFile{}
    editCmd := app.Command("edit", "Edit a (encrypted) file and encrypt it afterwards.").Action(editArg.edit)
    editCmd.Arg("file", "The file to edit and encrypt.").StringVar(&editArg.file)

    app.Flag("v", "Verbose").BoolVar(&verbose)
}

func main() {
    app := kingpin.New("saidumlo", "SaiDumLo secret management tool.")
    configureCommands(app)
    kingpin.MustParse(app.Parse(os.Args[1:]))
}
