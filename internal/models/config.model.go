package models

type AppConfig struct {
	Version  int             `yaml:"version"`
	Projects []ProjectConfig `yaml:"projects"`
}

type ProjectConfig struct {
	Name   string       `yaml:"name"`
	Docker DockerConfig `yaml:"docker"`
	Server ServerConfig `yaml:"server"`
}

type DockerConfig struct {
	ImageName  string            `yaml:"image_name"`
	ImageTag   string            `yaml:"image_tag"`
	Dockerfile string            `yaml:"docker_file"`
	BuildArgs  map[string]string `yaml:"build_args"`
}

type ServerConfig struct {
	Host              string `yaml:"host"`
	Port              int    `yaml:"port"`
	User              string `yaml:"user"`
	KeyPath           string `yaml:"key_path"`
	DockerComposePath string `yaml:"docker_compose_path"`
}
