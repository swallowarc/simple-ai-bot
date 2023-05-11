package repositories

import (
	"testing"
)

func TestJson(t *testing.T) {
	type (
		ChatMessage struct {
			ID   string `json:"id"`
			Text string `json:"text"`
		}
		ChatMessages []ChatMessage
	)

	cms := ChatMessages{
		ChatMessage{
			ID:   "1",
			Text: "Hello",
		},
		ChatMessage{
			ID:   "2",
			Text: "World",
		},
	}

	j, err := toJson[ChatMessages](cms)
	if err != nil {
		t.Fatalf("failed to marshal json: %v", err)
	}
	// fmt.Printf("%+v\n", j)

	_, err = fromJson[ChatMessages](j)
	if err != nil {
		t.Fatalf("failed to unmarshal json: %v", err)
	}
	// fmt.Printf("%+v\n", cms2)
}
