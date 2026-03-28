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

// GetAll - Bütün şirkətləri gətir (hiyerarşiksiz)
func (r *CompanyRepository) GetAll() ([]models.Company, error) {
        query := `SELECT id, name, parent_id, is_holding, company_type, tax_id, address, created_at FROM companies ORDER BY is_holding DESC, parent_id IS NOT NULL, parent_id, name`

        rows, err := r.db.Query(query)
        if err != nil {
                return nil, err
        }
        defer rows.Close()

        var companies []models.Company
        for rows.Next() {
                c, err := r.scanCompany(rows)
                if err != nil {
                        return nil, err
                }
                companies = append(companies, *c)
        }

        return companies, nil
}

// GetTopLevel - Yalnız əsas şirkətləri gətir (parent_id = NULL olanlar)
func (r *CompanyRepository) GetTopLevel() ([]models.Company, error) {
        query := `SELECT id, name, parent_id, is_holding, company_type, tax_id, address, created_at FROM companies WHERE parent_id IS NULL ORDER BY is_holding DESC, name`

        rows, err := r.db.Query(query)
        if err != nil {
                return nil, err
        }
        defer rows.Close()

        var companies []models.Company
        for rows.Next() {
                c, err := r.scanCompany(rows)
                if err != nil {
                        return nil, err
                }
                companies = append(companies, *c)
        }

        return companies, nil
}

// GetWithHierarchy - Şirkətləri ierarxik qaydada gətir (alt-şirkətlərlə birlikdə)
func (r *CompanyRepository) GetWithHierarchy() ([]models.Company, error) {
        allCompanies, err := r.GetAll()
        if err != nil {
                return nil, err
        }

        // Employee saylarını hesabla
        counts, err := r.getEmployeeCounts()
        if err != nil {
                counts = make(map[int]int)
        }

        // Alt-şirkətləri parent-lərə bağla
        topLevel := make([]models.Company, 0)
        childMap := make(map[int][]models.Company)

        for _, c := range allCompanies {
                c.EmployeeCount = counts[c.ID]
                if c.ParentID != nil {
                        childMap[*c.ParentID] = append(childMap[*c.ParentID], c)
                } else {
                        topLevel = append(topLevel, c)
                }
        }

        // Rekursiv olaraq SubCompanies doldur
        for i := range topLevel {
                r.buildHierarchy(&topLevel[i], childMap)
        }

        return topLevel, nil
}

// buildHierarchy - Rekursiv olaraq alt-şirkətləri əlavə edir
func (r *CompanyRepository) buildHierarchy(company *models.Company, childMap map[int][]models.Company) {
        if children, ok := childMap[company.ID]; ok {
                company.SubCompanies = children
                // Parent-in sayına uşaqların saylarını da əlavə et
                for _, child := range children {
                        r.buildHierarchy(&child, childMap)
                }
                // Ümumi sayı hesabla (bu şirkətin + alt-şirkətlərin işçiləri)
                totalCount := company.EmployeeCount
                for _, child := range children {
                        totalCount += countAllEmployees(child)
                }
                company.EmployeeCount = totalCount
        }
}

// countAllEmployees - Rekursiv olaraq bütün alt-şirkətlərin işçi sayını hesablayır
func countAllEmployees(c models.Company) int {
        total := c.EmployeeCount
        for _, sub := range c.SubCompanies {
                total += countAllEmployees(sub)
        }
        return total
}

// GetSubCompanies - Bir şirkətin birbaşa alt-şirkətlərini gətir
func (r *CompanyRepository) GetSubCompanies(parentID int) ([]models.Company, error) {
        query := `SELECT id, name, parent_id, is_holding, company_type, tax_id, address, created_at FROM companies WHERE parent_id = ? ORDER BY name`

        rows, err := r.db.Query(query, parentID)
        if err != nil {
                return nil, err
        }
        defer rows.Close()

        var companies []models.Company
        for rows.Next() {
                c, err := r.scanCompany(rows)
                if err != nil {
                        return nil, err
                }
                companies = append(companies, *c)
        }

        return companies, nil
}

// GetByID - ID-yə görə şirkət gətir
func (r *CompanyRepository) GetByID(id int) (*models.Company, error) {
        query := `SELECT id, name, parent_id, is_holding, company_type, tax_id, address, created_at FROM companies WHERE id = ?`

        c := &models.Company{}
        var parentID sql.NullInt64
        var taxID, address, companyType sql.NullString
        err := r.db.QueryRow(query, id).Scan(
                &c.ID,
                &c.Name,
                &parentID,
                &c.IsHolding,
                &companyType,
                &taxID,
                &address,
                &c.CreatedAt,
        )

        if err != nil {
                return nil, err
        }

        if parentID.Valid {
                pid := int(parentID.Int64)
                c.ParentID = &pid
        }
        if taxID.Valid {
                c.TaxID = taxID.String
        }
        if address.Valid {
                c.Address = address.String
        }
        if companyType.Valid {
                c.CompanyType = companyType.String
        }

        return c, nil
}

// GetParentCompany - Şirkətin ana şirkətini gətir
func (r *CompanyRepository) GetParentCompany(companyID int) (*models.Company, error) {
        query := `SELECT c.id, c.name, c.parent_id, c.is_holding, c.company_type, c.tax_id, c.address, c.created_at
                  FROM companies c
                  INNER JOIN companies sub ON sub.parent_id = c.id
                  WHERE sub.id = ?`

        c := &models.Company{}
        var parentID sql.NullInt64
        var taxID, address, companyType sql.NullString
        err := r.db.QueryRow(query, companyID).Scan(
                &c.ID,
                &c.Name,
                &parentID,
                &c.IsHolding,
                &companyType,
                &taxID,
                &address,
                &c.CreatedAt,
        )

        if err != nil {
                return nil, err
        }

        if parentID.Valid {
                pid := int(parentID.Int64)
                c.ParentID = &pid
        }
        if taxID.Valid {
                c.TaxID = taxID.String
        }
        if address.Valid {
                c.Address = address.String
        }
        if companyType.Valid {
                c.CompanyType = companyType.String
        }

        return c, nil
}

// HasSubCompanies - Şirkətin alt-şirkətləri var mı?
func (r *CompanyRepository) HasSubCompanies(companyID int) (bool, error) {
        var count int
        err := r.db.QueryRow("SELECT COUNT(*) FROM companies WHERE parent_id = ?", companyID).Scan(&count)
        return count > 0, err
}

// GetCompaniesOrEmpty - Hierarchy ilə şirkətləri gətir, xəta olsa boz qaytar
func (r *CompanyRepository) GetCompaniesOrEmpty() []models.Company {
        companies, err := r.GetWithHierarchy()
        if err != nil {
                return nil
        }
        return companies
}

// GetEmployeeCount - Şirkətin işçi sayını gətir
func (r *CompanyRepository) GetEmployeeCount(companyID int) (int, error) {
        var count int
        err := r.db.QueryRow("SELECT COUNT(*) FROM employees WHERE company_id = ? AND status = 'ACTIVE'", companyID).Scan(&count)
        return count, err
}

// GetEmployeeCountRecursive - Şirkət və bütün alt-şirkətlərinin işçi sayını gətir
func (r *CompanyRepository) GetEmployeeCountRecursive(companyID int) (int, error) {
        // Bu şirkətin direkt işçiləri
        var directCount int
        err := r.db.QueryRow("SELECT COUNT(*) FROM employees WHERE company_id = ? AND status = 'ACTIVE'", companyID).Scan(&directCount)
        if err != nil {
                return 0, err
        }

        // Alt-şirkətlərin sayları
        subs, err := r.GetSubCompanies(companyID)
        if err != nil {
                return directCount, nil
        }

        total := directCount
        for _, sub := range subs {
                subCount, err := r.GetEmployeeCountRecursive(sub.ID)
                if err == nil {
                        total += subCount
                }
        }

        return total, nil
}

// GetAllSubCompanyIDs - Şirkət və bütün alt-şirkətlərinin ID-lərini gətir
func (r *CompanyRepository) GetAllSubCompanyIDs(companyID int) ([]int, error) {
        ids := []int{companyID}

        subs, err := r.GetSubCompanies(companyID)
        if err != nil {
                return ids, nil
        }

        for _, sub := range subs {
                subIDs, err := r.GetAllSubCompanyIDs(sub.ID)
                if err == nil {
                        ids = append(ids, subIDs...)
                }
        }

        return ids, nil
}

func (r *CompanyRepository) Create(company *models.Company) error {
        query := `INSERT INTO companies (name, parent_id, is_holding, company_type, tax_id, address) VALUES (?, ?, ?, ?, ?, ?)`

        result, err := r.db.Exec(query, company.Name, company.ParentID, company.IsHolding, company.CompanyType, company.TaxID, company.Address)
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
        query := `UPDATE companies SET name = ?, parent_id = ?, is_holding = ?, company_type = ?, tax_id = ?, address = ? WHERE id = ?`

        _, err := r.db.Exec(query, company.Name, company.ParentID, company.IsHolding, company.CompanyType, company.TaxID, company.Address, company.ID)
        return err
}

func (r *CompanyRepository) Delete(id int) error {
        query := `DELETE FROM companies WHERE id = ?`
        _, err := r.db.Exec(query, id)
        return err
}

// Subsidiaries - Yalnız əsas (top-level) alt şirkətləri qaytarır
func (r *CompanyRepository) GetSubsidiaries() ([]models.Company, error) {
        query := `SELECT id, name, parent_id, is_holding, company_type, tax_id, address, created_at FROM companies WHERE parent_id IS NULL AND is_holding = FALSE ORDER BY name`

        rows, err := r.db.Query(query)
        if err != nil {
                return nil, err
        }
        defer rows.Close()

        var companies []models.Company
        for rows.Next() {
                c, err := r.scanCompany(rows)
                if err != nil {
                        return nil, err
                }
                companies = append(companies, *c)
        }

        return companies, nil
}

// scanCompany - Rows-dan Company skan edir
func (r *CompanyRepository) scanCompany(rows *sql.Rows) (*models.Company, error) {
        c := &models.Company{}
        var parentID sql.NullInt64
        var taxID, address, companyType sql.NullString
        err := rows.Scan(
                &c.ID,
                &c.Name,
                &parentID,
                &c.IsHolding,
                &companyType,
                &taxID,
                &address,
                &c.CreatedAt,
        )
        if err != nil {
                return nil, err
        }
        if parentID.Valid {
                pid := int(parentID.Int64)
                c.ParentID = &pid
        }
        if taxID.Valid {
                c.TaxID = taxID.String
        }
        if address.Valid {
                c.Address = address.String
        }
        if companyType.Valid {
                c.CompanyType = companyType.String
        }
        return c, nil
}

// getEmployeeCounts - Bütün şirkətlərin ACTIVE işçi saylarını gətir
func (r *CompanyRepository) getEmployeeCounts() (map[int]int, error) {
        query := `SELECT company_id, COUNT(*) as cnt FROM employees WHERE status = 'ACTIVE' GROUP BY company_id`
        rows, err := r.db.Query(query)
        if err != nil {
                return nil, err
        }
        defer rows.Close()

        counts := make(map[int]int)
        for rows.Next() {
                var companyID, count int
                if err := rows.Scan(&companyID, &count); err != nil {
                        continue
                }
                counts[companyID] = count
        }
        return counts, nil
}
