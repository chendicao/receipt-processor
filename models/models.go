package models

type Item struct {
	ShortDescription string  `json:"shortDescription" validate:"required"`
	Price            float64 `json:"price" validate:"required"`
}

type Receipt struct {
	Retailer     string  `json:"retailer" validate:"required"`
	PurchaseDate string  `json:"purchaseDate" validate:"required,datetime=2006-01-02"`
	PurchaseTime string  `json:"purchaseTime" validate:"required,datetime=15:04"`
	Items        []Item  `json:"items" validate:"required,dive,required"`
	Total        float64 `json:"total" validate:"required"`
}
