package main

// SaiDumLo root structure
type SaiDumLo struct {
	ConfigFileName string
	ConfigDir      string
	Config         Config
}

// Config ...
type Config struct {
	Groups       map[string]SecretGroup `yaml:"groups"`
	VaultAddress string                 `yaml:"vault_address"`
	VaultBin     string                 `yaml:"vault_bin"`
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
