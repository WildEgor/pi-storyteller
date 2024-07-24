//go:build wireinject

package app

import (
	"github.com/google/wire"
)

// ServerSet
var ServerSet = wire.NewSet(Set)

// NewServer
func NewServer() (*App, error) {
	wire.Build(ServerSet)
	return nil, nil
}
