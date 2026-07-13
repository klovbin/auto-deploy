package cli

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/chzyer/readline"
	"github.com/manifoldco/promptui"
	"github.com/trust/deploy/internal/domain/config"
)

func IsExit(err error) bool {
	if err == nil {
		return false
	}
	if errors.Is(err, io.EOF) || errors.Is(err, readline.ErrInterrupt) {
		return true
	}
	switch err.Error() {
	case "^D", "^C":
		return true
	}
	return false
}

func ClearScreen() {
	fmt.Print("\033[H\033[2J")
}

func PrintHeader(cfg config.Config) {
	fmt.Println("Trust Deploy")
	fmt.Println()
	fmt.Printf("Рабочая папка: %s\n", cfg.WorkDirectory())
	if cfg.Repository != "" {
		fmt.Printf("Текущий репозиторий: %s\n", cfg.Repository)
	} else {
		fmt.Println("Текущий репозиторий: не задан")
	}
	fmt.Println()
}

func Select(items []Item) (int, error) {
	labels := itemLabels(items)

	tty, err := openTTY()
	if err != nil {
		return selectSimple(labels, os.Stdin)
	}
	defer tty.Close()

	prompt := promptui.Select{
		Label: "Выберите пункт",
		Items: labels,
		Size:  len(labels),
		Stdin: tty,
	}

	index, _, err := prompt.Run()
	if err == nil {
		return index, nil
	}

	if index, simpleErr := selectSimple(labels, tty); simpleErr == nil {
		return index, nil
	}
	if IsExit(err) {
		return 0, err
	}

	return 0, err
}

func PromptRepositoryURL() (string, error) {
	tty, err := openTTY()
	if err != nil {
		return promptSimple("Ссылка на репозиторий", os.Stdin)
	}
	defer tty.Close()

	prompt := promptui.Prompt{
		Label: "Ссылка на репозиторий",
		Stdin: tty,
		Validate: func(input string) error {
			if strings.TrimSpace(input) == "" {
				return fmt.Errorf("ссылка не может быть пустой")
			}
			return nil
		},
	}

	url, err := prompt.Run()
	if err == nil {
		return strings.TrimSpace(url), nil
	}
	if IsExit(err) {
		return "", err
	}

	return promptSimple("Ссылка на репозиторий", tty)
}

func WaitForEnter() {
	fmt.Println()
	fmt.Print("Нажмите Enter для возврата в меню...")

	tty, err := openTTY()
	if err != nil {
		_, _ = fmt.Scanln()
		return
	}
	defer tty.Close()

	scanner := bufio.NewScanner(tty)
	scanner.Scan()
}
