package config

type Node struct {
	Type     string `yaml:"type"`
	Name     string `yaml:"name"`
	Children []Node `yaml:"children, omitempty"`
}

type Config struct {
	ProjectName string `yaml:"projectName"`
	Structure   []Node `yaml:"structure"`
}
