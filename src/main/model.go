package main

// SaiDumLo root structure
type SaiDumLo struct {
	ConfigFileName string
	ConfigDir      string
	Config         Config
}

// Config ...
type Config struct {
	Vaults       map[string]Vault       `yaml:"vaults"`
	SecretGroups map[string]SecretGroup `yaml:"secrets"`
}

// Vault ...
type Vault struct {
	Default   bool      `yaml:"default"`
	Address   string    `yaml:"address"`
	Bin       string    `yaml:"bin"`
	VaultAuth VaultAuth `yaml:"auth"`
}

// VaultAuth ...
type VaultAuth struct {
	Method             string `yaml:"method"`
	CredentialFilePath string `yaml:"credential_file"`
}

// SecretGroup ...
type SecretGroup struct {
	Mappings []SecretMapping `yaml:"mappings"`
	LeaseTTL string          `yaml:"lease_ttl"`
}

// SecretMapping ...
type SecretMapping struct {
	Local string `yaml:"local"`
	Vault string `yaml:"vault"`
}
