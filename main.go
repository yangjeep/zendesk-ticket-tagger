package main

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/yangjeep/zendesk-ticket-tagger/config"
	"github.com/yangjeep/zendesk-ticket-tagger/server"
	"github.com/yangjeep/zendesk-ticket-tagger/zendesk"
)

func main() {
	// Initialize logger
	logger := logrus.New()
	// Configure logger settings if needed
	logger.SetLevel(logrus.DebugLevel) // Set to debug level to see the webhook details

	// Initialize the global logger in zendesk package
	zendesk.InitLogger(logger)

	cfg := config.Load()

	// Register webhook before starting server
	webhookURL := fmt.Sprintf("http://%s:%d%s", cfg.WebhookHost, cfg.WebhookPort, cfg.WebhookEndpoint)
	if err := zendesk.RegisterWebhook(
		cfg,
		webhookURL,
		"Ticket Tagger Webhook",
		cfg.ZendeskToken,
	); err != nil {
		logger.Fatalf("Failed to register webhook: %v", err)
	}
	logger.Infof("Successfully registered webhook at %s", webhookURL)

	// Start the server
	if err := server.Start(cfg); err != nil {
		logger.Fatalf("Failed to start server: %v", err)
	}
}
