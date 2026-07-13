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

	for {
		cli.ClearScreen()
		cli.PrintHeader(cfg)

		index, err := cli.Select([]string{
			"1. Добавить репозиторий",
			"2. Сгенерировать деплой ключи",
			"3. Склонировать репозиторий",
			"4. Сгенерировать CI/CD",
			"5. Сгенерировать ключ для CI/CD",
		})
		if err != nil {
			if cli.IsExit(err) {
				fmt.Fprintln(os.Stderr, "выход")
				return
			}
			fmt.Fprintf(os.Stderr, "ошибка: %v\n", err)
			os.Exit(1)
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
		default:
			fmt.Printf("Выбрано: %d\n", index+1)
			return
		}
	}
}
