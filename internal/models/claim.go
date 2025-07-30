package models

// Claim represents a medication claim.
type Claim struct {
	ID        string  `json:"id" db:"id"`               // Unique ID of the claim (UUID)
	NDC       string  `json:"ndc" db:"ndc"`             // National Drug Code of the medication
	NPI       string  `json:"npi" db:"npi"`             // National Provider Identifier of the pharmacy
	Quantity  float64 `json:"quantity" db:"quantity"`   // Quantity of the medication
	Price     float64 `json:"price" db:"price"`         // Price of the medication
	Timestamp string  `json:"timestamp" db:"timestamp"` // Date and time of claim submission
	Reverted  bool    `json:"reverted" db:"reverted"`   // Indicates if the claim has been reverted
}

// ClaimSubmissionRequest represents the input payload for creating a new claim.
type ClaimSubmissionRequest struct {
	NDC      string  `json:"ndc"`      // National Drug Code of the medication
	Quantity float64 `json:"quantity"` // Quantity of the medication
	NPI      string  `json:"npi"`      // National Provider Identifier of the pharmacy
	Price    float64 `json:"price"`    // Price of the medication
}

// ClaimResponse represents the response payload after a claim submission.
type ClaimResponse struct {
	Status  string `json:"status"`   // Operation status (e.g., "claim submitted")
	ClaimID string `json:"claim_id"` // ID of the created claim
}

// ClaimReversalRequest represents the input payload for reverting a claim.
type ClaimReversalRequest struct {
	ClaimID string `json:"claim_id"` // ID of the claim to be reverted
}

// ClaimReversalResponse represents the response payload after a claim reversal.
type ClaimReversalResponse struct {
	Status  string `json:"status"`   // Operation status (e.g., "claim reversed")
	ClaimID string `json:"claim_id"` // ID of the reverted claim
}
