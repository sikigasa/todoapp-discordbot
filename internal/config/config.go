package config

import (
	"fmt"
	"os"
)

// Config はアプリケーション設定を保持する
type Config struct {
	// DiscordToken はDiscord Botのトークン
	DiscordToken string
	// APIBaseURL はgithub-task-controllerのベースURL
	APIBaseURL string
	// AuthCookie は認証用セッションCookie値
	AuthCookie string
	// DefaultUserID はプロジェクト操作に使用するデフォルトユーザーID
	DefaultUserID string
}

// Load は環境変数から設定を読み込む
func Load() (*Config, error) {
	token := os.Getenv("DISCORD_TOKEN")
	if token == "" {
		return nil, fmt.Errorf("DISCORD_TOKEN is required")
	}

	apiBaseURL := os.Getenv("API_BASE_URL")
	if apiBaseURL == "" {
		apiBaseURL = "http://localhost:8080"
	}

	authCookie := os.Getenv("AUTH_COOKIE")
	if authCookie == "" {
		return nil, fmt.Errorf("AUTH_COOKIE is required")
	}

	defaultUserID := os.Getenv("DEFAULT_USER_ID")

	return &Config{
		DiscordToken:  token,
		APIBaseURL:    apiBaseURL,
		AuthCookie:    authCookie,
		DefaultUserID: defaultUserID,
	}, nil
}
