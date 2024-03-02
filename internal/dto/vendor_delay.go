package dto

type VendorDelay struct {
	VendorID   uint   `json:"vendor_id"`
	VendorName string `json:"vendor_name"`
	TotalDelay int    `json:"total_delay"`
}
