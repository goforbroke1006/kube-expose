package domain

type ComposeConfig struct {
	Services map[string]Service `yaml:"services"`
}

type Service struct {
	Ports []string `yaml:"ports"`
}
