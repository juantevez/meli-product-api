package dto

import (
	"meli-product-api/internal/domain/model"
	"time"
)

type ProductDetailsResponse struct {
	Product         ProductDTO          `json:"product"`
	Seller          SellerDTO           `json:"seller"`
	Shipping        ShippingDTO         `json:"shipping"`
	Reviews         ReviewsDTO          `json:"reviews"`
	Questions       []QuestionDTO       `json:"questions"`
	RelatedProducts []RelatedProductDTO `json:"related_products"`
}

type ProductDTO struct {
	ID                string         `json:"id"`
	Title             string         `json:"title"`
	Description       string         `json:"description"`
	Price             float64        `json:"price"`
	OriginalPrice     *float64       `json:"original_price,omitempty"`
	DiscountPercent   *int           `json:"discount_percentage,omitempty"`
	Condition         string         `json:"condition"`
	AvailableQuantity int            `json:"available_quantity"`
	SoldQuantity      int            `json:"sold_quantity"`
	Images            []string       `json:"images"`
	Category          string         `json:"category"`
	Attributes        []AttributeDTO `json:"attributes"`
	Brand             string         `json:"brand"`
	Model             string         `json:"model"`
}

type AttributeDTO struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type SellerDTO struct {
	ID              string  `json:"id"`
	Nickname        string  `json:"nickname"`
	ReputationLevel string  `json:"reputation_level"`
	TotalSales      int     `json:"total_sales"`
	ReputationScore float64 `json:"reputation_score"`
	YearsActive     int     `json:"years_active"`
	IsOfficialStore bool    `json:"is_official_store"`
}

type ShippingDTO struct {
	FreeShipping      bool    `json:"free_shipping"`
	ShippingMode      string  `json:"shipping_mode"`
	Cost              float64 `json:"cost"`
	EstimatedDelivery string  `json:"estimated_delivery"`
	FullFulfillment   bool    `json:"full_fulfillment"`
	PickupAvailable   string  `json:"pickup_available"`
}

type ReviewsDTO struct {
	AverageRating float64     `json:"average_rating"`
	TotalReviews  int         `json:"total_reviews"`
	Items         []ReviewDTO `json:"items"`
}

type ReviewDTO struct {
	ID           string    `json:"id"`
	UserName     string    `json:"user_name"`
	Rating       int       `json:"rating"`
	Title        string    `json:"title"`
	Comment      string    `json:"comment"`
	CreatedAt    time.Time `json:"created_at"`
	HelpfulCount int       `json:"helpful_count"`
}

type QuestionDTO struct {
	ID           string     `json:"id"`
	UserName     string     `json:"user_name"`
	Question     string     `json:"question"`
	Answer       string     `json:"answer"`
	QuestionDate time.Time  `json:"question_date"`
	AnswerDate   *time.Time `json:"answer_date,omitempty"`
	Likes        int        `json:"likes"`
}

type RelatedProductDTO struct {
	ID           string  `json:"id"`
	Title        string  `json:"title"`
	Price        float64 `json:"price"`
	Image        string  `json:"image,omitempty"`
	SoldQuantity int     `json:"sold_quantity"`
}

// Mapper functions
func ToProductDetailsResponse(details *model.ProductDetails) *ProductDetailsResponse {
	return &ProductDetailsResponse{
		Product:         toProductDTO(details.Product),
		Seller:          toSellerDTO(details.Seller),
		Shipping:        toShippingDTO(details.Shipping),
		Reviews:         toReviewsDTO(details.Reviews, details.AverageRating, details.TotalReviews),
		Questions:       toQuestionDTOs(details.Questions),
		RelatedProducts: toRelatedProductDTOs(details.RelatedProducts),
	}
}

func toProductDTO(p model.Product) ProductDTO {
	attrs := make([]AttributeDTO, len(p.Attributes))
	for i, attr := range p.Attributes {
		attrs[i] = AttributeDTO{
			Name:  attr.Name,
			Value: attr.Value,
		}
	}

	return ProductDTO{
		ID:                p.ID,
		Title:             p.Title,
		Description:       p.Description,
		Price:             p.Price,
		OriginalPrice:     p.OriginalPrice,
		DiscountPercent:   p.DiscountPercent,
		Condition:         p.Condition,
		AvailableQuantity: p.AvailableQuantity,
		SoldQuantity:      p.SoldQuantity,
		Images:            p.Images,
		Category:          p.Category,
		Attributes:        attrs,
		Brand:             p.Brand,
		Model:             p.Model,
	}
}

func toSellerDTO(s model.Seller) SellerDTO {
	return SellerDTO{
		ID:              s.ID,
		Nickname:        s.Nickname,
		ReputationLevel: s.ReputationLevel,
		TotalSales:      s.TotalSales,
		ReputationScore: s.ReputationScore,
		YearsActive:     s.YearsActive,
		IsOfficialStore: s.IsOfficialStore,
	}
}

func toShippingDTO(s model.Shipping) ShippingDTO {
	return ShippingDTO{
		FreeShipping:      s.FreeShipping,
		ShippingMode:      s.ShippingMode,
		Cost:              s.Cost,
		EstimatedDelivery: s.EstimatedDelivery,
		FullFulfillment:   s.FullFulfillment,
		PickupAvailable:   s.PickupAvailable,
	}
}

func toReviewsDTO(reviews []model.Review, avgRating float64, total int) ReviewsDTO {
	items := make([]ReviewDTO, len(reviews))
	for i, r := range reviews {
		items[i] = ReviewDTO{
			ID:           r.ID,
			UserName:     r.UserName,
			Rating:       r.Rating,
			Title:        r.Title,
			Comment:      r.Comment,
			CreatedAt:    r.CreatedAt,
			HelpfulCount: r.HelpfulCount,
		}
	}

	return ReviewsDTO{
		AverageRating: avgRating,
		TotalReviews:  total,
		Items:         items,
	}
}

func toQuestionDTOs(questions []model.Question) []QuestionDTO {
	dtos := make([]QuestionDTO, len(questions))
	for i, q := range questions {
		dtos[i] = QuestionDTO{
			ID:           q.ID,
			UserName:     q.UserName,
			Question:     q.Question,
			Answer:       q.Answer,
			QuestionDate: q.QuestionDate,
			AnswerDate:   q.AnswerDate,
			Likes:        q.Likes,
		}
	}
	return dtos
}

func toRelatedProductDTOs(products []model.Product) []RelatedProductDTO {
	dtos := make([]RelatedProductDTO, len(products))
	for i, p := range products {
		image := ""
		if len(p.Images) > 0 {
			image = p.Images[0]
		}

		dtos[i] = RelatedProductDTO{
			ID:           p.ID,
			Title:        p.Title,
			Price:        p.Price,
			Image:        image,
			SoldQuantity: p.SoldQuantity,
		}
	}
	return dtos
}
