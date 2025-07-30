package database

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/diogocarasco/go-pharmacy-service/internal/models"
	_ "github.com/mattn/go-sqlite3"
)

// DBRepository defines the interface for database operations.
type DBRepository interface {
	SavePharmacy(pharmacy models.Pharmacy) error
	GetPharmacyByNPI(npi string) (*models.Pharmacy, error)
	SaveClaim(claim models.Claim) error
	GetClaimByID(id string) (*models.Claim, error)
	UpdateClaimRevertedStatus(id string, reverted bool) error
	SaveRevert(revert models.Revert) error
	Close() error
	SaveClaims(claims []models.Claim) error
	SaveReverts(reverts []models.Revert) error
}

// SQLiteRepository implements DBRepository for SQLite.
type SQLiteRepository struct {
	DB *sql.DB
}

// InitDB initializes the SQLite database connection.
func InitDB(dataSourceName string) (DBRepository, error) {
	db, err := sql.Open("sqlite3", dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("error connecting to the database: %w", err)
	}

	log.Printf("SQLite database connection established at: %s", dataSourceName)
	return &SQLiteRepository{DB: db}, nil
}

// Close closes the database connection.
func (s *SQLiteRepository) Close() error {
	return s.DB.Close()
}

// SavePharmacy inserts or updates a pharmacy in the database.
func (s *SQLiteRepository) SavePharmacy(pharmacy models.Pharmacy) error {
	stmt, err := s.DB.Prepare("INSERT OR REPLACE INTO pharmacies(chain, npi) VALUES(?, ?)")
	if err != nil {
		return fmt.Errorf("error preparing SavePharmacy statement: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(pharmacy.Chain, pharmacy.NPI)
	if err != nil {
		return fmt.Errorf("error executing SavePharmacy: %w", err)
	}
	return nil
}

// GetPharmacyByNPI fetches a pharmacy by its NPI.
func (s *SQLiteRepository) GetPharmacyByNPI(npi string) (*models.Pharmacy, error) {
	row := s.DB.QueryRow("SELECT chain, npi FROM pharmacies WHERE npi = ?", npi)

	var pharmacy models.Pharmacy
	err := row.Scan(&pharmacy.Chain, &pharmacy.NPI)
	if err == sql.ErrNoRows {
		return nil, nil // Not found
	}
	if err != nil {
		return nil, fmt.Errorf("error scanning pharmacy by NPI %s: %w", npi, err)
	}
	return &pharmacy, nil
}

// SaveClaim inserts a new claim into the database.
func (s *SQLiteRepository) SaveClaim(claim models.Claim) error {
	stmt, err := s.DB.Prepare(`
        INSERT INTO claims(id, ndc, npi, quantity, price, timestamp, reverted)
        VALUES(?, ?, ?, ?, ?, ?, ?)
        ON CONFLICT(id) DO UPDATE SET
            ndc = excluded.ndc,
            npi = excluded.npi,
            quantity = excluded.quantity,
            price = excluded.price,
            timestamp = excluded.timestamp,
            reverted = excluded.reverted;
    `)
	if err != nil {
		return fmt.Errorf("error preparing statement to insert/update claim: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(
		claim.ID,
		claim.NDC,
		claim.NPI,
		claim.Quantity,
		claim.Price,
		claim.Timestamp,
		claim.Reverted,
	)
	if err != nil {
		return fmt.Errorf("error executing insert/update for claim %s: %w", claim.ID, err)
	}
	return nil
}

// SaveClaims inserts multiple claims into the database within a transaction.
func (s *SQLiteRepository) SaveClaims(claims []models.Claim) error {
	tx, err := s.DB.Begin()
	if err != nil {
		return fmt.Errorf("error starting transaction for claims: %w", err)
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(`
        INSERT INTO claims (id, ndc, npi, quantity, price, timestamp, reverted)
        VALUES (?, ?, ?, ?, ?, ?, ?)
        ON CONFLICT(id) DO UPDATE SET
            ndc = excluded.ndc,
            npi = excluded.npi,
            quantity = excluded.quantity,
            price = excluded.price,
            timestamp = excluded.timestamp,
            reverted = excluded.reverted;
    `)
	if err != nil {
		return fmt.Errorf("error preparing statement to save claims in batch: %w", err)
	}
	defer stmt.Close()

	for _, claim := range claims {
		_, err := stmt.Exec(claim.ID, claim.NDC, claim.NPI, claim.Quantity, claim.Price, claim.Timestamp, claim.Reverted)
		if err != nil {
			return fmt.Errorf("error executing insert/update for claim %s: %w", claim.ID, err)
		}
	}

	return tx.Commit()
}

// GetClaimByID fetches a claim by its ID.
func (s *SQLiteRepository) GetClaimByID(id string) (*models.Claim, error) {
	row := s.DB.QueryRow("SELECT id, ndc, npi, quantity, price, timestamp, reverted FROM claims WHERE id = ?", id)

	var claim models.Claim

	err := row.Scan(
		&claim.ID,
		&claim.NDC,
		&claim.NPI,
		&claim.Quantity,
		&claim.Price,
		&claim.Timestamp,
		&claim.Reverted,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("error scanning claim by ID %s: %w", id, err)
	}
	return &claim, nil
}

// UpdateClaimRevertedStatus updates the 'reverted' status of a claim.
func (s *SQLiteRepository) UpdateClaimRevertedStatus(id string, reverted bool) error {
	stmt, err := s.DB.Prepare("UPDATE claims SET reverted = ? WHERE id = ?")
	if err != nil {
		return fmt.Errorf("error preparing statement to update claim status: %w", err)
	}
	defer stmt.Close()

	res, err := stmt.Exec(reverted, id)
	if err != nil {
		return fmt.Errorf("error executing status update for claim %s: %w", id, err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("claim with ID '%s' not found for status update", id)
	}
	return nil
}

// SaveRevert inserts a new revert record into the database.
func (s *SQLiteRepository) SaveRevert(revert models.Revert) error {
	stmt, err := s.DB.Prepare(`
        INSERT INTO reverts(id, claim_id, timestamp)
        VALUES(?, ?, ?)
        ON CONFLICT(id) DO UPDATE SET
            claim_id = excluded.claim_id,
            timestamp = excluded.timestamp;
    `)
	if err != nil {
		return fmt.Errorf("error preparing statement to insert/update revert: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(
		revert.ID,
		revert.ClaimID,
		revert.Timestamp,
	)
	if err != nil {
		return fmt.Errorf("error executing insert/update for revert %s: %w", revert.ID, err)
	}
	return nil
}

// SaveReverts inserts multiple reverts into the database within a transaction.
func (s *SQLiteRepository) SaveReverts(reverts []models.Revert) error {
	tx, err := s.DB.Begin()
	if err != nil {
		return fmt.Errorf("error starting transaction for reverts: %w", err)
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(`
        INSERT INTO reverts (id, claim_id, timestamp)
        VALUES (?, ?, ?)
        ON CONFLICT(id) DO UPDATE SET
            claim_id = excluded.claim_id,
            timestamp = excluded.timestamp;
    `)
	if err != nil {
		return fmt.Errorf("error preparing statement to save reverts in batch: %w", err)
	}
	defer stmt.Close()

	for _, revert := range reverts {
		_, err := stmt.Exec(revert.ID, revert.ClaimID, revert.Timestamp)
		if err != nil {
			return fmt.Errorf("error executing insert/update for revert %s: %w", revert.ID, err)
		}
	}

	return tx.Commit()
}
