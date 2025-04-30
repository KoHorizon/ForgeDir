package generator

import "github.com/KoHorizon/ForgeDir/internal/config"

type BoilerplateGenerator interface {
	Generate(cfg *config.Config, projectRoot string) error
	GetLanguage() string
}
