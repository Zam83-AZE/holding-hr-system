package repository

import (
        "database/sql"
        "fmt"
        "holding-hr-system/internal/models"
        "time"
)

type EmployeeRepository struct {
        db *sql.DB
}

func NewEmployeeRepository(db *sql.DB) *EmployeeRepository {
        return &EmployeeRepository{db: db}
}

// GetByStatus - Statusa görə işçiləri gətir
func (r *EmployeeRepository) GetByStatus(companyID *int, status models.EmployeeStatus) ([]models.Employee, error) {
        var query string
        var args []interface{}

        if companyID != nil {
                query = `SELECT e.id, e.company_id, e.first_name, e.last_name, e.father_name, e.fin_code,
                                e.birth_date, e.gender, e.photo_path, e.phone, e.email, e.address, e.status,
                                e.department_id, e.position_id, e.hire_date, e.termination_date, e.termination_reason,
                                e.created_at, e.updated_at,
                                c.name as company_name, d.name as department_name, p.name as position_name
                                FROM employees e
                                LEFT JOIN companies c ON e.company_id = c.id
                                LEFT JOIN departments d ON e.department_id = d.id
                                LEFT JOIN positions p ON e.position_id = p.id
                                WHERE e.company_id = ? AND e.status = ?
                                ORDER BY e.created_at DESC`
                args = []interface{}{*companyID, string(status)}
        } else {
                query = `SELECT e.id, e.company_id, e.first_name, e.last_name, e.father_name, e.fin_code,
                                e.birth_date, e.gender, e.photo_path, e.phone, e.email, e.address, e.status,
                                e.department_id, e.position_id, e.hire_date, e.termination_date, e.termination_reason,
                                e.created_at, e.updated_at,
                                c.name as company_name, d.name as department_name, p.name as position_name
                                FROM employees e
                                LEFT JOIN companies c ON e.company_id = c.id
                                LEFT JOIN departments d ON e.department_id = d.id
                                LEFT JOIN positions p ON e.position_id = p.id
                                WHERE e.status = ?
                                ORDER BY e.created_at DESC`
                args = []interface{}{string(status)}
        }

        return r.queryEmployees(query, args...)
}

// GetByID - ID-yə görə işçi gətir
func (r *EmployeeRepository) GetByID(id int) (*models.Employee, error) {
        query := `SELECT e.id, e.company_id, e.first_name, e.last_name, e.father_name, e.fin_code,
                        e.birth_date, e.gender, e.photo_path, e.phone, e.email, e.address, e.status,
                        e.department_id, e.position_id, e.hire_date, e.termination_date, e.termination_reason,
                        e.created_at, e.updated_at,
                        c.name as company_name, d.name as department_name, p.name as position_name
                        FROM employees e
                        LEFT JOIN companies c ON e.company_id = c.id
                        LEFT JOIN departments d ON e.department_id = d.id
                        LEFT JOIN positions p ON e.position_id = p.id
                        WHERE e.id = ?`

        employee := &models.Employee{}
        var birthDate, hireDate, termDate sql.NullTime
        var gender sql.NullString
        var deptID, posID sql.NullInt64
        var deptName, posName, termReason sql.NullString

        err := r.db.QueryRow(query, id).Scan(
                &employee.ID, &employee.CompanyID, &employee.FirstName, &employee.LastName,
                &employee.FatherName, &employee.FINCode, &birthDate, &gender, &employee.PhotoPath,
                &employee.Phone, &employee.Email, &employee.Address, &employee.Status,
                &deptID, &posID, &hireDate, &termDate, &termReason,
                &employee.CreatedAt, &employee.UpdatedAt,
                &employee.CompanyName, &deptName, &posName,
        )

        if err != nil {
                return nil, err
        }

        if birthDate.Valid {
                employee.BirthDate = &birthDate.Time
        }
        if gender.Valid {
                g := models.Gender(gender.String)
                employee.Gender = &g
        }
        if deptID.Valid {
                id := int(deptID.Int64)
                employee.DepartmentID = &id
        }
        if posID.Valid {
                id := int(posID.Int64)
                employee.PositionID = &id
        }
        if hireDate.Valid {
                employee.HireDate = &hireDate.Time
        }
        if termDate.Valid {
                employee.TerminationDate = &termDate.Time
        }
        if termReason.Valid {
                employee.TerminationReason = termReason.String
        }
        if deptName.Valid {
                employee.DepartmentName = deptName.String
        }
        if posName.Valid {
                employee.PositionName = posName.String
        }

        return employee, nil
}

// Create - Yeni işçi yarat
func (r *EmployeeRepository) Create(employee *models.Employee) error {
        query := `INSERT INTO employees (company_id, first_name, last_name, father_name, fin_code,
                        birth_date, gender, photo_path, phone, email, address, status, department_id, position_id, hire_date)
                        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

        result, err := r.db.Exec(query,
                employee.CompanyID,
                employee.FirstName,
                employee.LastName,
                employee.FatherName,
                employee.FINCode,
                employee.BirthDate,
                employee.Gender,
                employee.PhotoPath,
                employee.Phone,
                employee.Email,
                employee.Address,
                employee.Status,
                employee.DepartmentID,
                employee.PositionID,
                employee.HireDate,
        )

        if err != nil {
                return err
        }

        id, err := result.LastInsertId()
        if err != nil {
                return err
        }

        employee.ID = int(id)
        return nil
}

// Update - İşçi məlumatlarını yenilə
func (r *EmployeeRepository) Update(employee *models.Employee) error {
        query := `UPDATE employees SET
                        first_name = ?, last_name = ?, father_name = ?, fin_code = ?,
                        birth_date = ?, gender = ?, photo_path = ?, phone = ?, email = ?,
                        address = ?, department_id = ?, position_id = ?
                        WHERE id = ?`

        _, err := r.db.Exec(query,
                employee.FirstName,
                employee.LastName,
                employee.FatherName,
                employee.FINCode,
                employee.BirthDate,
                employee.Gender,
                employee.PhotoPath,
                employee.Phone,
                employee.Email,
                employee.Address,
                employee.DepartmentID,
                employee.PositionID,
                employee.ID,
        )

        return err
}

// UpdateStatus - Statusu dəyiş
func (r *EmployeeRepository) UpdateStatus(id int, status models.EmployeeStatus, hireDate, terminationDate *time.Time, terminationReason string) error {
        query := `UPDATE employees SET status = ?, hire_date = ?, termination_date = ?, termination_reason = ? WHERE id = ?`
        _, err := r.db.Exec(query, string(status), hireDate, terminationDate, terminationReason, id)
        return err
}

// Delete - İşçi sil
func (r *EmployeeRepository) Delete(id int) error {
        query := `DELETE FROM employees WHERE id = ?`
        _, err := r.db.Exec(query, id)
        return err
}

// Search - Axtarış
func (r *EmployeeRepository) Search(companyID *int, query string) ([]models.Employee, error) {
        var sqlQuery string
        var args []interface{}

        searchTerm := "%" + query + "%"

        if companyID != nil {
                sqlQuery = `SELECT e.id, e.company_id, e.first_name, e.last_name, e.father_name, e.fin_code,
                                e.birth_date, e.gender, e.photo_path, e.phone, e.email, e.address, e.status,
                                e.department_id, e.position_id, e.hire_date, e.termination_date, e.termination_reason,
                                e.created_at, e.updated_at,
                                c.name as company_name, d.name as department_name, p.name as position_name
                                FROM employees e
                                LEFT JOIN companies c ON e.company_id = c.id
                                LEFT JOIN departments d ON e.department_id = d.id
                                LEFT JOIN positions p ON e.position_id = p.id
                                WHERE e.company_id = ? AND (e.first_name LIKE ? OR e.last_name LIKE ? OR e.fin_code LIKE ?)
                                ORDER BY e.created_at DESC`
                args = []interface{}{*companyID, searchTerm, searchTerm, searchTerm}
        } else {
                sqlQuery = `SELECT e.id, e.company_id, e.first_name, e.last_name, e.father_name, e.fin_code,
                                e.birth_date, e.gender, e.photo_path, e.phone, e.email, e.address, e.status,
                                e.department_id, e.position_id, e.hire_date, e.termination_date, e.termination_reason,
                                e.created_at, e.updated_at,
                                c.name as company_name, d.name as department_name, p.name as position_name
                                FROM employees e
                                LEFT JOIN companies c ON e.company_id = c.id
                                LEFT JOIN departments d ON e.department_id = d.id
                                LEFT JOIN positions p ON e.position_id = p.id
                                WHERE e.first_name LIKE ? OR e.last_name LIKE ? OR e.fin_code LIKE ?
                                ORDER BY e.created_at DESC`
                args = []interface{}{searchTerm, searchTerm, searchTerm}
        }

        return r.queryEmployees(sqlQuery, args...)
}

// queryEmployees - ümumi query funksiyası
func (r *EmployeeRepository) queryEmployees(query string, args ...interface{}) ([]models.Employee, error) {
        rows, err := r.db.Query(query, args...)
        if err != nil {
                return nil, err
        }
        defer rows.Close()

        var employees []models.Employee
        for rows.Next() {
                var e models.Employee
                var birthDate, hireDate, termDate sql.NullTime
                var gender sql.NullString
                var deptID, posID sql.NullInt64
                var deptName, posName, termReason sql.NullString

                err := rows.Scan(
                        &e.ID, &e.CompanyID, &e.FirstName, &e.LastName,
                        &e.FatherName, &e.FINCode, &birthDate, &gender, &e.PhotoPath,
                        &e.Phone, &e.Email, &e.Address, &e.Status,
                        &deptID, &posID, &hireDate, &termDate, &termReason,
                        &e.CreatedAt, &e.UpdatedAt,
                        &e.CompanyName, &deptName, &posName,
                )

                if err != nil {
                        return nil, fmt.Errorf("scan error: %w", err)
                }

                if birthDate.Valid {
                        e.BirthDate = &birthDate.Time
                }
                if gender.Valid {
                        g := models.Gender(gender.String)
                        e.Gender = &g
                }
                if deptID.Valid {
                        id := int(deptID.Int64)
                        e.DepartmentID = &id
                }
                if posID.Valid {
                        id := int(posID.Int64)
                        e.PositionID = &id
                }
                if hireDate.Valid {
                        e.HireDate = &hireDate.Time
                }
                if termDate.Valid {
                        e.TerminationDate = &termDate.Time
                }
                if termReason.Valid {
                        e.TerminationReason = termReason.String
                }
                if deptName.Valid {
                        e.DepartmentName = deptName.String
                }
                if posName.Valid {
                        e.PositionName = posName.String
                }

                employees = append(employees, e)
        }

        return employees, nil
}

// GetEducation - Təhsil məlumatları
func (r *EmployeeRepository) GetEducation(employeeID int) ([]models.EmployeeEducation, error) {
        query := `SELECT id, employee_id, institution, specialty, degree, start_year, end_year, diploma_number
                        FROM employee_education WHERE employee_id = ? ORDER BY end_year DESC`

        rows, err := r.db.Query(query, employeeID)
        if err != nil {
                return nil, err
        }
        defer rows.Close()

        var education []models.EmployeeEducation
        for rows.Next() {
                var e models.EmployeeEducation
                var specialty, diploma sql.NullString
                err := rows.Scan(&e.ID, &e.EmployeeID, &e.Institution, &specialty, &e.Degree, &e.StartYear, &e.EndYear, &diploma)
                if err != nil {
                        return nil, err
                }
                if specialty.Valid {
                        e.Specialty = specialty.String
                }
                if diploma.Valid {
                        e.DiplomaNumber = diploma.String
                }
                education = append(education, e)
        }

        return education, nil
}

// AddEducation - Təhsil əlavə et
func (r *EmployeeRepository) AddEducation(edu *models.EmployeeEducation) error {
        query := `INSERT INTO employee_education (employee_id, institution, specialty, degree, start_year, end_year, diploma_number)
                        VALUES (?, ?, ?, ?, ?, ?, ?)`

        result, err := r.db.Exec(query, edu.EmployeeID, edu.Institution, edu.Specialty, edu.Degree, edu.StartYear, edu.EndYear, edu.DiplomaNumber)
        if err != nil {
                return err
        }

        id, _ := result.LastInsertId()
        edu.ID = int(id)
        return nil
}

// DeleteEducation - Təhsil sil
func (r *EmployeeRepository) DeleteEducation(id int) error {
        _, err := r.db.Exec(`DELETE FROM employee_education WHERE id = ?`, id)
        return err
}

// UpdateEducation - Təhsil yenilə
func (r *EmployeeRepository) UpdateEducation(edu *models.EmployeeEducation) error {
        query := `UPDATE employee_education SET 
                        institution = ?, specialty = ?, degree = ?, 
                        start_year = ?, end_year = ?, diploma_number = ?
                        WHERE id = ?`
        _, err := r.db.Exec(query, edu.Institution, edu.Specialty, edu.Degree, edu.StartYear, edu.EndYear, edu.DiplomaNumber, edu.ID)
        return err
}

// GetEducationByID - ID-yə görə təhsil gətir
func (r *EmployeeRepository) GetEducationByID(id int) (*models.EmployeeEducation, error) {
        query := `SELECT id, employee_id, institution, specialty, degree, start_year, end_year, diploma_number
                        FROM employee_education WHERE id = ?`
        
        var e models.EmployeeEducation
        var specialty, diploma sql.NullString
        err := r.db.QueryRow(query, id).Scan(&e.ID, &e.EmployeeID, &e.Institution, &specialty, &e.Degree, &e.StartYear, &e.EndYear, &diploma)
        if err != nil {
                return nil, err
        }
        if specialty.Valid {
                e.Specialty = specialty.String
        }
        if diploma.Valid {
                e.DiplomaNumber = diploma.String
        }
        return &e, nil
}

// GetExperience - İş təcrübəsi
func (r *EmployeeRepository) GetExperience(employeeID int) ([]models.EmployeeExperience, error) {
        query := `SELECT id, employee_id, company_name, position, start_date, end_date, leaving_reason
                        FROM employee_experience WHERE employee_id = ? ORDER BY start_date DESC`

        rows, err := r.db.Query(query, employeeID)
        if err != nil {
                return nil, err
        }
        defer rows.Close()

        var experiences []models.EmployeeExperience
        for rows.Next() {
                var e models.EmployeeExperience
                var position, leavingReason sql.NullString
                var startDate, endDate sql.NullTime
                err := rows.Scan(&e.ID, &e.EmployeeID, &e.CompanyName, &position, &startDate, &endDate, &leavingReason)
                if err != nil {
                        return nil, err
                }
                if position.Valid {
                        e.Position = position.String
                }
                if leavingReason.Valid {
                        e.LeavingReason = leavingReason.String
                }
                if startDate.Valid {
                        e.StartDate = &startDate.Time
                }
                if endDate.Valid {
                        e.EndDate = &endDate.Time
                }
                experiences = append(experiences, e)
        }

        return experiences, nil
}

// AddExperience - İş təcrübəsi əlavə et
func (r *EmployeeRepository) AddExperience(exp *models.EmployeeExperience) error {
        query := `INSERT INTO employee_experience (employee_id, company_name, position, start_date, end_date, leaving_reason)
                        VALUES (?, ?, ?, ?, ?, ?)`

        result, err := r.db.Exec(query, exp.EmployeeID, exp.CompanyName, exp.Position, exp.StartDate, exp.EndDate, exp.LeavingReason)
        if err != nil {
                return err
        }

        id, _ := result.LastInsertId()
        exp.ID = int(id)
        return nil
}

// DeleteExperience - İş təcrübəsi sil
func (r *EmployeeRepository) DeleteExperience(id int) error {
        _, err := r.db.Exec(`DELETE FROM employee_experience WHERE id = ?`, id)
        return err
}

// UpdateExperience - İş təcrübəsi yenilə
func (r *EmployeeRepository) UpdateExperience(exp *models.EmployeeExperience) error {
        query := `UPDATE employee_experience SET 
                        company_name = ?, position = ?, start_date = ?, 
                        end_date = ?, leaving_reason = ?
                        WHERE id = ?`
        _, err := r.db.Exec(query, exp.CompanyName, exp.Position, exp.StartDate, exp.EndDate, exp.LeavingReason, exp.ID)
        return err
}

// GetExperienceByID - ID-yə görə təcrübə gətir
func (r *EmployeeRepository) GetExperienceByID(id int) (*models.EmployeeExperience, error) {
        query := `SELECT id, employee_id, company_name, position, start_date, end_date, leaving_reason
                        FROM employee_experience WHERE id = ?`
        
        var e models.EmployeeExperience
        var position, leavingReason sql.NullString
        var startDate, endDate sql.NullTime
        err := r.db.QueryRow(query, id).Scan(&e.ID, &e.EmployeeID, &e.CompanyName, &position, &startDate, &endDate, &leavingReason)
        if err != nil {
                return nil, err
        }
        if position.Valid {
                e.Position = position.String
        }
        if leavingReason.Valid {
                e.LeavingReason = leavingReason.String
        }
        if startDate.Valid {
                e.StartDate = &startDate.Time
        }
        if endDate.Valid {
                e.EndDate = &endDate.Time
        }
        return &e, nil
}

// GetFamily - Ailə məlumatları
func (r *EmployeeRepository) GetFamily(employeeID int) ([]models.EmployeeFamily, error) {
        query := `SELECT id, employee_id, relation_type, full_name, birth_date, contact_number
                        FROM employee_family WHERE employee_id = ?`

        rows, err := r.db.Query(query, employeeID)
        if err != nil {
                return nil, err
        }
        defer rows.Close()

        var family []models.EmployeeFamily
        for rows.Next() {
                var f models.EmployeeFamily
                var birthDate sql.NullTime
                var contact sql.NullString
                err := rows.Scan(&f.ID, &f.EmployeeID, &f.RelationType, &f.FullName, &birthDate, &contact)
                if err != nil {
                        return nil, err
                }
                if birthDate.Valid {
                        f.BirthDate = &birthDate.Time
                }
                if contact.Valid {
                        f.ContactNumber = contact.String
                }
                family = append(family, f)
        }

        return family, nil
}

// AddFamily - Ailə üzvü əlavə et
func (r *EmployeeRepository) AddFamily(fam *models.EmployeeFamily) error {
        query := `INSERT INTO employee_family (employee_id, relation_type, full_name, birth_date, contact_number)
                        VALUES (?, ?, ?, ?, ?)`

        result, err := r.db.Exec(query, fam.EmployeeID, fam.RelationType, fam.FullName, fam.BirthDate, fam.ContactNumber)
        if err != nil {
                return err
        }

        id, _ := result.LastInsertId()
        fam.ID = int(id)
        return nil
}

// DeleteFamily - Ailə üzvü sil
func (r *EmployeeRepository) DeleteFamily(id int) error {
        _, err := r.db.Exec(`DELETE FROM employee_family WHERE id = ?`, id)
        return err
}

// UpdateFamily - Ailə üzvü yenilə
func (r *EmployeeRepository) UpdateFamily(fam *models.EmployeeFamily) error {
        query := `UPDATE employee_family SET 
                        relation_type = ?, full_name = ?, birth_date = ?, contact_number = ?
                        WHERE id = ?`
        _, err := r.db.Exec(query, fam.RelationType, fam.FullName, fam.BirthDate, fam.ContactNumber, fam.ID)
        return err
}

// GetFamilyByID - ID-yə görə ailə üzvü gətir
func (r *EmployeeRepository) GetFamilyByID(id int) (*models.EmployeeFamily, error) {
        query := `SELECT id, employee_id, relation_type, full_name, birth_date, contact_number
                        FROM employee_family WHERE id = ?`
        
        var f models.EmployeeFamily
        var birthDate sql.NullTime
        var contact sql.NullString
        err := r.db.QueryRow(query, id).Scan(&f.ID, &f.EmployeeID, &f.RelationType, &f.FullName, &birthDate, &contact)
        if err != nil {
                return nil, err
        }
        if birthDate.Valid {
                f.BirthDate = &birthDate.Time
        }
        if contact.Valid {
                f.ContactNumber = contact.String
        }
        return &f, nil
}

// GetLifecycleLogs - Yaşam dövrü tarixçəsi
func (r *EmployeeRepository) GetLifecycleLogs(employeeID int) ([]models.EmployeeLifecycleLog, error) {
        query := `SELECT id, employee_id, user_id, old_status, new_status, event_description, created_at
                        FROM employee_lifecycle_logs WHERE employee_id = ? ORDER BY created_at DESC`

        rows, err := r.db.Query(query, employeeID)
        if err != nil {
                return nil, err
        }
        defer rows.Close()

        var logs []models.EmployeeLifecycleLog
        for rows.Next() {
                var l models.EmployeeLifecycleLog
                var oldStatus, newStatus sql.NullString
                err := rows.Scan(&l.ID, &l.EmployeeID, &l.UserID, &oldStatus, &newStatus, &l.EventDescription, &l.CreatedAt)
                if err != nil {
                        return nil, err
                }
                if oldStatus.Valid {
                        l.OldStatus = oldStatus.String
                }
                if newStatus.Valid {
                        l.NewStatus = newStatus.String
                }
                logs = append(logs, l)
        }

        return logs, nil
}

// AddLifecycleLog - Tarixçə əlavə et
func (r *EmployeeRepository) AddLifecycleLog(log *models.EmployeeLifecycleLog) error {
        query := `INSERT INTO employee_lifecycle_logs (employee_id, user_id, old_status, new_status, event_description)
                        VALUES (?, ?, ?, ?, ?)`

        result, err := r.db.Exec(query, log.EmployeeID, log.UserID, log.OldStatus, log.NewStatus, log.EventDescription)
        if err != nil {
                return err
        }

        id, _ := result.LastInsertId()
        log.ID = int(id)
        return nil
}

// GetStats - Dashboard statistikası
func (r *EmployeeRepository) GetStats(companyID *int) (*models.DashboardStats, error) {
        stats := &models.DashboardStats{}

        var whereClause string
        var args []interface{}

        if companyID != nil {
                whereClause = "WHERE company_id = ?"
                args = []interface{}{*companyID}
        }

        // Ümumi say
        query := fmt.Sprintf("SELECT COUNT(*) FROM employees %s", whereClause)
        r.db.QueryRow(query, args...).Scan(&stats.TotalEmployees)

        // Statuslara görə
        query = fmt.Sprintf("SELECT COUNT(*) FROM employees WHERE status = 'ACTIVE' %s", andClause(whereClause))
        r.db.QueryRow(query, append([]interface{}{"ACTIVE"}, args...)...).Scan(&stats.ActiveEmployees)

        query = fmt.Sprintf("SELECT COUNT(*) FROM employees WHERE status = 'CANDIDATE' %s", andClause(whereClause))
        r.db.QueryRow(query, append([]interface{}{"CANDIDATE"}, args...)...).Scan(&stats.Candidates)

        query = fmt.Sprintf("SELECT COUNT(*) FROM employees WHERE status = 'TERMINATED' %s", andClause(whereClause))
        r.db.QueryRow(query, append([]interface{}{"TERMINATED"}, args...)...).Scan(&stats.Terminated)

        // Bu ay işə qəbul
        query = fmt.Sprintf("SELECT COUNT(*) FROM employees WHERE status = 'ACTIVE' AND MONTH(hire_date) = MONTH(CURRENT_DATE()) AND YEAR(hire_date) = YEAR(CURRENT_DATE()) %s", andClause(whereClause))
        r.db.QueryRow(query, args...).Scan(&stats.ThisMonthHired)

        // Bu ay işdən çıxan
        query = fmt.Sprintf("SELECT COUNT(*) FROM employees WHERE status = 'TERMINATED' AND MONTH(termination_date) = MONTH(CURRENT_DATE()) AND YEAR(termination_date) = YEAR(CURRENT_DATE()) %s", andClause(whereClause))
        r.db.QueryRow(query, args...).Scan(&stats.ThisMonthTerminated)

        return stats, nil
}

func andClause(whereClause string) string {
        if whereClause != "" {
                return "AND company_id = ?"
        }
        return ""
}

// GetWorkLocationsByCompany - Şirkətə aid iş yerlərini gətir
func (r *EmployeeRepository) GetWorkLocationsByCompany(companyID int) ([]models.WorkLocation, error) {
	query := "SELECT id, company_id, name, address, type, is_active, created_at FROM work_locations WHERE company_id = ? AND is_active = TRUE ORDER BY name"
	rows, err := r.db.Query(query, companyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var locations []models.WorkLocation
	for rows.Next() {
		var loc models.WorkLocation
		if err := rows.Scan(&loc.ID, &loc.CompanyID, &loc.Name, &loc.Address, &loc.Type, &loc.IsActive, &loc.CreatedAt); err != nil {
			return nil, err
		}
		locations = append(locations, loc)
	}
	return locations, nil
}

// GetWorkLocationByID - İş yeri ID ilə gətir
func (r *EmployeeRepository) GetWorkLocationByID(id int) (*models.WorkLocation, error) {
	query := "SELECT id, company_id, name, address, type, is_active, created_at FROM work_locations WHERE id = ?"
	var loc models.WorkLocation
	err := r.db.QueryRow(query, id).Scan(&loc.ID, &loc.CompanyID, &loc.Name, &loc.Address, &loc.Type, &loc.IsActive, &loc.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &loc, nil
}

// GetCertificatesByEmployee - İşçinin sertifikatlarını gətir
func (r *EmployeeRepository) GetCertificatesByEmployee(employeeID int) ([]models.EmployeeCertificate, error) {
	query := "SELECT id, employee_id, certificate_type, certificate_number, issued_by, issue_date, expiry_date, status, notes, created_at, updated_at FROM employee_certificates WHERE employee_id = ? ORDER BY expiry_date ASC"
	rows, err := r.db.Query(query, employeeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var certs []models.EmployeeCertificate
	for rows.Next() {
		var c models.EmployeeCertificate
		var issueDate, expiryDate sql.NullTime
		var certNumber, issuedBy, notes, status sql.NullString
		if err := rows.Scan(&c.ID, &c.EmployeeID, &c.CertificateType, &certNumber, &issuedBy, &issueDate, &expiryDate, &status, &notes, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		if certNumber.Valid {
			c.CertificateNumber = certNumber.String
		}
		if issuedBy.Valid {
			c.IssuedBy = issuedBy.String
		}
		if notes.Valid {
			c.Notes = notes.String
		}
		if status.Valid {
			c.Status = status.String
		}
		if issueDate.Valid {
			c.IssueDate = &issueDate.Time
		}
		if expiryDate.Valid {
			c.ExpiryDate = &expiryDate.Time
		}
		certs = append(certs, c)
	}
	return certs, nil
}

// AddCertificate - Sertifikat əlavə et
func (r *EmployeeRepository) AddCertificate(cert *models.EmployeeCertificate) error {
	query := "INSERT INTO employee_certificates (employee_id, certificate_type, certificate_number, issued_by, issue_date, expiry_date, notes) VALUES (?, ?, ?, ?, ?, ?, ?)"
	result, err := r.db.Exec(query, cert.EmployeeID, cert.CertificateType, cert.CertificateNumber, cert.IssuedBy, cert.IssueDate, cert.ExpiryDate, cert.Notes)
	if err != nil {
		return err
	}
	id, _ := result.LastInsertId()
	cert.ID = int(id)
	return nil
}

// UpdateCertificate - Sertifikat yenilə
func (r *EmployeeRepository) UpdateCertificate(cert *models.EmployeeCertificate) error {
	query := "UPDATE employee_certificates SET certificate_type=?, certificate_number=?, issued_by=?, issue_date=?, expiry_date=?, notes=? WHERE id=?"
	_, err := r.db.Exec(query, cert.CertificateType, cert.CertificateNumber, cert.IssuedBy, cert.IssueDate, cert.ExpiryDate, cert.Notes, cert.ID)
	return err
}

// DeleteCertificate - Sertifikat sil
func (r *EmployeeRepository) DeleteCertificate(id int) error {
	_, err := r.db.Exec("DELETE FROM employee_certificates WHERE id=?", id)
	return err
}

// GetCertificateByID - Sertifikat ID ilə gətir
func (r *EmployeeRepository) GetCertificateByID(id int) (*models.EmployeeCertificate, error) {
	query := "SELECT id, employee_id, certificate_type, certificate_number, issued_by, issue_date, expiry_date, status, notes, created_at, updated_at FROM employee_certificates WHERE id = ?"
	var c models.EmployeeCertificate
	var issueDate, expiryDate sql.NullTime
	var certNumber, issuedBy, notes, status sql.NullString
	err := r.db.QueryRow(query, id).Scan(&c.ID, &c.EmployeeID, &c.CertificateType, &certNumber, &issuedBy, &issueDate, &expiryDate, &status, &notes, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return nil, err
	}
	if certNumber.Valid {
		c.CertificateNumber = certNumber.String
	}
	if issuedBy.Valid {
		c.IssuedBy = issuedBy.String
	}
	if notes.Valid {
		c.Notes = notes.String
	}
	if status.Valid {
		c.Status = status.String
	}
	if issueDate.Valid {
		c.IssueDate = &issueDate.Time
	}
	if expiryDate.Valid {
		c.ExpiryDate = &expiryDate.Time
	}
	return &c, nil
}
