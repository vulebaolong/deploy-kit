package usecase

import (
	"deploy-kit/common/assets"
	"deploy-kit/common/interfaces"
	"deploy-kit/internal/models"
	"errors"
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type ConfigUsecase struct {
	ConfigPath string
	Config     *models.AppConfig
}

func NewConfigUsecase() interfaces.ConfigUsecase {
	log.Println("ConfigUsecase Init called")

	const configFileName = "deploy-kit.yaml"

	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get working directory: %v", err)
	}
	wdConfigPath := filepath.Join(wd, configFileName)

	exePath, err := os.Executable()
	if err != nil {
		log.Fatalf("Failed to get executable path: %v", err)
	}
	exeDir := filepath.Dir(exePath)
	exeConfigPath := filepath.Join(exeDir, configFileName)

	log.Printf("Working directory config path: %s", wdConfigPath)
	log.Printf("Executable directory config path: %s", exeConfigPath)

	var configPath string

	// Ưu tiên config ở working directory
	if fileExists(wdConfigPath) {
		configPath = wdConfigPath
		log.Printf("Using config file from working directory: %s", configPath)
	} else if fileExists(exeConfigPath) {
		// fallback sang config cạnh binary
		configPath = exeConfigPath
		log.Printf("Using config file next to executable: %s", configPath)
	} else {
		// nếu cả hai nơi đều chưa có thì tạo mới cạnh binary
		configPath = exeConfigPath

		err = os.WriteFile(configPath, assets.ConfigExample, 0644)
		if err != nil {
			log.Fatalf("Failed to create config file: %v", err)
		}

		log.Printf("Config file created successfully at: %s", configPath)
	}

	// đọc file config
	raw, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatalf("Failed to read config file: %v", err)
	}

	var cfg models.AppConfig
	err = yaml.Unmarshal(raw, &cfg)
	if err != nil {
		log.Fatalf("Failed to parse yaml config: %v", err)
	}

	return &ConfigUsecase{
		ConfigPath: configPath,
		Config:     &cfg,
	}
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	if err == nil {
		return !info.IsDir()
	}

	if errors.Is(err, os.ErrNotExist) {
		return false
	}

	log.Fatalf("Failed to check file %s: %v", path, err)
	return false
}

// GetConfig implements [interfaces.ConfigUsecase].
func (c *ConfigUsecase) GetConfig() *models.AppConfig {
	return c.Config
}

// GetConfigPath implements [interfaces.ConfigUsecase].
func (c *ConfigUsecase) GetConfigPath() string {
	return c.ConfigPath
}
