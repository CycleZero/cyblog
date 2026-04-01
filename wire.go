//go:build wireinject
// +build wireinject

package main

import (
	"cyblog/internal"
	"cyblog/pkg/infra"
	"cyblog/pkg/log"
	"cyblog/pkg/repo"

	"github.com/google/wire"
	"github.com/spf13/viper"
)

func initApp(vc *viper.Viper, logger *log.Logger) *MainApp {
	panic(wire.Build(
		NewMainApp,
		infra.ProviderSet,
		internal.ProviderSet,
		repo.ProviderSet,
	))
}
