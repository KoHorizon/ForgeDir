package generator

import "maps"

// BoilerplateGenerator is the interface that every language-specific
// generator implementation must satisfy. Using this interface allows the
// registry to accept any supported language instance, enabling easy
// dependency injection into the Coordinator.
var (
	registered = make(map[string]BoilerplateGenerator)
)

// Register adds a new BoilerplateGenerator to the registry under its
// language key. It panics on duplicate registrations.
func Register(g BoilerplateGenerator) {
	lang := g.GetLanguage()
	if _, exists := registered[lang]; exists {
		panic("generator: duplicate registration for language " + lang)
	}
	registered[lang] = g
}

// All returns a shallow copy of the registry map, so callers can inject
// all registered generators without mutating the internal state.
func All() map[string]BoilerplateGenerator {
	copy := make(map[string]BoilerplateGenerator, len(registered))
	maps.Copy(copy, registered)
	return copy
}
