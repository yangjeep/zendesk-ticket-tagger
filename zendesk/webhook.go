package zendesk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/yangjeep/zendesk-ticket-tagger/config"
)

var log *logrus.Logger

// InitLogger initializes the global logger
func InitLogger(logger *logrus.Logger) {
	log = logger
}

type WebhookRequest struct {
	Webhook struct {
		Name           string   `json:"name"`
		Status         string   `json:"status"`
		Endpoint       string   `json:"endpoint"`
		HTTPMethod     string   `json:"http_method"`
		RequestFormat  string   `json:"request_format"`
		Subscriptions  []string `json:"subscriptions"`
		Authentication struct {
			Type string `json:"type"`
			Data struct {
				Token string `json:"token"`
			} `json:"data"`
			AddPosition string `json:"add_position"`
		} `json:"authentication"`
	} `json:"webhook"`
}

type WebhookResponse struct {
	Webhook struct {
		ID     string `json:"id"`
		Status string `json:"status"`
	} `json:"webhook"`
	Error struct {
		Title   string   `json:"title"`
		Message string   `json:"message"`
		Details []string `json:"details"`
	} `json:"error"`
}

// RegisterWebhook registers a new webhook with Zendesk
func RegisterWebhook(cfg *config.Config, endpoint, webhookName, bearerToken string) error {
	// Prepare webhook request
	webhookReq := WebhookRequest{}
	webhookReq.Webhook.Name = webhookName
	webhookReq.Webhook.Status = "active"

	// Ensure the endpoint uses HTTPS
	if !strings.HasPrefix(endpoint, "https://") {
		return fmt.Errorf("webhook endpoint must use HTTPS: %s", endpoint)
	}
	webhookReq.Webhook.Endpoint = endpoint
	webhookReq.Webhook.HTTPMethod = "POST"
	webhookReq.Webhook.RequestFormat = "json"
	webhookReq.Webhook.Subscriptions = []string{"conditional_ticket_events"}
	webhookReq.Webhook.Authentication.Type = "bearer_token"
	webhookReq.Webhook.Authentication.Data.Token = bearerToken
	webhookReq.Webhook.Authentication.AddPosition = "header"

	// Convert to JSON
	jsonData, err := json.Marshal(webhookReq)
	if err != nil {
		return fmt.Errorf("failed to marshal webhook request: %w", err)
	}

	// Log the webhook request JSON
	log.Debugf("Webhook request: %s", string(jsonData))

	// Create request
	url := fmt.Sprintf("https://%s.zendesk.com/api/v2/webhooks", cfg.ZendeskSubdomain)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Add headers
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(fmt.Sprintf("%s/token", cfg.ZendeskEmail), cfg.ZendeskToken)

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Read and parse response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	// Log both raw response and status code
	log.Debugf("Webhook registration response (status %d):\nHeaders: %v\nBody: %s",
		resp.StatusCode,
		resp.Header,
		string(body))

	var webhookResp WebhookResponse
	if err := json.Unmarshal(body, &webhookResp); err != nil {
		return fmt.Errorf("failed to parse response: %w, body: %s", err, string(body))
	}

	// Check response status
	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		if webhookResp.Error.Message != "" {
			return fmt.Errorf("webhook registration failed: %s - %s (details: %v)",
				webhookResp.Error.Title,
				webhookResp.Error.Message,
				webhookResp.Error.Details)
		}
		return fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}

	if webhookResp.Webhook.Status != "active" {
		return fmt.Errorf("webhook created but not active, status: %s", webhookResp.Webhook.Status)
	}

	log.Infof("Successfully registered webhook (ID: %s) with status: %s",
		webhookResp.Webhook.ID,
		webhookResp.Webhook.Status)
	return nil
}
