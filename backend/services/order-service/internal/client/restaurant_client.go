package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"github.com/google/uuid"
)

type MenuItem struct {
	ID           uuid.UUID `json:"id"`
	RestaurantID uuid.UUID `json:"restaurant_id"`
	Name         string    `json:"name"`
	Price        float64   `json:"price"`
	IsAvailable  bool      `json:"is_available"`
}

type MenuResponse struct {
	Success bool `json:"success"`
	Data    struct {
		Categories []struct {
			Items []MenuItem `json:"items"`
		} `json:"categories"`
		Uncategorized []MenuItem `json:"uncategorized"`
	} `json:"data"`
}

type RestaurantClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewRestaurantClient(baseURL string) *RestaurantClient {
	return &RestaurantClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

// FetchMenu gets the entire menu of a restaurant to validate items and grab price snapshots
func (c *RestaurantClient) FetchMenu(restaurantID uuid.UUID) (map[uuid.UUID]MenuItem, error) {
	url := fmt.Sprintf("%s/api/restaurants/%s/menu", c.baseURL, restaurantID.String())
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to call restaurant service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("restaurant service returned status: %d", resp.StatusCode)
	}

	var menuResp MenuResponse
	if err := json.NewDecoder(resp.Body).Decode(&menuResp); err != nil {
		return nil, err
	}

	itemMap := make(map[uuid.UUID]MenuItem)
	
	// Map items inside categories
	for _, cat := range menuResp.Data.Categories {
		for _, item := range cat.Items {
			itemMap[item.ID] = item
		}
	}
	// Map uncategorized items
	for _, item := range menuResp.Data.Uncategorized {
		itemMap[item.ID] = item
	}

	return itemMap, nil
}
