package dto

import "meli-product-api/internal/domain/model"

type ProductSearchResponse struct {
	Query        string              `json:"query"`
	TotalResults int                 `json:"total_results"`
	Limit        int                 `json:"limit"`
	Offset       int                 `json:"offset"`
	Results      []ProductSummaryDTO `json:"results"`
}

type ProductSummaryDTO struct {
	ID                string   `json:"id"`
	Title             string   `json:"title"`
	Price             float64  `json:"price"`
	OriginalPrice     *float64 `json:"original_price,omitempty"`
	DiscountPercent   *int     `json:"discount_percentage,omitempty"`
	Condition         string   `json:"condition"`
	Thumbnail         string   `json:"thumbnail,omitempty"`
	SoldQuantity      int      `json:"sold_quantity"`
	AvailableQuantity int      `json:"available_quantity"`
	Category          string   `json:"category"`
	Brand             string   `json:"brand"`
	FreeShipping      bool     `json:"free_shipping"`
}

func ToProductSearchResponse(query string, products []model.Product, total, limit, offset int) *ProductSearchResponse {
	summaries := make([]ProductSummaryDTO, len(products))

	for i, p := range products {
		thumbnail := ""
		if len(p.Images) > 0 {
			thumbnail = p.Images[0]
		}

		freeShipping := p.Price > 50000

		summaries[i] = ProductSummaryDTO{
			ID:                p.ID,
			Title:             p.Title,
			Price:             p.Price,
			OriginalPrice:     p.OriginalPrice,
			DiscountPercent:   p.DiscountPercent,
			Condition:         p.Condition,
			Thumbnail:         thumbnail,
			SoldQuantity:      p.SoldQuantity,
			AvailableQuantity: p.AvailableQuantity,
			Category:          p.Category,
			Brand:             p.Brand,
			FreeShipping:      freeShipping,
		}
	}

	return &ProductSearchResponse{
		Query:        query,
		TotalResults: total,
		Limit:        limit,
		Offset:       offset,
		Results:      summaries,
	}
}
