package model

type Seller struct {
	ID              string  `json:"id"`
	Nickname        string  `json:"nickname"`
	ReputationLevel string  `json:"reputation_level"`
	TotalSales      int     `json:"total_sales"`
	ReputationScore float64 `json:"reputation_score"`
	YearsActive     int     `json:"years_active"`
	IsOfficialStore bool    `json:"is_official_store"`
}
