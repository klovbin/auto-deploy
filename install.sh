#!/usr/bin/env bash
set -euo pipefail

REPO_URL="https://github.com/klovbin/auto-deploy.git"
INSTALL_DIR="${INSTALL_DIR:-$HOME/auto-deploy}"

install_go() {
	if command -v go >/dev/null 2>&1; then
		return 0
	fi

	echo "Go не найден, устанавливаю..."

	if command -v apt-get >/dev/null 2>&1; then
		sudo apt-get update -qq
		sudo apt-get install -y golang-go git
	elif command -v apk >/dev/null 2>&1; then
		sudo apk add --no-cache go git
	elif command -v dnf >/dev/null 2>&1; then
		sudo dnf install -y golang git
	elif command -v yum >/dev/null 2>&1; then
		sudo yum install -y golang git
	elif command -v brew >/dev/null 2>&1; then
		brew install go git
	else
		echo "Не удалось установить Go автоматически. Установите Go 1.22+ вручную: https://go.dev/dl/"
		exit 1
	fi
}

install_git() {
	if command -v git >/dev/null 2>&1; then
		return 0
	fi

	echo "Git не найден, устанавливаю..."

	if command -v apt-get >/dev/null 2>&1; then
		sudo apt-get update -qq
		sudo apt-get install -y git
	elif command -v apk >/dev/null 2>&1; then
		sudo apk add --no-cache git
	elif command -v dnf >/dev/null 2>&1; then
		sudo dnf install -y git
	elif command -v yum >/dev/null 2>&1; then
		sudo yum install -y git
	elif command -v brew >/dev/null 2>&1; then
		brew install git
	else
		echo "Не удалось установить Git автоматически."
		exit 1
	fi
}

install_git
install_go

if [ -d "$INSTALL_DIR/.git" ]; then
	echo "Обновляю $INSTALL_DIR..."
	git -C "$INSTALL_DIR" pull --ff-only
else
	echo "Клонирую в $INSTALL_DIR..."
	git clone "$REPO_URL" "$INSTALL_DIR"
fi

cd "$INSTALL_DIR"
go build -o deploy .

echo "Запуск..."
exec ./deploy
