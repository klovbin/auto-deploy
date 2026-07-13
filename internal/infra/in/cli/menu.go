package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/trust/deploy/internal/domain/config"
)

func ClearScreen() {
	fmt.Print("\033[H\033[2J")
}

func PrintHeader(cfg config.Config) {
	fmt.Println("Trust Deploy")
	fmt.Println()
	if cfg.Repository != "" {
		fmt.Printf("Текущий репозиторий: %s\n", cfg.Repository)
	} else {
		fmt.Println("Текущий репозиторий: не задан")
	}
	fmt.Println()
}

func Select(items []string) (int, error) {
	prompt := promptui.Select{
		Label: "Выберите пункт",
		Items: items,
	}

	index, _, err := prompt.Run()
	return index, err
}

func PromptRepositoryURL() (string, error) {
	prompt := promptui.Prompt{
		Label: "Ссылка на репозиторий",
		Validate: func(input string) error {
			input = strings.TrimSpace(input)
			if input == "" {
				return fmt.Errorf("ссылка не может быть пустой")
			}
			return nil
		},
	}

	url, err := prompt.Run()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(url), nil
}

func WaitForEnter() {
	fmt.Println()
	fmt.Print("Нажмите Enter для возврата в меню...")
	if _, err := fmt.Scanln(); err != nil {
		os.Exit(0)
	}
}
