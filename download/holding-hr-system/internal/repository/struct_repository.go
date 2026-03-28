package repository

import (
	"database/sql"
	"holding-hr-system/internal/models"
)

type DepartmentRepository struct {
	db *sql.DB
}

func NewDepartmentRepository(db *sql.DB) *DepartmentRepository {
	return &DepartmentRepository{db: db}
}

func (r *DepartmentRepository) GetByCompanyID(companyID int) ([]models.Department, error) {
	query := `SELECT id, company_id, name FROM departments WHERE company_id = ? ORDER BY name`

	rows, err := r.db.Query(query, companyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var departments []models.Department
	for rows.Next() {
		var d models.Department
		err := rows.Scan(&d.ID, &d.CompanyID, &d.Name)
		if err != nil {
			return nil, err
		}
		departments = append(departments, d)
	}

	return departments, nil
}

func (r *DepartmentRepository) Create(dept *models.Department) error {
	query := `INSERT INTO departments (company_id, name) VALUES (?, ?)`

	result, err := r.db.Exec(query, dept.CompanyID, dept.Name)
	if err != nil {
		return err
	}

	id, _ := result.LastInsertId()
	dept.ID = int(id)
	return nil
}

func (r *DepartmentRepository) Delete(id int) error {
	_, err := r.db.Exec(`DELETE FROM departments WHERE id = ?`, id)
	return err
}

type PositionRepository struct {
	db *sql.DB
}

func NewPositionRepository(db *sql.DB) *PositionRepository {
	return &PositionRepository{db: db}
}

func (r *PositionRepository) GetByCompanyID(companyID int) ([]models.Position, error) {
	query := `SELECT id, company_id, name FROM positions WHERE company_id = ? ORDER BY name`

	rows, err := r.db.Query(query, companyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var positions []models.Position
	for rows.Next() {
		var p models.Position
		err := rows.Scan(&p.ID, &p.CompanyID, &p.Name)
		if err != nil {
			return nil, err
		}
		positions = append(positions, p)
	}

	return positions, nil
}

func (r *PositionRepository) Create(pos *models.Position) error {
	query := `INSERT INTO positions (company_id, name) VALUES (?, ?)`

	result, err := r.db.Exec(query, pos.CompanyID, pos.Name)
	if err != nil {
		return err
	}

	id, _ := result.LastInsertId()
	pos.ID = int(id)
	return nil
}

func (r *PositionRepository) Delete(id int) error {
	_, err := r.db.Exec(`DELETE FROM positions WHERE id = ?`, id)
	return err
}
