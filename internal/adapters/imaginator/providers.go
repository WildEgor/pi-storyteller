package imaginator

import (
	"github.com/WildEgor/pi-storyteller/internal/configs"
	"github.com/WildEgor/pi-storyteller/pkg/kandinsky"
	"github.com/WildEgor/pi-storyteller/pkg/kandinsky/mocks"
)

// KandinskyClientProvider wrapper for client
type KandinskyClientProvider struct {
	kandinsky.Client
}

// NewKandinskyClientProvider create client
func NewKandinskyClientProvider(
	config *configs.KandinskyConfig,
	appConfig *configs.AppConfig,
) *KandinskyClientProvider {
	if appConfig.IsDebug() {
		return &KandinskyClientProvider{mocks.NewKandinskyDummyClient(config.Config)}
	}

	return &KandinskyClientProvider{kandinsky.New(config.Config)}
}
