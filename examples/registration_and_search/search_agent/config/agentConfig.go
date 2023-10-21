package config

type AgentConfig struct {
	AgentName   string `yaml:"agent_name"`
	AgentURL    string `yaml:"agent_url"`
	RegistryURL string `yaml:"registry_url"`
	UniqueID    string `yaml:"unique_id"`
	AgentType   string `yaml:"agent_type"`
}
