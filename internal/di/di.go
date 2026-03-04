package di

import (
	"deploy-kit/internal/delivery"
	"deploy-kit/internal/usecase"
	"log"
)

func NewApp() *delivery.CLI {
	configUsecase := usecase.NewConfigUsecase()
	configPath := configUsecase.GetConfigPath()
	config := configUsecase.GetConfig()

	log.Printf("Config loaded successfully from: %s", configPath)
	log.Printf("Total projects loaded: %d", len(config.Projects))

	projectUsecase := usecase.NewProjectUsecase()

	cli := delivery.NewCLI(config, projectUsecase)
	return cli
}
