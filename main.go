package main

import (
	"fmt"
	"os"

	"github.com/trust/deploy/internal/app/commands"
	"github.com/trust/deploy/internal/infra/in/cli"
	configstore "github.com/trust/deploy/internal/infra/out/config"
)

func main() {
	cfg, err := configstore.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "ошибка чтения конфига: %v\n", err)
		os.Exit(1)
	}

	addRepo := commands.AddRepositoryHandler{}
	generateDeployKeys := commands.GenerateDeployKeysHandler{}
	cloneRepo := commands.CloneRepositoryHandler{}
	generateGitLabCI := commands.GenerateGitLabCIHandler{}
	generateCICDKey := commands.GenerateCICDKeyHandler{}
	installDocker := commands.InstallDockerHandler{}
	runService := commands.RunServiceHandler{}

	menuItems := []cli.Item{
		{Label: "1. Добавить репозиторий"},
		{Label: "2. Сгенерировать деплой ключи"},
		{Label: "3. Склонировать репозиторий"},
		{Label: "4. Сгенерировать CI/CD", Disabled: commands.GitLabCIExists(cfg), DisabledText: "уже есть"},
		{Label: "5. Сгенерировать ключ для CI/CD"},
		{Label: "6. Установить Docker"},
		{Label: "7. Запустить сервис с билдом", Disabled: !commands.CanRunService(cfg)},
	}

	for {
		cli.ClearScreen()
		cli.PrintHeader(cfg)

		index, err := cli.Select(menuItems)
		if err != nil {
			if cli.IsExit(err) {
				fmt.Fprintln(os.Stderr, "выход")
				return
			}
			fmt.Fprintf(os.Stderr, "ошибка: %v\n", err)
			os.Exit(1)
		}

		if menuItems[index].Disabled {
			fmt.Println("Пункт недоступен")
			cli.WaitForEnter()
			continue
		}

		switch index {
		case 0:
			url, err := cli.PromptRepositoryURL()
			if err != nil {
				if cli.IsExit(err) {
					continue
				}
				fmt.Fprintf(os.Stderr, "ошибка: %v\n", err)
				os.Exit(1)
			}
			if err := addRepo.Handle(&cfg, url); err != nil {
				fmt.Fprintf(os.Stderr, "ошибка: %v\n", err)
				os.Exit(1)
			}
		case 1:
			if err := generateDeployKeys.Handle(cfg); err != nil {
				fmt.Fprintf(os.Stderr, "ошибка: %v\n", err)
				os.Exit(1)
			}
			cli.WaitForEnter()
		case 2:
			if err := cloneRepo.Handle(cfg); err != nil {
				fmt.Fprintf(os.Stderr, "ошибка: %v\n", err)
				os.Exit(1)
			}
			cli.WaitForEnter()
		case 3:
			if err := generateGitLabCI.Handle(cfg); err != nil {
				fmt.Fprintf(os.Stderr, "ошибка: %v\n", err)
				os.Exit(1)
			}
			cli.WaitForEnter()
		case 4:
			if err := generateCICDKey.Handle(cfg); err != nil {
				fmt.Fprintf(os.Stderr, "ошибка: %v\n", err)
				os.Exit(1)
			}
			cli.WaitForEnter()
		case 5:
			if err := installDocker.Handle(); err != nil {
				fmt.Fprintf(os.Stderr, "ошибка: %v\n", err)
				os.Exit(1)
			}
			cli.WaitForEnter()
		case 6:
			if err := runService.Handle(cfg); err != nil {
				fmt.Fprintf(os.Stderr, "ошибка: %v\n", err)
				os.Exit(1)
			}
			cli.WaitForEnter()
		default:
			fmt.Printf("Выбрано: %d\n", index+1)
			return
		}
	}
}
