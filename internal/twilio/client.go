package twilio

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Client is a tiny Twilio REST client implemented with net/http.
// This avoids coupling the example to a specific Twilio SDK version.
type Client struct {
	accountSid string
	authToken  string
	httpClient *http.Client
}

func NewClient(accountSid, authToken string) *Client {
	return &Client{
		accountSid: accountSid,
		authToken:  authToken,
		httpClient: &http.Client{Timeout: 15 * time.Second},
	}
}

// SendSMS sends an outbound SMS using Twilio's Messages API.
//
// statusCallbackURL is optional; when provided, Twilio will POST delivery status
// updates to that endpoint.
func (c *Client) SendSMS(ctx context.Context, from, to, body, statusCallbackURL string) (map[string]any, error) {
	endpoint := fmt.Sprintf("https://api.twilio.com/2010-04-01/Accounts/%s/Messages.json", c.accountSid)

	form := url.Values{}
	form.Set("From", from)
	form.Set("To", to)
	form.Set("Body", body)
	if statusCallbackURL != "" {
		form.Set("StatusCallback", statusCallbackURL)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, io.NopCloser(bytes.NewBufferString(form.Encode())))
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}
	req.SetBasicAuth(c.accountSid, c.authToken)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("twilio request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("twilio error: status=%d body=%s", resp.StatusCode, strings.TrimSpace(string(respBody)))
	}

	var out map[string]any
	if err := json.Unmarshal(respBody, &out); err != nil {
		return nil, fmt.Errorf("decode response JSON: %w", err)
	}
	return out, nil
}

