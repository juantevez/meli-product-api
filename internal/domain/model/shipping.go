package model

type Shipping struct {
	FreeShipping      bool    `json:"free_shipping"`
	ShippingMode      string  `json:"shipping_mode"`
	Cost              float64 `json:"cost"`
	EstimatedDelivery string  `json:"estimated_delivery"`
	FullFulfillment   bool    `json:"full_fulfillment"`
	PickupAvailable   string  `json:"pickup_available"`
}
