package repository

import (
	"database/sql"
	"holding-hr-system/internal/models"
)

type CompanyRepository struct {
	db *sql.DB
}

func NewCompanyRepository(db *sql.DB) *CompanyRepository {
	return &CompanyRepository{db: db}
}

func (r *CompanyRepository) GetAll() ([]models.Company, error) {
	query := `SELECT id, name, is_holding, tax_id, address, created_at FROM companies ORDER BY is_holding DESC, name`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var companies []models.Company
	for rows.Next() {
		var c models.Company
		var taxID, address sql.NullString
		err := rows.Scan(
			&c.ID,
			&c.Name,
			&c.IsHolding,
			&taxID,
			&address,
			&c.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		if taxID.Valid {
			c.TaxID = taxID.String
		}
		if address.Valid {
			c.Address = address.String
		}
		companies = append(companies, c)
	}

	return companies, nil
}

func (r *CompanyRepository) GetByID(id int) (*models.Company, error) {
	query := `SELECT id, name, is_holding, tax_id, address, created_at FROM companies WHERE id = ?`

	var c models.Company
	var taxID, address sql.NullString
	err := r.db.QueryRow(query, id).Scan(
		&c.ID,
		&c.Name,
		&c.IsHolding,
		&taxID,
		&address,
		&c.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	if taxID.Valid {
		c.TaxID = taxID.String
	}
	if address.Valid {
		c.Address = address.String
	}

	return &c, nil
}

func (r *CompanyRepository) Create(company *models.Company) error {
	query := `INSERT INTO companies (name, is_holding, tax_id, address) VALUES (?, ?, ?, ?)`

	result, err := r.db.Exec(query, company.Name, company.IsHolding, company.TaxID, company.Address)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	company.ID = int(id)
	return nil
}

func (r *CompanyRepository) Update(company *models.Company) error {
	query := `UPDATE companies SET name = ?, is_holding = ?, tax_id = ?, address = ? WHERE id = ?`

	_, err := r.db.Exec(query, company.Name, company.IsHolding, company.TaxID, company.Address, company.ID)
	return err
}

func (r *CompanyRepository) Delete(id int) error {
	query := `DELETE FROM companies WHERE id = ?`
	_, err := r.db.Exec(query, id)
	return err
}

// Subsidiaries - Yalnız alt şirkətləri qaytarır
func (r *CompanyRepository) GetSubsidiaries() ([]models.Company, error) {
	query := `SELECT id, name, is_holding, tax_id, address, created_at FROM companies WHERE is_holding = FALSE ORDER BY name`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var companies []models.Company
	for rows.Next() {
		var c models.Company
		var taxID, address sql.NullString
		err := rows.Scan(
			&c.ID,
			&c.Name,
			&c.IsHolding,
			&taxID,
			&address,
			&c.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		if taxID.Valid {
			c.TaxID = taxID.String
		}
		if address.Valid {
			c.Address = address.String
		}
		companies = append(companies, c)
	}

	return companies, nil
}
