package models

// Revert represents a reversal of a claim.
type Revert struct {
	ID        string `json:"id" db:"id"`               // Unique ID of the reversal (UUID)
	ClaimID   string `json:"claim_id" db:"claim_id"`   // ID of the claim that was reverted
	Timestamp string `json:"timestamp" db:"timestamp"` // Date and time of the reversal
}
