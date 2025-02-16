package responses

import "avito-shop/internal/models"

type InfoResponse struct {
	Coins       int                    `json:"coins"`
	Inventory   []models.InventoryItem `json:"inventory"`
	CoinHistory models.CoinHistory     `json:"coinHistory"`
}
