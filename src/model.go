package main

// SaiDumLo root structure
type SaiDumLo struct {
	ConfigFileName string
	ConfigDir      string
	Config         Config
}

// Config ...
type Config struct {
	Groups map[string]SecretGroup `yaml:"groups"`
	Vault  Vault
}

// Vault ...
type Vault struct {
	Address   string
	Bin       string
	VaultAuth VaultAuth `yaml:"auth"`
}

// VaultAuth ...
type VaultAuth struct {
	Method             string
	CredentialFilePath string `yaml:"credential_file"`
}

// SecretGroup ...
type SecretGroup struct {
	Mappings []SecretMapping `yaml:"secrets"`
}

// SecretMapping ...
type SecretMapping struct {
	Local string `yaml:"local"`
	Vault string `yaml:"vault"`
}
