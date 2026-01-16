package model

import "time"

type Product struct {
	ID                string      `json:"id"`
	Title             string      `json:"title"`
	Description       string      `json:"description"`
	Price             float64     `json:"price"`
	OriginalPrice     *float64    `json:"original_price,omitempty"`
	DiscountPercent   *int        `json:"discount_percentage,omitempty"`
	Condition         string      `json:"condition"`
	AvailableQuantity int         `json:"available_quantity"`
	SoldQuantity      int         `json:"sold_quantity"`
	Images            []string    `json:"images"`
	Category          string      `json:"category"`
	Attributes        []Attribute `json:"attributes"`
	Brand             string      `json:"brand"`
	Model             string      `json:"model"`
	CreatedAt         time.Time   `json:"created_at"`
	UpdatedAt         time.Time   `json:"updated_at"`
}

type Attribute struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type ProductDetails struct {
	Product         Product    `json:"product"`
	Seller          Seller     `json:"seller"`
	Shipping        Shipping   `json:"shipping"`
	Reviews         []Review   `json:"reviews"`
	AverageRating   float64    `json:"average_rating"`
	TotalReviews    int        `json:"total_reviews"`
	Questions       []Question `json:"questions"`
	RelatedProducts []Product  `json:"related_products"`
}
