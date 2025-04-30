package config

const (
	TypeDir  = "dir"
	TypeFile = "file"
	// Other type to expand
)

type StructureNode struct {
	Type     string          `yaml:"type"`
	Name     string          `yaml:"name"`
	Children []StructureNode `yaml:"children,omitempty"`
}

type Config struct {
	ProjectName string          `yaml:"projectName"`
	Language    string          `yaml:"language"`
	Structure   []StructureNode `yaml:"structure"`
}
