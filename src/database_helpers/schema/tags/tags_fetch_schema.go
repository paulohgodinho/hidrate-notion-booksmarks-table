package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	apiKey := os.Getenv("NOTION_API_KEY")
	if apiKey == "" {
		log.Fatal("NOTION_API_KEY environment variable not set")
	}

	databaseID := "2faa8389a62b802a82ece77040166e52"

	req, err := http.NewRequest("GET", fmt.Sprintf("https://api.notion.com/v1/databases/%s", databaseID), nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Notion-Version", "2022-06-28")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("API request failed with status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		log.Fatal(err)
	}

	properties, ok := result["properties"].(map[string]interface{})
	if !ok {
		log.Fatal("Properties not found in response")
	}

	prettyJSON, err := json.MarshalIndent(properties, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(prettyJSON))
}
