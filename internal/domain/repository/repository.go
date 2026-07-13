package repository

import (
	"fmt"
	"strings"
)

func NameFromURL(raw string) (string, error) {
	raw = strings.TrimSpace(raw)
	raw = strings.TrimSuffix(raw, ".git")
	raw = strings.TrimSuffix(raw, "/")

	if raw == "" {
		return "", fmt.Errorf("ссылка на репозиторий пустая")
	}

	if strings.HasPrefix(raw, "git@") {
		if idx := strings.LastIndex(raw, ":"); idx >= 0 {
			raw = raw[idx+1:]
		}
	} else {
		raw = strings.TrimPrefix(raw, "https://")
		raw = strings.TrimPrefix(raw, "http://")
		raw = strings.TrimPrefix(raw, "ssh://")

		if idx := strings.Index(raw, "/"); idx >= 0 {
			raw = raw[idx+1:]
		}
	}

	parts := strings.Split(raw, "/")
	name := parts[len(parts)-1]
	if name == "" {
		return "", fmt.Errorf("не удалось определить имя репозитория")
	}

	return name, nil
}
