package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type QueueInfo struct {
	Arguments  map[string]any `json:"arguments"`
	AutoDelete bool           `json:"auto_delete"`
	Durable    bool           `json:"durable"`
	Exclusive  bool           `json:"exclusive"`
	Name       string         `json:"name"`
	Node       string         `json:"node"`
	State      string         `json:"state"`
	Type       string         `json:"type"`
	VHost      string         `json:"vhost"`
}

func main() {
	queueName := "app.queue.a"
	vhost := "goapp-vhost"

	username := "goadmin"
	password := "123456"

	request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://localhost:15672/api/queues/%s", vhost), nil)
	if err != nil {
		log.Fatal("Failed to new request")
	}

	authorizationValue := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", username, password)))
	request.Header.Set("Authorization", fmt.Sprintf("Basic %s", authorizationValue))
	client := &http.Client{}
	response, err := client.Do(request)
	fmt.Printf("HTTP status: %d\n", response.StatusCode)

	body, err := io.ReadAll(response.Body)
	queues := make([]QueueInfo, 0)

	err = json.Unmarshal(body, &queues)
	if err != nil {
		log.Fatal(err)
	}

	for _, d := range queues {
		if d.Name == queueName {
			fmt.Println("exists")
			os.Exit(0)
		}
	}
	fmt.Println("not exist")
}
