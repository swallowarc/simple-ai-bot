package e2e

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"
)

const channelSecret = "xxxxx"
const localWebhookURL = "http://localhost:8080/callback"

func TestChat(t *testing.T) {
	t.Skip("Skip e2e test(a mock server for LINE needs to be implemented)")

	webhookRequest := `{
		"events": [
			{
				"type": "message",
				"replyToken": "testReplyToken",
				"source": {
					"userId": "testUserId",
					"type": "user"
				},
				"timestamp": 1462629479859,
				"message": {
					"type": "text",
					"id": "testMessageId",
					"text": "ドラえもんの誕生日は？"
				}
			}
		]
	}`

	signature := generateSignature(channelSecret, webhookRequest)

	req, err := http.NewRequest("POST", localWebhookURL, bytes.NewBufferString(webhookRequest))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Line-Signature", signature)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	fmt.Println("Response:", resp.Status, string(body))
}

func generateSignature(channelSecret, requestBody string) string {
	mac := hmac.New(sha256.New, []byte(channelSecret))
	mac.Write([]byte(requestBody))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}
