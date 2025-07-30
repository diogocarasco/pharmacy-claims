package models

// Pharmacy represents a pharmacy
type Pharmacy struct {
	Chain string `json:"chain" db:"chain"` // Name of the pharmacy chain (e.g., health, saint, doctor)
	NPI   string `json:"npi" db:"npi"`     // National Provider Identifier of the pharmacy
}
