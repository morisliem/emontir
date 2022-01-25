package fcm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// still missing the fcm key
// and the registred device id

type FCMMessage struct {
	To            string      `json:"to,omitempty"`
	RegisteredIDs []string    `json:"registred_ids,omitempty"`
	Data          interface{} `json:"data,omitempty"`
}

// inside data return
// title, message and redirect link (/api/v1/orders/{order_id})

func SendNotification(ctx context.Context, msg FCMMessage) {
	payload, err := json.Marshal(msg)
	if err != nil {
		log.Printf("failed to marshal json: %s", err.Error())
	}

	req, err := http.NewRequestWithContext(ctx, "POST", "https://fcm.googleapis.com/fcm/send", bytes.NewReader(payload))
	if err != nil {
		log.Printf("failed to create http request: %s", err.Error())
	}

	req.Header.Set("Authorization", "fcm key")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("request to https://fcm.googleapis.com/fcm/send failed: %s", err.Error())
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("failed to parse resp body content: %s", err.Error())
	}

	fmt.Printf("Request returned status %s:\n\tHeader: %s\n\tBody: %s",
		resp.Status,
		resp.Header,
		body,
	)
}
