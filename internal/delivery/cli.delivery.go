package delivery

import (
	"bufio"
	"deploy-kit/common/interfaces"
	"deploy-kit/internal/models"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type CLI struct {
	config         *models.AppConfig
	projectUsecase interfaces.ProjectUsecase
}

func NewCLI(config *models.AppConfig, projectUsecase interfaces.ProjectUsecase) *CLI {
	return &CLI{
		config:         config,
		projectUsecase: projectUsecase,
	}
}

func (c *CLI) Run() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("\n====== Deploy-toolkit ======")
		for i, project := range c.config.Projects {
			fmt.Printf("%d) %s\n", i+1, project.Name)
		}
		fmt.Println("0) Thoát")
		fmt.Println("")
		fmt.Print(">> Chọn: ")

		choice, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("Lỗi đọc input: %v", err)
			continue
		}

		choice = strings.TrimSpace(choice)

		if choice == "0" {
			fmt.Println("Tạm biệt!")
			return
		}

		selectedIndex, err := strconv.Atoi(choice)
		if err != nil {
			fmt.Println("Vui lòng nhập số hợp lệ.")
			continue
		}

		if selectedIndex < 1 || selectedIndex > len(c.config.Projects) {
			fmt.Println("Lựa chọn không hợp lệ, vui lòng thử lại.")
			continue
		}

		projectConfig := c.config.Projects[selectedIndex-1]

		log.Printf("Bạn đã chọn project: %s", projectConfig.Name)

		c.projectUsecase.RunProject(projectConfig)
	}
}
