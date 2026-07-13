package cli

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

func openTTY() (*os.File, error) {
	return os.Open("/dev/tty")
}

func selectSimple(items []string, in io.Reader) (int, error) {
	fmt.Println("Выберите пункт:")
	for i, item := range items {
		fmt.Printf("  %d) %s\n", i+1, item)
	}
	fmt.Print("> ")

	scanner := bufio.NewScanner(in)
	if !scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return 0, err
		}
		return 0, io.EOF
	}

	n, err := strconv.Atoi(strings.TrimSpace(scanner.Text()))
	if err != nil || n < 1 || n > len(items) {
		return 0, fmt.Errorf("неверный выбор")
	}

	return n - 1, nil
}

func promptSimple(label string, in io.Reader) (string, error) {
	fmt.Printf("%s: ", label)
	scanner := bufio.NewScanner(in)
	if !scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return "", err
		}
		return "", io.EOF
	}

	value := strings.TrimSpace(scanner.Text())
	if value == "" {
		return "", fmt.Errorf("значение не может быть пустым")
	}

	return value, nil
}
