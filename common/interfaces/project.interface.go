package interfaces

import "deploy-kit/internal/models"

type ProjectUsecase interface {
	RunProject(projectConfig models.ProjectConfig)
}
