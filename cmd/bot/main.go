package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/sikigasa/todoapp-discordbot/internal/api"
	"github.com/sikigasa/todoapp-discordbot/internal/bot"
	"github.com/sikigasa/todoapp-discordbot/internal/config"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	// .envファイルを読み込む（存在しなくてもエラーにしない）
	if err := godotenv.Load(); err != nil {
		logger.Info("No .env file found, using environment variables")
	}

	cfg, err := config.Load()
	if err != nil {
		logger.Error("Failed to load config", "error", err)
		os.Exit(1)
	}

	apiClient := api.NewClient(cfg.APIBaseURL, cfg.AuthCookie)

	b, err := bot.New(cfg.DiscordToken, apiClient, cfg.DefaultUserID, logger)
	if err != nil {
		logger.Error("Failed to create bot", "error", err)
		os.Exit(1)
	}

	if err := b.Start(); err != nil {
		logger.Error("Failed to start bot", "error", err)
		os.Exit(1)
	}

	logger.Info("Bot is running. Press Ctrl+C to exit.")

	// シグナル待機
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	if err := b.Stop(); err != nil {
		logger.Error("Failed to stop bot", "error", err)
	}

	logger.Info("Bot stopped")
}
