package main

/********
* Structures
*********/

// SaiDumLo ...
type SaiDumLo struct {
    ConfigFileName string
    ConfigDir string
    Config Config
}

// Config ...
type Config struct {
    Encryptions map[string]Encryption
    SecretGroups map[string]SecretGroup `yaml:"secretGroups"`
}

// SecretGroup ...
type SecretGroup struct {
    Encryption string
    Key string
    Secrets []Secret
}

// Secret ...
type Secret struct {
    Encrypted string
    Decrypted string
}

// CommandSetting ...
type CommandSetting struct {
    Args []string
    Executable string
    Env []EnvVar
}

// Encryption ...
type Encryption struct {
    Encrypt CommandSetting
    Decrypt CommandSetting
}

// EnvVar ...
type EnvVar struct {
    Name string
    Value string
}

// Mapping ...
type Mapping struct {
    Key string
    In string
    Out string
}
