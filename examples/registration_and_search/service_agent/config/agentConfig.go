package config

type AgentConfig struct {
	AgentName    string      `yaml:"agent_name"`
	AgentURL     string      `yaml:"agent_url"`
	RegistryURL  string      `yaml:"registry_url"`
	KafkaURL     string      `yaml:"kafka_url"`
	UniqueID     string      `yaml:"unique_id"`
	AgentType    string      `yaml:"agent_type"`
	Capabilities []string    `yaml:"capabilities"`
	Metadata     Metadata    `yaml:"metadata"`
	Status       Status      `yaml:"status"`
	AuthData     AuthData    `yaml:"auth_data"`
	ContactInfo  ContactInfo `yaml:"contact_info"`
	Tags         []string    `yaml:"tags"`
}

type Metadata struct {
	SoftwareVersion string   `yaml:"software_version"`
	Location        Location `yaml:"location"`
	Region          string   `yaml:"region"`
	Organization    string   `yaml:"organization"`
}

type Location struct {
	Latitude  float64 `yaml:"latitude"`
	Longitude float64 `yaml:"longitude"`
}

type Status struct {
	CurrentStatus string `yaml:"current_status"`
}

type AuthData struct {
	Token     string `yaml:"token"`
	PublicKey string `yaml:"public_key"`
}

type ContactInfo struct {
	Email            string `yaml:"email"`
	SecondaryChannel string `yaml:"secondary_channel"`
}
