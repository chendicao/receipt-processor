package models

// Item represents a receipt item
type Item struct {
	ID               int     `json:"id" db:"id"`                                                  // Mapping the ID field in the Item table
	ReceiptID        string  `json:"receiptId" db:"receipt_id"`                                   // Foreign key linking to the Receipt table
	ShortDescription string  `json:"shortDescription" validate:"required" db:"short_description"` // Item's short description
	Price            float64 `json:"price" validate:"required" db:"price"`                        // Item price
}

// Receipt represents the receipt information
type Receipt struct {
	ID           string  `json:"id" db:"id"`                                                              // Mapping the ID field in the Receipt table
	Retailer     string  `json:"retailer" validate:"required" db:"retailer"`                              // Retailer's name
	PurchaseDate string  `json:"purchaseDate" validate:"required,datetime=2006-01-02" db:"purchase_date"` // Date of purchase
	PurchaseTime string  `json:"purchaseTime" validate:"required,datetime=15:04" db:"purchase_time"`      // Time of purchase
	Total        float64 `json:"total" validate:"required" db:"total"`                                    // Total receipt amount
	Items        []Item  `json:"items" validate:"required,dive,required" db:"-"`                          // List of items (not stored in the receipt table directly)
}
