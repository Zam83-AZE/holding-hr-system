package handler

import (
	"holding-hr-system/internal/middleware"
	"holding-hr-system/internal/models"
	"holding-hr-system/internal/repository"
	"net/http"
	"strconv"
)

type CompanyHandler struct {
	companyRepo *repository.CompanyRepository
	deptRepo    *repository.DepartmentRepository
	posRepo     *repository.PositionRepository
	templates   *template.Template
}

func NewCompanyHandler(
	companyRepo *repository.CompanyRepository,
	deptRepo *repository.DepartmentRepository,
	posRepo *repository.PositionRepository,
	templates *template.Template,
) *CompanyHandler {
	return &CompanyHandler{
		companyRepo: companyRepo,
		deptRepo:    deptRepo,
		posRepo:     posRepo,
		templates:   templates,
	}
}

// ShowStructure - Struktur səhifəsi
func (h *CompanyHandler) ShowStructure(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetCurrentUser(r)

	companies, _ := h.companyRepo.GetAll()

	// İlk şirkətin məlumatlarını gətir
	var selectedCompanyID int
	var departments []models.Department
	var positions []models.Position

	if len(companies) > 0 {
		if user.CompanyID != nil {
			selectedCompanyID = *user.CompanyID
		} else {
			selectedCompanyID = companies[0].ID
		}
		departments, _ = h.deptRepo.GetByCompanyID(selectedCompanyID)
		positions, _ = h.posRepo.GetByCompanyID(selectedCompanyID)
	}

	data := PageData{
		Title:       "Struktur",
		User:        user,
		Companies:   companies,
		Departments: departments,
		Positions:   positions,
		SelectedCompany: selectedCompanyID,
	}

	h.templates.ExecuteTemplate(w, "structure.html", data)
}

// ShowCompanyStructure - Şirkət strukturunu göstər (AJAX)
func (h *CompanyHandler) ShowCompanyStructure(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetCurrentUser(r)
	companyID, _ := strconv.Atoi(r.URL.Query().Get("company_id"))

	// İcazə yoxla
	if !middleware.CanAccessCompany(user, companyID) {
		http.Error(w, "Səlahiyyətsiz", http.StatusForbidden)
		return
	}

	departments, _ := h.deptRepo.GetByCompanyID(companyID)
	positions, _ := h.posRepo.GetByCompanyID(companyID)

	data := PageData{
		User:        user,
		Departments: departments,
		Positions:   positions,
		SelectedCompany: companyID,
	}

	h.templates.ExecuteTemplate(w, "structure_list", data)
}

// CreateDepartment - Departament yarat
func (h *CompanyHandler) CreateDepartment(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetCurrentUser(r)

	if user.Role == models.RoleSubsidiaryHR {
		http.Error(w, "Səlahiyyətsiz", http.StatusForbidden)
		return
	}

	companyID, _ := strconv.Atoi(r.FormValue("company_id"))
	name := r.FormValue("name")

	dept := &models.Department{
		CompanyID: companyID,
		Name:      name,
	}

	if err := h.deptRepo.Create(dept); err != nil {
		http.Error(w, "Departament yaradıla bilmədi", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/structure?company_id="+strconv.Itoa(companyID), http.StatusSeeOther)
}

// CreatePosition - Vəzifə yarat
func (h *CompanyHandler) CreatePosition(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetCurrentUser(r)

	if user.Role == models.RoleSubsidiaryHR {
		http.Error(w, "Səlahiyyətsiz", http.StatusForbidden)
		return
	}

	companyID, _ := strconv.Atoi(r.FormValue("company_id"))
	name := r.FormValue("name")

	pos := &models.Position{
		CompanyID: companyID,
		Name:      name,
	}

	if err := h.posRepo.Create(pos); err != nil {
		http.Error(w, "Vəzifə yaradıla bilmədi", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/structure?company_id="+strconv.Itoa(companyID), http.StatusSeeOther)
}

// DeleteDepartment - Departament sil
func (h *CompanyHandler) DeleteDepartment(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetCurrentUser(r)

	if user.Role != models.RoleAdmin {
		http.Error(w, "Səlahiyyətsiz", http.StatusForbidden)
		return
	}

	id, _ := strconv.Atoi(r.FormValue("id"))
	companyID := r.FormValue("company_id")

	h.deptRepo.Delete(id)
	http.Redirect(w, r, "/structure?company_id="+companyID, http.StatusSeeOther)
}

// DeletePosition - Vəzifə sil
func (h *CompanyHandler) DeletePosition(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetCurrentUser(r)

	if user.Role != models.RoleAdmin {
		http.Error(w, "Səlahiyyətsiz", http.StatusForbidden)
		return
	}

	id, _ := strconv.Atoi(r.FormValue("id"))
	companyID := r.FormValue("company_id")

	h.posRepo.Delete(id)
	http.Redirect(w, r, "/structure?company_id="+companyID, http.StatusSeeOther)
}

// ShowSettings - Ayarlar səhifəsi
func (h *CompanyHandler) ShowSettings(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetCurrentUser(r)

	// Yalnız Admin görə bilər
	if user.Role != models.RoleAdmin {
		http.Error(w, "Səlahiyyətsiz", http.StatusForbidden)
		return
	}

	companies, _ := h.companyRepo.GetAll()

	data := PageData{
		Title:     "Ayarlar",
		User:      user,
		Companies: companies,
	}

	h.templates.ExecuteTemplate(w, "settings.html", data)
}

// CreateCompany - Şirkət yarat
func (h *CompanyHandler) CreateCompany(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetCurrentUser(r)

	if user.Role != models.RoleAdmin {
		http.Error(w, "Səlahiyyətsiz", http.StatusForbidden)
		return
	}

	name := r.FormValue("name")
	taxID := r.FormValue("tax_id")
	address := r.FormValue("address")
	isHolding := r.FormValue("is_holding") == "on"

	company := &models.Company{
		Name:      name,
		TaxID:     taxID,
		Address:   address,
		IsHolding: isHolding,
	}

	if err := h.companyRepo.Create(company); err != nil {
		http.Error(w, "Şirkət yaradıla bilmədi", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/settings", http.StatusSeeOther)
}

// DeleteCompany - Şirkət sil
func (h *CompanyHandler) DeleteCompany(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetCurrentUser(r)

	if user.Role != models.RoleAdmin {
		http.Error(w, "Səlahiyyətsiz", http.StatusForbidden)
		return
	}

	id, _ := strconv.Atoi(r.FormValue("id"))
	h.companyRepo.Delete(id)
	http.Redirect(w, r, "/settings", http.StatusSeeOther)
}
