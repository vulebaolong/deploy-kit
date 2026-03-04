package usecase

import (
	"deploy-kit/common/interfaces"
	"deploy-kit/common/ui"
	"deploy-kit/internal/models"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
)

type projectUsecase struct {
}

func NewProjectUsecase() interfaces.ProjectUsecase {
	return &projectUsecase{}
}

// RunProject implements [interfaces.ProjectUsecase].
func (p *projectUsecase) RunProject(projectConfig models.ProjectConfig) {
	ui.PrintStruct("Thông tin project đã chọn", projectConfig)

	imageFullName := fmt.Sprintf("%s:%s", projectConfig.Docker.ImageName, projectConfig.Docker.ImageTag)
	buildContext := filepath.Dir(projectConfig.Docker.Dockerfile)

	// 1. Xóa image cũ ở local nếu có
	ui.Step(fmt.Sprintf("Remove old local image: %s", imageFullName))
	if err := runCommand(buildContext, "docker", "rmi", imageFullName); err != nil {
		ui.Warn(fmt.Sprintf("Không xóa được image cũ hoặc image không tồn tại: %v", err))
	}

	// 2. Build image mới ở local
	ui.Step(fmt.Sprintf("Build new local image: %s", imageFullName))
	args := []string{
		"build",
		"--platform", "linux/amd64",
		"-t", imageFullName,
		"-f", projectConfig.Docker.Dockerfile,
	}
	ui.Info(fmt.Sprintf("Number of build args: %d", len(projectConfig.Docker.BuildArgs)))
	if len(projectConfig.Docker.BuildArgs) > 0 {
		for key, value := range projectConfig.Docker.BuildArgs {
			ui.Info(fmt.Sprintf("Build arg: %s=%s", key, value))
			args = append(args, "--build-arg", fmt.Sprintf("%s=%s", key, value))
		}
	}
	args = append(args, ".")
	err := runCommand(buildContext, "docker", args...)
	if err != nil {
		ui.Error(fmt.Sprintf("Build image thất bại: %v", err))
		os.Exit(1)
	}
	ui.Success("Build image thành công")

	// 3. Test connect server
	ui.Step(fmt.Sprintf("Connect to EC2: %s@%s", projectConfig.Server.User, projectConfig.Server.Host))
	err = runCommand(
		"",
		"ssh",
		"-i", projectConfig.Server.KeyPath,
		"-o", "StrictHostKeyChecking=no",
		fmt.Sprintf("%s@%s", projectConfig.Server.User, projectConfig.Server.Host),
		"echo connected successfully",
	)
	if err != nil {
		ui.Error(fmt.Sprintf("Kết nối EC2 thất bại: %v", err))
		os.Exit(1)
	}
	ui.Success("Kết nối EC2 thành công")

	// 4. Stream image từ local sang server, không tạo file tar
	ui.Step(fmt.Sprintf("Deploy image and restart service on server: %s", imageFullName))
	err = deployImageToServer(
		buildContext,
		imageFullName,
		projectConfig.Server.DockerComposePath,
		projectConfig.Server.KeyPath,
		projectConfig.Server.User,
		projectConfig.Server.Host,
	)
	if err != nil {
		ui.Error(fmt.Sprintf("Deploy to server thất bại: %v", err))
		os.Exit(1)
	}
	ui.Success("Deploy to server thành công")

	// 5 . Xóa image cũ ở local nếu có
	ui.Step(fmt.Sprintf("Remove local image: %s", imageFullName))
	if err := runCommand("", "docker", "rmi", imageFullName); err != nil {
		ui.Warn(fmt.Sprintf("Không xóa được image cũ hoặc image không tồn tại: %v", err))
	} else {
		ui.Success("Xóa local image thành công")
	}
}

func runCommand(dir string, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	return cmd.Run()
}

func deployImageToServer(dir, imageFullName, dockerComposePath, keyPath, user, host string) error {
	saveCmd := exec.Command("docker", "save", imageFullName)
	saveCmd.Dir = dir
	saveCmd.Stderr = os.Stderr

	sshCmd := exec.Command(
		"ssh",
		"-i", keyPath,
		"-o", "StrictHostKeyChecking=no",
		fmt.Sprintf("%s@%s", user, host),
		fmt.Sprintf(
			"sudo docker load && sudo docker compose -f %s up -d && sudo docker image prune -f",
			dockerComposePath,
		),
	)
	sshCmd.Stdout = os.Stdout
	sshCmd.Stderr = os.Stderr

	pipeReader, pipeWriter := io.Pipe()
	saveCmd.Stdout = pipeWriter
	sshCmd.Stdin = pipeReader

	if err := sshCmd.Start(); err != nil {
		_ = pipeWriter.Close()
		_ = pipeReader.Close()
		return err
	}

	if err := saveCmd.Start(); err != nil {
		_ = pipeWriter.Close()
		_ = pipeReader.Close()
		return err
	}

	saveErr := saveCmd.Wait()
	_ = pipeWriter.Close()

	sshErr := sshCmd.Wait()
	_ = pipeReader.Close()

	if saveErr != nil {
		return saveErr
	}
	if sshErr != nil {
		return sshErr
	}

	return nil
}
