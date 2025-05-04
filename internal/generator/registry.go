package generator

import "maps"

// BoilerplateGenerator is the interface that every language-specific
// generator implementation must satisfy. Using this interface allows the
// registry to accept any supported language instance, enabling easy
// dependency injection into the Coordinator.
var (
	registered = make(map[string]BoilerplateGenerator)
)

// Register registers a new generator under its language key
// and panics if that key is already in use.
func Register(g BoilerplateGenerator) {
	lang := g.GetLanguage()
	if _, exists := registered[lang]; exists {
		panic("generator: duplicate registration for language " + lang)
	}
	registered[lang] = g
}

// All returns a copy of the registry so callers can safely enumerate
// without mutating internal state.
func All() map[string]BoilerplateGenerator {
	copy := make(map[string]BoilerplateGenerator, len(registered))
	maps.Copy(copy, registered)
	return copy
}
