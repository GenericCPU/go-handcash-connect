package handcash

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// Item represents an item retrieved from HandCash
type Item struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Price string `json:"price"`
	// Add more fields as needed
}

// GetItems retrieves the items for the associated auth token
func (c *Client) GetItems(ctx context.Context, token string) ([]Item, error) {
	// Make sure we have an auth token
	if len(token) == 0 {
		return nil, fmt.Errorf("missing auth token")
	}

	// Get the signed request
	signed, err := c.getSignedRequest(
		http.MethodGet,
		endpointGetItems,
		token,
		&requestBody{authToken: token},
		currentISOTimestamp(),
	)
	if err != nil {
		return nil, fmt.Errorf("error creating signed request: %w", err)
	}

	// Make the HTTP request
	response := httpRequest(
		ctx,
		c,
		&httpPayload{
			Data:           []byte(emptyBody),
			ExpectedStatus: http.StatusOK,
			Method:         signed.Method,
			URL:            signed.URI,
		},
		signed,
	)

	// Error in request?
	if response.Error != nil {
		return nil, response.Error
	}

	// Unmarshal into a slice of items
	var items []Item
	if err = json.Unmarshal(response.BodyContents, &items); err != nil {
		return nil, fmt.Errorf("failed to unmarshal: %w", err)
	}
	return items, nil
}
