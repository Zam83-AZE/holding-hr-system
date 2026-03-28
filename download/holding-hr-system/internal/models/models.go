package models

import (
	"time"
)

// Role - İstifadəçi rolları
type Role string

const (
	RoleAdmin        Role = "ADMIN"
	RoleHoldingHR    Role = "HOLDING_HR"
	RoleSubsidiaryHR Role = "SUBSIDIARY_HR"
)

// EmployeeStatus - İşçi statusları
type EmployeeStatus string

const (
	StatusCandidate  EmployeeStatus = "CANDIDATE"
	StatusActive     EmployeeStatus = "ACTIVE"
	StatusTerminated EmployeeStatus = "TERMINATED"
)

// Gender - Cins
type Gender string

const (
	GenderMale   Gender = "MALE"
	GenderFemale Gender = "FEMALE"
)

// Company - Şirkət
type Company struct {
	ID        int       `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	IsHolding bool      `json:"is_holding" db:"is_holding"`
	TaxID     string    `json:"tax_id" db:"tax_id"`
	Address   string    `json:"address" db:"address"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// Department - Departament
type Department struct {
	ID        int    `json:"id" db:"id"`
	CompanyID int    `json:"company_id" db:"company_id"`
	Name      string `json:"name" db:"name"`
}

// Position - Vəzifə
type Position struct {
	ID        int    `json:"id" db:"id"`
	CompanyID int    `json:"company_id" db:"company_id"`
	Name      string `json:"name" db:"name"`
}

// User - İstifadəçi
type User struct {
	ID           int       `json:"id" db:"id"`
	CompanyID    *int      `json:"company_id" db:"company_id"`
	FullName     string    `json:"full_name" db:"full_name"`
	Email        string    `json:"email" db:"email"`
	PasswordHash string    `json:"-" db:"password_hash"`
	Role         Role      `json:"role" db:"role"`
	IsActive     bool      `json:"is_active" db:"is_active"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}

// Employee - İşçi (əsas model)
type Employee struct {
	ID               int            `json:"id" db:"id"`
	CompanyID        int            `json:"company_id" db:"company_id"`
	FirstName        string         `json:"first_name" db:"first_name"`
	LastName         string         `json:"last_name" db:"last_name"`
	FatherName       string         `json:"father_name" db:"father_name"`
	FINCode          string         `json:"fin_code" db:"fin_code"`
	BirthDate        *time.Time     `json:"birth_date" db:"birth_date"`
	Gender           *Gender        `json:"gender" db:"gender"`
	PhotoPath        string         `json:"photo_path" db:"photo_path"`
	Phone            string         `json:"phone" db:"phone"`
	Email            string         `json:"email" db:"email"`
	Address          string         `json:"address" db:"address"`
	Status           EmployeeStatus `json:"status" db:"status"`
	DepartmentID     *int           `json:"department_id" db:"department_id"`
	PositionID       *int           `json:"position_id" db:"position_id"`
	HireDate         *time.Time     `json:"hire_date" db:"hire_date"`
	TerminationDate  *time.Time     `json:"termination_date" db:"termination_date"`
	TerminationReason string        `json:"termination_reason" db:"termination_reason"`
	CreatedAt        time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at" db:"updated_at"`

	// Əlaqəli məlumatlar
	CompanyName    string `json:"company_name" db:"company_name"`
	DepartmentName string `json:"department_name" db:"department_name"`
	PositionName   string `json:"position_name" db:"position_name"`
}

// EmployeeEducation - Təhsil
type EmployeeEducation struct {
	ID            int    `json:"id" db:"id"`
	EmployeeID    int    `json:"employee_id" db:"employee_id"`
	Institution   string `json:"institution" db:"institution"`
	Specialty     string `json:"specialty" db:"specialty"`
	Degree        string `json:"degree" db:"degree"`
	StartYear     int    `json:"start_year" db:"start_year"`
	EndYear       int    `json:"end_year" db:"end_year"`
	DiplomaNumber string `json:"diploma_number" db:"diploma_number"`
}

// EmployeeExperience - İş təcrübəsi
type EmployeeExperience struct {
	ID            int        `json:"id" db:"id"`
	EmployeeID    int        `json:"employee_id" db:"employee_id"`
	CompanyName   string     `json:"company_name" db:"company_name"`
	Position      string     `json:"position" db:"position"`
	StartDate     *time.Time `json:"start_date" db:"start_date"`
	EndDate       *time.Time `json:"end_date" db:"end_date"`
	LeavingReason string     `json:"leaving_reason" db:"leaving_reason"`
}

// EmployeeFamily - Ailə məlumatları
type EmployeeFamily struct {
	ID             int        `json:"id" db:"id"`
	EmployeeID     int        `json:"employee_id" db:"employee_id"`
	RelationType   string     `json:"relation_type" db:"relation_type"`
	FullName       string     `json:"full_name" db:"full_name"`
	BirthDate      *time.Time `json:"birth_date" db:"birth_date"`
	ContactNumber  string     `json:"contact_number" db:"contact_number"`
}

// EmployeeLifecycleLog - Yaşam dövrü tarixçəsi
type EmployeeLifecycleLog struct {
	ID               int       `json:"id" db:"id"`
	EmployeeID       int       `json:"employee_id" db:"employee_id"`
	UserID           int       `json:"user_id" db:"user_id"`
	OldStatus        string    `json:"old_status" db:"old_status"`
	NewStatus        string    `json:"new_status" db:"new_status"`
	EventDescription string    `json:"event_description" db:"event_description"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
}

// EmployeeFull - Tam işçi məlumatı (kartoçka üçün)
type EmployeeFull struct {
	Employee    Employee              `json:"employee"`
	Education   []EmployeeEducation   `json:"education"`
	Experience  []EmployeeExperience  `json:"experience"`
	Family      []EmployeeFamily      `json:"family"`
	Lifecycle   []EmployeeLifecycleLog `json:"lifecycle"`
}

// DashboardStats - Dashboard statistikası
type DashboardStats struct {
	TotalEmployees    int `json:"total_employees"`
	ActiveEmployees   int `json:"active_employees"`
	Candidates        int `json:"candidates"`
	Terminated        int `json:"terminated"`
	TotalCompanies    int `json:"total_companies"`
	ThisMonthHired    int `json:"this_month_hired"`
	ThisMonthTerminated int `json:"this_month_terminated"`
}

// Claims - JWT claims
type Claims struct {
	UserID    int    `json:"user_id"`
	CompanyID *int   `json:"company_id"`
	Email     string `json:"email"`
	Role      Role   `json:"role"`
}
