package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/manifoldco/promptui"
)

const configFileName = ".trust-deploy.json"

type Config struct {
	Repository string `json:"repository"`
}

func configPath() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return filepath.Join(cwd, configFileName), nil
}

func loadConfig() (Config, error) {
	path, err := configPath()
	if err != nil {
		return Config{}, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return Config{}, nil
		}
		return Config{}, err
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func saveConfig(cfg Config) error {
	path, err := configPath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0o644)
}

func clearScreen() {
	fmt.Print("\033[H\033[2J")
}

func printHeader(cfg Config) {
	fmt.Println("Trust Deploy")
	fmt.Println()
	if cfg.Repository != "" {
		fmt.Printf("Текущий репозиторий: %s\n", cfg.Repository)
	} else {
		fmt.Println("Текущий репозиторий: не задан")
	}
	fmt.Println()
}

func addRepository(cfg *Config) error {
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
		return err
	}

	cfg.Repository = strings.TrimSpace(url)
	return saveConfig(*cfg)
}

func waitForEnter() {
	fmt.Println()
	fmt.Print("Нажмите Enter для возврата в меню...")
	if _, err := fmt.Scanln(); err != nil {
		os.Exit(0)
	}
}

func main() {
	cfg, err := loadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "ошибка чтения конфига: %v\n", err)
		os.Exit(1)
	}

	for {
		clearScreen()
		printHeader(cfg)

		items := []string{
			"1. Добавить репозиторий",
			"2. Сгенерировать деплой ключи",
			"3. Склонировать репозиторий",
			"4. Сгенерировать CI/CD",
			"5. Сгенерировать ключ для CI/CD",
		}

		prompt := promptui.Select{
			Label: "Выберите пункт",
			Items: items,
		}

		index, _, err := prompt.Run()
		if err != nil {
			fmt.Fprintf(os.Stderr, "ошибка: %v\n", err)
			os.Exit(1)
		}

		switch index {
		case 0:
			if err := addRepository(&cfg); err != nil {
				fmt.Fprintf(os.Stderr, "ошибка: %v\n", err)
				os.Exit(1)
			}
		case 1:
			if err := generateDeployKeys(cfg); err != nil {
				fmt.Fprintf(os.Stderr, "ошибка: %v\n", err)
				os.Exit(1)
			}
			waitForEnter()
		case 2:
			if err := cloneRepository(cfg); err != nil {
				fmt.Fprintf(os.Stderr, "ошибка: %v\n", err)
				os.Exit(1)
			}
			waitForEnter()
		case 3:
			if err := generateGitLabCI(cfg); err != nil {
				fmt.Fprintf(os.Stderr, "ошибка: %v\n", err)
				os.Exit(1)
			}
			waitForEnter()
		case 4:
			if err := generateCICDKey(cfg); err != nil {
				fmt.Fprintf(os.Stderr, "ошибка: %v\n", err)
				os.Exit(1)
			}
			waitForEnter()
		default:
			fmt.Printf("Выбрано: %d\n", index+1)
			return
		}
	}
}
