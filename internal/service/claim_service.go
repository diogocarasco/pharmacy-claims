package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/diogocarasco/go-pharmacy-service/internal/database"
	"github.com/diogocarasco/go-pharmacy-service/internal/logger"
	"github.com/diogocarasco/go-pharmacy-service/internal/models"
)

// ClaimService defines the interface for claim service operations.
// This interface specifies the methods that any claim service implementation must have.
type ClaimService interface {
	SubmitClaim(req models.ClaimSubmissionRequest) (*models.Claim, error)
	ReverseClaim(req models.ClaimReversalRequest) (*models.Revert, error)
	GetClaimByID(id string) (*models.Claim, error)
	// Add other methods that your ClaimService might have in the future here
}

// claimService is the concrete implementation of the ClaimService interface.
// The lowercase 'c' is a convention to differentiate it from the interface of the same name.
type claimService struct {
	logger logger.Logger
	dbRepo database.DBRepository
}

// NewClaimService creates and returns a new instance of the ClaimService interface.
// It returns a POINTER to the concrete 'claimService' struct, which satisfies the interface.
func NewClaimService(log logger.Logger, dbRepo database.DBRepository) ClaimService {
	return &claimService{ // Returns a pointer to the concrete implementation
		logger: log,
		dbRepo: dbRepo,
	}
}

// SubmitClaim processes the submission of a new claim.
// The '*claimService' receiver means this method operates on a pointer to the struct.
func (s *claimService) SubmitClaim(req models.ClaimSubmissionRequest) (*models.Claim, error) {
	if req.NDC == "" || req.NPI == "" || req.Quantity <= 0 || req.Price <= 0 {
		return nil, errors.New("invalid claim data: NDC, NPI, Quantity, and Price are required and must be positive")
	}

	pharmacy, err := s.dbRepo.GetPharmacyByNPI(req.NPI)
	if err != nil {
		s.logger.Error("Error fetching pharmacy with NPI %s: %v", req.NPI, err)
		return nil, fmt.Errorf("internal error processing claim")
	}
	if pharmacy == nil {
		return nil, fmt.Errorf("invalid NPI '%s'", req.NPI)
	}

	newClaim := models.Claim{
		ID:        uuid.New().String(),
		NDC:       req.NDC,
		NPI:       req.NPI,
		Quantity:  req.Quantity,
		Price:     req.Price,
		Timestamp: time.Now().Format("2006-01-02T15:04:05"), // String format for the timestamp
		Reverted:  false,
	}

	if err := s.dbRepo.SaveClaim(newClaim); err != nil {
		s.logger.Error("Error saving new claim %s: %v", newClaim.ID, err)
		return nil, errors.New("internal error saving claim")
	}

	s.logger.Info("Claim %s submitted successfully for NPI %s", newClaim.ID, newClaim.NPI)
	return &newClaim, nil
}

// ReverseClaim processes the reversal of an existing claim.
// The '*claimService' receiver means this method operates on a pointer to the struct.
func (s *claimService) ReverseClaim(req models.ClaimReversalRequest) (*models.Revert, error) {
	if req.ClaimID == "" {
		return nil, errors.New("invalid reversal claim ID")
	}

	claim, err := s.dbRepo.GetClaimByID(req.ClaimID) // This already exists and works in dbRepo
	if err != nil {
		s.logger.Error("Error fetching claim %s for reversal: %v", req.ClaimID, err)
		return nil, errors.New("internal error reverting claim")
	}
	if claim == nil {
		return nil, fmt.Errorf("claim with ID '%s' not found for reversal", req.ClaimID)
	}
	if claim.Reverted {
		return nil, fmt.Errorf("claim with ID '%s' is already reverted", req.ClaimID)
	}

	// Update the claim status as reverted
	if err := s.dbRepo.UpdateClaimRevertedStatus(claim.ID, true); err != nil {
		s.logger.Error("Error updating claim reversal status for claim %s: %v", claim.ID, err)
		return nil, errors.New("internal error reverting claim")
	}

	newRevert := models.Revert{
		ID:        uuid.New().String(),
		ClaimID:   claim.ID,
		Timestamp: time.Now().Format("2006-01-02T15:04:05"), // String format for the timestamp
	}

	if err := s.dbRepo.SaveRevert(newRevert); err != nil {
		s.logger.Error("Error saving reversal record for claim %s: %v", newRevert.ClaimID, err)
		return nil, errors.New("internal error saving claim reversal")
	}

	s.logger.Info("Claim %s reverted successfully. Revert ID: %s", claim.ID, newRevert.ID)
	return &newRevert, nil
}

// GetClaimByID fetches a claim by its ID.
// This method is now part of the concrete implementation and satisfies the interface.
func (s *claimService) GetClaimByID(id string) (*models.Claim, error) {
	// The actual fetch logic should be in your DBRepository
	claim, err := s.dbRepo.GetClaimByID(id)
	if err != nil {
		s.logger.Error("DB error fetching claim %s: %v", id, err)
		return nil, fmt.Errorf("error fetching claim: %w", err)
	}
	return claim, nil
}
