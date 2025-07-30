package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/diogocarasco/go-pharmacy-service/internal/logger"
	"github.com/diogocarasco/go-pharmacy-service/internal/models"
	"github.com/diogocarasco/go-pharmacy-service/internal/service"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type Handlers struct {
	claimService service.ClaimService
	logger       logger.Logger
	validator    *validator.Validate
}

func NewHandlers(claimService service.ClaimService, log logger.Logger) *Handlers {
	return &Handlers{
		claimService: claimService,
		logger:       log,
		validator:    validator.New(),
	}
}

// HealthCheckHandler responds with an OK status for application health checks.
// @Summary Checks application health
// @Description Returns an "ok" status if the application is running.
// @Tags health
// @Produce json
// @Success 200 {object} map[string]string "Status OK"
// @Router /health [get]
func (h *Handlers) HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	h.logger.Info("Health check performed.")
}

// SubmitClaimHandler handles claim submission via HTTP POST.
// @Summary Submit a new claim
// @Description Receives claim data and processes it
// @Tags claims
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param claim body models.ClaimSubmissionRequest true "Claim data to submit"
// @Success 200 {object} models.Claim "Claim submitted successfully"
// @Failure 400 "Invalid request"
// @Failure 500 "Internal server error"
// @Router /claims [post]
func (h *Handlers) SubmitClaimHandler(w http.ResponseWriter, r *http.Request) {
	var req models.ClaimSubmissionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("Error decoding SubmitClaim request: %v", err)
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	if err := h.validator.Struct(req); err != nil {
		h.logger.Error("Validation error for SubmitClaim: %v", err)
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	claim, err := h.claimService.SubmitClaim(req)
	if err != nil {
		h.logger.Error("Error submitting claim: %v", err)
		if strings.Contains(err.Error(), "pharmacy with NPI") { // Original: "farmácia com NPI"
			http.Error(w, "", http.StatusBadRequest)
		} else {
			http.Error(w, "", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(claim)
	h.logger.Info("Claim %s processed successfully via API.", claim.ID)
}

// GetClaimByIDHandler fetches a claim by its ID via HTTP GET.
// @Summary Get claim by ID
// @Description Returns the details of a specific claim by its ID
// @Tags claims
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Claim ID"
// @Success 200 {object} models.Claim "Claim details"
// @Failure 400 "Claim ID not provided"
// @Failure 404 "Claim not found"
// @Failure 500 "Internal server error"
// @Router /claims/{id} [get]
func (h *Handlers) GetClaimByIDHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if id == "" {
		h.logger.Error("Error: Claim ID not provided in the request.")
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	claim, err := h.claimService.GetClaimByID(id)
	if err != nil {
		h.logger.Error("Error fetching claim %s: %v", id, err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	if claim == nil {
		h.logger.Info("Claim with ID %s not found.", id)
		http.Error(w, "", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(claim)
}

// ReverseClaimHandler handles claim reversal via HTTP POST.
// @Summary Reverse an existing claim
// @Description Reverts an already submitted claim and records the reversal
// @Tags claims
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param reversal body models.ClaimReversalRequest true "Claim ID to be reverted"
// @Success 200 {object} object "Reversal successfully recorded"
// @Failure 400 "Invalid request or claim already reverted/not found"
// @Failure 500 "Internal server error"
// @Router /claims/reverse [post]
func (h *Handlers) ReverseClaimHandler(w http.ResponseWriter, r *http.Request) {
	var req models.ClaimReversalRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("Error decoding ReverseClaim request: %v", err)
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	if err := h.validator.Struct(req); err != nil {
		h.logger.Error("Validation error for ReverseClaim: %v", err)
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	revert, err := h.claimService.ReverseClaim(req)
	if err != nil {
		h.logger.Error("Error reverting claim: %v", err)
		if strings.Contains(err.Error(), "claim with ID") {
			http.Error(w, "", http.StatusBadRequest)
		} else {
			http.Error(w, "", http.StatusInternalServerError)
		}
		return
	}

	response := map[string]string{
		"status":   "claim reversed",
		"claim_id": revert.ClaimID,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
	h.logger.Info("Claim %s reverted successfully via API.", revert.ClaimID)
}

type ErrorResponse struct {
	Message string `json:"message"`
}
