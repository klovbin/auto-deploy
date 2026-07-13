package main

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/pem"
	"fmt"
	"os"
	"path/filepath"

	"golang.org/x/crypto/ssh"
)

func generateSSHKeyPair(privateKeyPath, publicKeyPath string) error {
	if _, err := os.Stat(privateKeyPath); err == nil {
		return fmt.Errorf("ключ уже существует: %s", privateKeyPath)
	}
	if _, err := os.Stat(publicKeyPath); err == nil {
		return fmt.Errorf("ключ уже существует: %s", publicKeyPath)
	}

	if err := os.MkdirAll(filepath.Dir(privateKeyPath), 0o755); err != nil {
		return fmt.Errorf("не удалось создать папку: %w", err)
	}

	_, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return fmt.Errorf("не удалось сгенерировать ключи: %w", err)
	}

	signer, err := ssh.NewSignerFromKey(priv)
	if err != nil {
		return fmt.Errorf("не удалось создать signer: %w", err)
	}

	privateKeyPEM, err := ssh.MarshalPrivateKey(priv, "")
	if err != nil {
		return fmt.Errorf("не удалось сериализовать приватный ключ: %w", err)
	}

	if err := os.WriteFile(privateKeyPath, pem.EncodeToMemory(privateKeyPEM), 0o600); err != nil {
		return fmt.Errorf("не удалось записать приватный ключ: %w", err)
	}

	publicKey := ssh.MarshalAuthorizedKey(signer.PublicKey())
	if err := os.WriteFile(publicKeyPath, publicKey, 0o644); err != nil {
		return fmt.Errorf("не удалось записать публичный ключ: %w", err)
	}

	return nil
}
