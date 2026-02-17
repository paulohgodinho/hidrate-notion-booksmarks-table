package notion

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// UpdatePageContentWithJSON replaces all content in a Notion page with a code block containing the provided JSON string.
// This erases all existing content before adding the new code block.
func UpdatePageContentWithJSON(ctx context.Context, apiKey, pageID, jsonContent string) error {
	client := &http.Client{}

	// Step 1: Get existing children
	getURL := fmt.Sprintf("https://api.notion.com/v1/blocks/%s/children", pageID)
	req, err := http.NewRequestWithContext(ctx, "GET", getURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create get children request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Notion-Version", "2022-06-28")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to get children: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("get children failed with status %d", resp.StatusCode)
	}

	var childrenResp struct {
		Results []struct {
			ID string `json:"id"`
		} `json:"results"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&childrenResp); err != nil {
		return fmt.Errorf("failed to decode children response: %w", err)
	}

	// Step 2: Delete all existing children
	for _, child := range childrenResp.Results {
		deleteURL := fmt.Sprintf("https://api.notion.com/v1/blocks/%s", child.ID)
		req, err := http.NewRequestWithContext(ctx, "DELETE", deleteURL, nil)
		if err != nil {
			return fmt.Errorf("failed to create delete request for %s: %w", child.ID, err)
		}
		req.Header.Set("Authorization", "Bearer "+apiKey)
		req.Header.Set("Notion-Version", "2022-06-28")

		resp, err := client.Do(req)
		if err != nil {
			return fmt.Errorf("failed to delete block %s: %w", child.ID, err)
		}
		resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("delete block %s failed with status %d", child.ID, resp.StatusCode)
		}
	}

	// Step 3: Append the new code block
	appendURL := fmt.Sprintf("https://api.notion.com/v1/blocks/%s/children", pageID)
	block := map[string]interface{}{
		"type": "code",
		"code": map[string]interface{}{
			"rich_text": []map[string]interface{}{
				{
					"type": "text",
					"text": map[string]interface{}{
						"content": jsonContent,
					},
				},
			},
			"language": "json",
		},
	}
	body := map[string]interface{}{
		"children": []map[string]interface{}{block},
	}
	bodyJSON, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("failed to marshal append body: %w", err)
	}

	req, err = http.NewRequestWithContext(ctx, "PATCH", appendURL, bytes.NewBuffer(bodyJSON))
	if err != nil {
		return fmt.Errorf("failed to create append request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Notion-Version", "2022-06-28")
	req.Header.Set("Content-Type", "application/json")

	resp, err = client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to append block: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("append block failed with status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}
