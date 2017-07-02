package main

// SaiDumLo root structure
type SaiDumLo struct {
	ConfigFileName string
	ConfigDir      string
	Config         Config
}

// Config ...
type Config struct {
	Vaults map[string]Vault `yaml:"vaults"`
}

// Vault ...
type Vault struct {
	Default      bool
	Address      string
	Bin          string
	VaultAuth    VaultAuth              `yaml:"auth"`
	SecretGroups map[string]SecretGroup `yaml:"secrets"`
}

// VaultAuth ...
type VaultAuth struct {
	Method             string
	CredentialFilePath string `yaml:"credential_file"`
}

// SecretGroup ...
type SecretGroup struct {
	Mappings []SecretMapping
}

// SecretMapping ...
type SecretMapping struct {
	Local string `yaml:"local"`
	Vault string `yaml:"vault"`
}

//
type SecretBackup struct {
	Local string `yaml:"local"`
	Vault string `yaml:"vault"`
}
