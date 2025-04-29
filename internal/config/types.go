package config

type StructureNode struct {
	Type     string          `yaml:"type"`
	Name     string          `yaml:"name"`
	Children []StructureNode `yaml:"children, omitempty"`
}

type Config struct {
	ProjectName string          `yaml:"projectName"`
	Structure   []StructureNode `yaml:"structure"`
}
