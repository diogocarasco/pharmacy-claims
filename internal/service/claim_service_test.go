package service_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/diogocarasco/go-pharmacy-service/internal/logger"
	"github.com/diogocarasco/go-pharmacy-service/internal/models"
	"github.com/diogocarasco/go-pharmacy-service/internal/service"
)

type MockDBRepository struct {
	mock.Mock
}

func (m *MockDBRepository) SavePharmacy(pharmacy models.Pharmacy) error {
	args := m.Called(pharmacy)
	return args.Error(0)
}

func (m *MockDBRepository) GetPharmacyByNPI(npi string) (*models.Pharmacy, error) {
	args := m.Called(npi)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Pharmacy), args.Error(1)
}

func (m *MockDBRepository) SaveClaim(claim models.Claim) error {
	args := m.Called(claim)
	return args.Error(0)
}

func (m *MockDBRepository) SaveClaims(claims []models.Claim) error {
	args := m.Called(claims)
	return args.Error(0)
}

func (m *MockDBRepository) GetClaimByID(id string) (*models.Claim, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Claim), args.Error(1)
}

func (m *MockDBRepository) UpdateClaimRevertedStatus(id string, reverted bool) error {
	args := m.Called(id, reverted)
	return args.Error(0)
}

func (m *MockDBRepository) SaveRevert(revert models.Revert) error {
	args := m.Called(revert)
	return args.Error(0)
}

func (m *MockDBRepository) SaveReverts(reverts []models.Revert) error {
	args := m.Called(reverts)
	return args.Error(0)
}

func (m *MockDBRepository) Close() error {
	args := m.Called()
	return args.Error(0)
}

func TestSubmitClaimSuccess(t *testing.T) {
	mockRepo := new(MockDBRepository)
	mockLogger := logger.NewLogger()

	mockRepo.On("GetPharmacyByNPI", "1234567890").Return(&models.Pharmacy{Chain: "health", NPI: "1234567890"}, nil).Once()
	mockRepo.On("SaveClaim", mock.AnythingOfType("models.Claim")).Return(nil).Once()
	mockRepo.On("Close").Return(nil).Maybe()

	claimService := service.NewClaimService(mockLogger, mockRepo)

	req := models.ClaimSubmissionRequest{
		NDC:      "00002323401",
		NPI:      "1234567890",
		Quantity: 10,
		Price:    50.00,
	}

	claim, err := claimService.SubmitClaim(req)

	assert.Nil(t, err, "Expected no error for successful claim submission")
	assert.NotNil(t, claim, "Expected a claim object to be returned")
	assert.Equal(t, req.NDC, claim.NDC, "Claim NDC should match request NDC")
	assert.Equal(t, req.NPI, claim.NPI, "Claim NPI should match request NPI")
	assert.Equal(t, req.Quantity, claim.Quantity, "Claim Quantity should match request Quantity")
	assert.Equal(t, req.Price, claim.Price, "Claim Price should match request Price")
	assert.False(t, claim.Reverted, "Claim should not be reverted initially")
	mockRepo.AssertExpectations(t)
}

func TestSubmitClaimInvalidData(t *testing.T) {
	mockRepo := new(MockDBRepository)
	mockLogger := logger.NewLogger()
	mockRepo.On("Close").Return(nil).Maybe()

	claimService := service.NewClaimService(mockLogger, mockRepo)

	req := models.ClaimSubmissionRequest{
		NDC:      "",
		NPI:      "1234567890",
		Quantity: 10,
		Price:    50.00,
	}

	claim, err := claimService.SubmitClaim(req)

	assert.Nil(t, claim, "Expected no claim to be returned for invalid data")
	assert.NotNil(t, err, "Expected an error for invalid data")
	assert.Contains(t, err.Error(), "invalid claim data", "Error message should indicate invalid claim data")
	mockRepo.AssertNotCalled(t, "GetPharmacyByNPI", mock.Anything)
	mockRepo.AssertNotCalled(t, "SaveClaim", mock.Anything)
	mockRepo.AssertExpectations(t)
}

func TestSubmitClaimNPINotSupported(t *testing.T) {
	mockRepo := new(MockDBRepository)
	mockLogger := logger.NewLogger()

	mockRepo.On("GetPharmacyByNPI", "9999999999").Return(nil, nil).Once()
	mockRepo.On("Close").Return(nil).Maybe()

	claimService := service.NewClaimService(mockLogger, mockRepo)

	req := models.ClaimSubmissionRequest{
		NDC:      "00002323401",
		NPI:      "9999999999",
		Quantity: 10,
		Price:    50.00,
	}

	claim, err := claimService.SubmitClaim(req)

	assert.Nil(t, claim, "Expected no claim to be returned for unsupported NPI")
	assert.NotNil(t, err, "Expected an error for unsupported NPI")
	assert.Contains(t, err.Error(), "invalid NPI '9999999999'", "Error message should indicate invalid NPI")
	mockRepo.AssertNotCalled(t, "SaveClaim", mock.Anything)
	mockRepo.AssertExpectations(t)
}

func TestReverseClaimSuccess(t *testing.T) {
	mockRepo := new(MockDBRepository)
	mockLogger := logger.NewLogger()

	claimID := "some-valid-claim-id"
	mockClaim := &models.Claim{
		ID:        claimID,
		NDC:       "some-ndc",
		NPI:       "some-npi",
		Quantity:  10,
		Price:     100.0,
		Timestamp: "2006-01-02T15:04:05",
		Reverted:  false,
	}

	mockRepo.On("GetClaimByID", claimID).Return(mockClaim, nil).Once()
	mockRepo.On("UpdateClaimRevertedStatus", claimID, true).Return(nil).Once()
	mockRepo.On("SaveRevert", mock.AnythingOfType("models.Revert")).Return(nil).Once()
	mockRepo.On("Close").Return(nil).Maybe()

	claimService := service.NewClaimService(mockLogger, mockRepo)

	req := models.ClaimReversalRequest{ClaimID: claimID}
	revert, err := claimService.ReverseClaim(req)

	assert.Nil(t, err, "Expected no error for successful claim reversal")
	assert.NotNil(t, revert, "Expected a revert object to be returned")
	assert.Equal(t, claimID, revert.ClaimID, "Revert ClaimID should match request ClaimID")
	mockRepo.AssertExpectations(t)
}

func TestReverseClaimNotFound(t *testing.T) {
	mockRepo := new(MockDBRepository)
	mockLogger := logger.NewLogger()

	claimID := "non-existent-id"

	mockRepo.On("GetClaimByID", claimID).Return(nil, nil).Once()
	mockRepo.On("Close").Return(nil).Maybe()

	claimService := service.NewClaimService(mockLogger, mockRepo)

	req := models.ClaimReversalRequest{ClaimID: claimID}
	revert, err := claimService.ReverseClaim(req)

	assert.Nil(t, revert, "Expected no revert object to be returned when claim is not found")
	assert.NotNil(t, err, "Expected an error when claim is not found")
	assert.Contains(t, err.Error(), "claim with ID 'non-existent-id' not found for reversal", "Error message should indicate claim not found")
	mockRepo.AssertNotCalled(t, "UpdateClaimRevertedStatus", mock.Anything, mock.Anything)
	mockRepo.AssertNotCalled(t, "SaveRevert", mock.Anything)
	mockRepo.AssertExpectations(t)
}

func TestReverseClaimAlreadyReverted(t *testing.T) {
	mockRepo := new(MockDBRepository)
	mockLogger := logger.NewLogger()

	claimID := "already-reverted-id"
	mockClaim := &models.Claim{
		ID:        claimID,
		NDC:       "some-ndc",
		NPI:       "some-npi",
		Quantity:  10,
		Price:     100.0,
		Timestamp: "2006-01-02T15:04:05",
		Reverted:  true,
	}

	mockRepo.On("GetClaimByID", claimID).Return(mockClaim, nil).Once()
	mockRepo.On("Close").Return(nil).Maybe()

	claimService := service.NewClaimService(mockLogger, mockRepo)

	req := models.ClaimReversalRequest{ClaimID: claimID}
	revert, err := claimService.ReverseClaim(req)

	assert.Nil(t, revert, "Expected no revert object to be returned when claim is already reverted")
	assert.NotNil(t, err, "Expected an error when claim is already reverted")
	assert.Contains(t, err.Error(), "claim with ID 'already-reverted-id' is already reverted", "Error message should indicate claim is already reverted")
	mockRepo.AssertNotCalled(t, "UpdateClaimRevertedStatus", mock.Anything, mock.Anything)
	mockRepo.AssertNotCalled(t, "SaveRevert", mock.Anything)
	mockRepo.AssertExpectations(t)
}

func TestReverseClaimUpdateDBError(t *testing.T) {
	mockRepo := new(MockDBRepository)
	mockLogger := logger.NewLogger()

	claimID := "claim-id-for-db-error"
	mockClaim := &models.Claim{
		ID:        claimID,
		NDC:       "some-ndc",
		NPI:       "some-npi",
		Quantity:  10,
		Price:     100.0,
		Timestamp: "2006-01-02T15:04:05",
		Reverted:  false,
	}
	dbError := errors.New("simulated DB error during update")

	mockRepo.On("GetClaimByID", claimID).Return(mockClaim, nil).Once()
	mockRepo.On("UpdateClaimRevertedStatus", claimID, true).Return(dbError).Once()
	mockRepo.On("Close").Return(nil).Maybe()

	claimService := service.NewClaimService(mockLogger, mockRepo)

	req := models.ClaimReversalRequest{ClaimID: claimID}
	revert, err := claimService.ReverseClaim(req)

	assert.Nil(t, revert, "Expected no revert object to be returned on DB update error")
	assert.NotNil(t, err, "Expected an error on DB update error")
	assert.Contains(t, err.Error(), "internal error reverting claim", "Error message should indicate internal error")
	mockRepo.AssertNotCalled(t, "SaveRevert", mock.Anything)
	mockRepo.AssertExpectations(t)
}
