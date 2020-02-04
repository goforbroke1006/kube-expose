package domain

type ComposeConfig struct {
	Namespaces []Namespace `yaml:"namespaces"`
}

type Namespace struct {
	Name      string              `yaml:"name"`
	Resources map[string]Resource `yaml:"resources"`
}

type Resource struct {
	Type  string   `yaml:"type"`
	Ports []string `yaml:"ports"`
}
