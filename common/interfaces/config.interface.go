package interfaces

import "deploy-kit/internal/models"

type ConfigHandler interface{}

type ConfigUsecase interface {
	GetConfigPath() string
	GetConfig() *models.AppConfig
}
