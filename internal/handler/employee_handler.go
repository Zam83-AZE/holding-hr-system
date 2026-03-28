package handler

import (
        "fmt"
        "holding-hr-system/internal/middleware"
        "holding-hr-system/internal/models"
        "holding-hr-system/internal/repository"
        "html/template"
        "net/http"
        "strconv"
        "time"
)

type EmployeeHandler struct {
        employeeRepo *repository.EmployeeRepository
        companyRepo  *repository.CompanyRepository
        deptRepo     *repository.DepartmentRepository
        posRepo      *repository.PositionRepository
        templates    *template.Template
}

func NewEmployeeHandler(
        employeeRepo *repository.EmployeeRepository,
        companyRepo *repository.CompanyRepository,
        deptRepo *repository.DepartmentRepository,
        posRepo *repository.PositionRepository,
        templates *template.Template,
) *EmployeeHandler {
        return &EmployeeHandler{
                employeeRepo: employeeRepo,
                companyRepo:  companyRepo,
                deptRepo:     deptRepo,
                posRepo:      posRepo,
                templates:    templates,
        }
}

// ShowDashboard - Ana səhifə
func (h *EmployeeHandler) ShowDashboard(w http.ResponseWriter, r *http.Request) {
        user := middleware.GetCurrentUser(r)
        companyFilter := middleware.GetCompanyFilter(user)

        stats, err := h.employeeRepo.GetStats(companyFilter)
        if err != nil {
                http.Error(w, "Statistika yüklənə bilmədi", http.StatusInternalServerError)
                return
        }

        companies, _ := h.companyRepo.GetAll()

        data := PageData{
                Title:     "Dashboard",
                User:      user,
                Stats:     stats,
                Companies: companies,
        }

        h.templates.ExecuteTemplate(w, "dashboard.html", data)
}

// ShowEmployees - İşçilər siyahısı
func (h *EmployeeHandler) ShowEmployees(w http.ResponseWriter, r *http.Request) {
        user := middleware.GetCurrentUser(r)
        companyFilter := middleware.GetCompanyFilter(user)

        status := models.EmployeeStatus(r.URL.Query().Get("status"))
        if status == "" {
                status = models.StatusActive
        }

        employees, err := h.employeeRepo.GetByStatus(companyFilter, status)
        if err != nil {
                http.Error(w, "İşçilər yüklənə bilmədi", http.StatusInternalServerError)
                return
        }

        companies, _ := h.companyRepo.GetAll()

        data := PageData{
                Title:     "Kadr Uçotu",
                User:      user,
                Employees: employees,
                Companies: companies,
                Status:    string(status),
        }

        h.templates.ExecuteTemplate(w, "employees.html", data)
}

// ShowEmployeeCard - İşçi kartoçkası
func (h *EmployeeHandler) ShowEmployeeCard(w http.ResponseWriter, r *http.Request) {
        user := middleware.GetCurrentUser(r)

        // ID-ni parse et
        idStr := r.URL.Query().Get("id")
        if idStr == "" {
                http.Error(w, "ID tələb olunur", http.StatusBadRequest)
                return
        }

        id, err := strconv.Atoi(idStr)
        if err != nil {
                http.Error(w, "Yanlış ID", http.StatusBadRequest)
                return
        }

        // İşçini gətir
        employee, err := h.employeeRepo.GetByID(id)
        if err != nil {
                http.Error(w, "İşçi tapılmadı", http.StatusNotFound)
                return
        }

        // İcazə yoxla
        if !middleware.CanAccessCompany(user, employee.CompanyID) {
                http.Error(w, "Səlahiyyətsiz", http.StatusForbidden)
                return
        }

        // Əlaqəli məlumatları gətir
        education, _ := h.employeeRepo.GetEducation(id)
        experience, _ := h.employeeRepo.GetExperience(id)
        family, _ := h.employeeRepo.GetFamily(id)
        lifecycle, _ := h.employeeRepo.GetLifecycleLogs(id)
	certificates, _ := h.employeeRepo.GetCertificatesByEmployee(id)

        // Departament və vəzifələri gətir
        departments, _ := h.deptRepo.GetByCompanyID(employee.CompanyID)
        positions, _ := h.posRepo.GetByCompanyID(employee.CompanyID)

        data := PageData{
                Title:       "İşçi Kartoçkası",
                User:        user,
                Employee:    employee,
                Education:   education,
                Experience:  experience,
                Family:      family,
                Lifecycle:   lifecycle,
		Certificates: certificates,
                Departments: departments,
                Positions:   positions,
        }

        h.templates.ExecuteTemplate(w, "employee_card.html", data)
}

// ShowNewEmployeeForm - Yeni işçi forması
func (h *EmployeeHandler) ShowNewEmployeeForm(w http.ResponseWriter, r *http.Request) {
        user := middleware.GetCurrentUser(r)

        var companies []models.Company
        var selectedCompanyID int

        if user.Role == models.RoleAdmin || user.Role == models.RoleHoldingHR {
                companies, _ = h.companyRepo.GetSubsidiaries()
                // URL-dən company_id parametri
                if cid := r.URL.Query().Get("company_id"); cid != "" {
                        selectedCompanyID, _ = strconv.Atoi(cid)
                } else if len(companies) > 0 {
                        selectedCompanyID = companies[0].ID
                }
        } else if user.CompanyID != nil {
                selectedCompanyID = *user.CompanyID
                company, _ := h.companyRepo.GetByID(selectedCompanyID)
                if company != nil {
                        companies = []models.Company{*company}
                }
        }

        departments, _ := h.deptRepo.GetByCompanyID(selectedCompanyID)
        positions, _ := h.posRepo.GetByCompanyID(selectedCompanyID)

        data := PageData{
                Title:            "Yeni İşçi",
                User:             user,
                Companies:        companies,
                Departments:      departments,
                Positions:        positions,
                SelectedCompany:  selectedCompanyID,
                IsNew:            true,
        }

        h.templates.ExecuteTemplate(w, "employee_form.html", data)
}

// CreateEmployee - Yeni işçi yarat
func (h *EmployeeHandler) CreateEmployee(w http.ResponseWriter, r *http.Request) {
        user := middleware.GetCurrentUser(r)

        // Form məlumatlarını al
        companyID, _ := strconv.Atoi(r.FormValue("company_id"))
        firstName := r.FormValue("first_name")
        lastName := r.FormValue("last_name")
        fatherName := r.FormValue("father_name")
        finCode := r.FormValue("fin_code")
        phone := r.FormValue("phone")
        email := r.FormValue("email")
        address := r.FormValue("address")

        // İcazə yoxla
        if !middleware.CanAccessCompany(user, companyID) {
                http.Error(w, "Səlahiyyətsiz", http.StatusForbidden)
                return
        }

        employee := &models.Employee{
                CompanyID:  companyID,
                FirstName:  firstName,
                LastName:   lastName,
                FatherName: fatherName,
                FINCode:    finCode,
                Phone:      phone,
                Email:      email,
                Address:    address,
                Status:     models.StatusCandidate,
        }

        // Doğum tarixi
        if birthDate := r.FormValue("birth_date"); birthDate != "" {
                t, err := time.Parse("2006-01-02", birthDate)
                if err == nil {
                        employee.BirthDate = &t
                }
        }

        // Cins
        if gender := r.FormValue("gender"); gender != "" {
                g := models.Gender(gender)
                employee.Gender = &g
        }

        if err := h.employeeRepo.Create(employee); err != nil {
                http.Error(w, "İşçi yaradıla bilmədi: "+err.Error(), http.StatusInternalServerError)
                return
        }

        // Tarixçəyə əlavə et
        log := &models.EmployeeLifecycleLog{
                EmployeeID:       employee.ID,
                UserID:           user.UserID,
                OldStatus:        "",
                NewStatus:        string(models.StatusCandidate),
                EventDescription: "Namizəd kimi sistemə əlavə edildi",
        }
        h.employeeRepo.AddLifecycleLog(log)

        http.Redirect(w, r, "/employees?status=CANDIDATE", http.StatusSeeOther)
}

// UpdateEmployee - İşçi məlumatlarını yenilə
func (h *EmployeeHandler) UpdateEmployee(w http.ResponseWriter, r *http.Request) {
        user := middleware.GetCurrentUser(r)

        id, _ := strconv.Atoi(r.FormValue("id"))
        employee, err := h.employeeRepo.GetByID(id)
        if err != nil {
                http.Error(w, "İşçi tapılmadı", http.StatusNotFound)
                return
        }

        // İcazə yoxla
        if !middleware.CanAccessCompany(user, employee.CompanyID) {
                http.Error(w, "Səlahiyyətsiz", http.StatusForbidden)
                return
        }

        // Form məlumatlarını al
        employee.FirstName = r.FormValue("first_name")
        employee.LastName = r.FormValue("last_name")
        employee.FatherName = r.FormValue("father_name")
        employee.FINCode = r.FormValue("fin_code")
        employee.Phone = r.FormValue("phone")
        employee.Email = r.FormValue("email")
        employee.Address = r.FormValue("address")

        // Doğum tarixi
        if birthDate := r.FormValue("birth_date"); birthDate != "" {
                t, err := time.Parse("2006-01-02", birthDate)
                if err == nil {
                        employee.BirthDate = &t
                }
        }

        // Cins
        if gender := r.FormValue("gender"); gender != "" {
                g := models.Gender(gender)
                employee.Gender = &g
        }

        // Departament və vəzifə
        if deptID := r.FormValue("department_id"); deptID != "" {
                id, _ := strconv.Atoi(deptID)
                employee.DepartmentID = &id
        }
        if posID := r.FormValue("position_id"); posID != "" {
                id, _ := strconv.Atoi(posID)
                employee.PositionID = &id
        }

        if err := h.employeeRepo.Update(employee); err != nil {
                http.Error(w, "Məlumatlar yenilənə bilmədi", http.StatusInternalServerError)
                return
        }

        http.Redirect(w, r, "/employee/card?id="+strconv.Itoa(employee.ID), http.StatusSeeOther)
}

// HireEmployee - Namizədi işə qəbul et
func (h *EmployeeHandler) HireEmployee(w http.ResponseWriter, r *http.Request) {
        user := middleware.GetCurrentUser(r)

        id, _ := strconv.Atoi(r.FormValue("id"))
        employee, err := h.employeeRepo.GetByID(id)
        if err != nil {
                http.Error(w, "İşçi tapılmadı", http.StatusNotFound)
                return
        }

        if !middleware.CanAccessCompany(user, employee.CompanyID) {
                http.Error(w, "Səlahiyyətsiz", http.StatusForbidden)
                return
        }

        // Departament və vəzifə
        deptID, _ := strconv.Atoi(r.FormValue("department_id"))
        posID, _ := strconv.Atoi(r.FormValue("position_id"))

        if deptID > 0 {
                employee.DepartmentID = &deptID
        }
        if posID > 0 {
                employee.PositionID = &posID
        }

        // İşə qəbul tarixi
        hireDateStr := r.FormValue("hire_date")
        if hireDateStr == "" {
                hireDateStr = time.Now().Format("2006-01-02")
        }
        hireDate, _ := time.Parse("2006-01-02", hireDateStr)
        employee.HireDate = &hireDate

        // Statusu dəyiş
        if err := h.employeeRepo.UpdateStatus(id, models.StatusActive, &hireDate, nil, ""); err != nil {
                http.Error(w, "Status dəyişdirilə bilmədi", http.StatusInternalServerError)
                return
        }

        // Tarixçə
        log := &models.EmployeeLifecycleLog{
                EmployeeID:       id,
                UserID:           user.UserID,
                OldStatus:        string(employee.Status),
                NewStatus:        string(models.StatusActive),
                EventDescription: fmt.Sprintf("İşə qəbul edildi. Tarix: %s", hireDate.Format("02.01.2006")),
        }
        h.employeeRepo.AddLifecycleLog(log)

        http.Redirect(w, r, "/employees?status=ACTIVE", http.StatusSeeOther)
}

// TerminateEmployee - İşçini işdən çıxar
func (h *EmployeeHandler) TerminateEmployee(w http.ResponseWriter, r *http.Request) {
        user := middleware.GetCurrentUser(r)

        id, _ := strconv.Atoi(r.FormValue("id"))
        employee, _ := h.employeeRepo.GetByID(id)

        if !middleware.CanAccessCompany(user, employee.CompanyID) {
                http.Error(w, "Səlahiyyətsiz", http.StatusForbidden)
                return
        }

        terminationReason := r.FormValue("termination_reason")
        terminationDateStr := r.FormValue("termination_date")
        if terminationDateStr == "" {
                terminationDateStr = time.Now().Format("2006-01-02")
        }
        terminationDate, _ := time.Parse("2006-01-02", terminationDateStr)

        if err := h.employeeRepo.UpdateStatus(id, models.StatusTerminated, nil, &terminationDate, terminationReason); err != nil {
                http.Error(w, "Status dəyişdirilə bilmədi", http.StatusInternalServerError)
                return
        }

        log := &models.EmployeeLifecycleLog{
                EmployeeID:       id,
                UserID:           user.UserID,
                OldStatus:        string(employee.Status),
                NewStatus:        string(models.StatusTerminated),
                EventDescription: fmt.Sprintf("İşdən çıxarıldı. Səbəb: %s", terminationReason),
        }
        h.employeeRepo.AddLifecycleLog(log)

        http.Redirect(w, r, "/employees?status=TERMINATED", http.StatusSeeOther)
}

// ReactivateEmployee - İşçini bərpa et
func (h *EmployeeHandler) ReactivateEmployee(w http.ResponseWriter, r *http.Request) {
        user := middleware.GetCurrentUser(r)

        id, _ := strconv.Atoi(r.FormValue("id"))
        employee, _ := h.employeeRepo.GetByID(id)

        if !middleware.CanAccessCompany(user, employee.CompanyID) {
                http.Error(w, "Səlahiyyətsiz", http.StatusForbidden)
                return
        }

        hireDate := time.Now()

        if err := h.employeeRepo.UpdateStatus(id, models.StatusActive, &hireDate, nil, ""); err != nil {
                http.Error(w, "Status dəyişdirilə bilmədi", http.StatusInternalServerError)
                return
        }

        log := &models.EmployeeLifecycleLog{
                EmployeeID:       id,
                UserID:           user.UserID,
                OldStatus:        string(models.StatusTerminated),
                NewStatus:        string(models.StatusActive),
                EventDescription: "Yenidən işə qəbul edildi",
        }
        h.employeeRepo.AddLifecycleLog(log)

        http.Redirect(w, r, "/employees?status=ACTIVE", http.StatusSeeOther)
}

// SearchEmployees - Axtarış
func (h *EmployeeHandler) SearchEmployees(w http.ResponseWriter, r *http.Request) {
        user := middleware.GetCurrentUser(r)
        companyFilter := middleware.GetCompanyFilter(user)

        query := r.URL.Query().Get("q")
        if query == "" {
                http.Redirect(w, r, "/employees", http.StatusSeeOther)
                return
        }

        employees, err := h.employeeRepo.Search(companyFilter, query)
        if err != nil {
                http.Error(w, "Axtarış xətası", http.StatusInternalServerError)
                return
        }

        data := PageData{
                Title:     "Axtarış Nəticələri",
                User:      user,
                Employees: employees,
                Query:     query,
        }

        h.templates.ExecuteTemplate(w, "employees.html", data)
}

// AddEducation - Təhsil əlavə et
func (h *EmployeeHandler) AddEducation(w http.ResponseWriter, r *http.Request) {
        employeeID, _ := strconv.Atoi(r.FormValue("employee_id"))

        edu := &models.EmployeeEducation{
                EmployeeID:    employeeID,
                Institution:   r.FormValue("institution"),
                Specialty:     r.FormValue("specialty"),
                Degree:        r.FormValue("degree"),
                DiplomaNumber: r.FormValue("diploma_number"),
        }

        edu.StartYear, _ = strconv.Atoi(r.FormValue("start_year"))
        edu.EndYear, _ = strconv.Atoi(r.FormValue("end_year"))

        h.employeeRepo.AddEducation(edu)
        http.Redirect(w, r, "/employee/card?id="+strconv.Itoa(employeeID), http.StatusSeeOther)
}

// AddExperience - İş təcrübəsi əlavə et
func (h *EmployeeHandler) AddExperience(w http.ResponseWriter, r *http.Request) {
        employeeID, _ := strconv.Atoi(r.FormValue("employee_id"))

        exp := &models.EmployeeExperience{
                EmployeeID:    employeeID,
                CompanyName:   r.FormValue("company_name"),
                Position:      r.FormValue("position"),
                LeavingReason: r.FormValue("leaving_reason"),
        }

        if startDate := r.FormValue("start_date"); startDate != "" {
                t, _ := time.Parse("2006-01-02", startDate)
                exp.StartDate = &t
        }
        if endDate := r.FormValue("end_date"); endDate != "" {
                t, _ := time.Parse("2006-01-02", endDate)
                exp.EndDate = &t
        }

        h.employeeRepo.AddExperience(exp)
        http.Redirect(w, r, "/employee/card?id="+strconv.Itoa(employeeID), http.StatusSeeOther)
}

// AddFamily - Ailə üzvü əlavə et
func (h *EmployeeHandler) AddFamily(w http.ResponseWriter, r *http.Request) {
        employeeID, _ := strconv.Atoi(r.FormValue("employee_id"))

        fam := &models.EmployeeFamily{
                EmployeeID:   employeeID,
                RelationType: r.FormValue("relation_type"),
                FullName:     r.FormValue("full_name"),
        }

        if birthDate := r.FormValue("birth_date"); birthDate != "" {
                t, _ := time.Parse("2006-01-02", birthDate)
                fam.BirthDate = &t
        }
        fam.ContactNumber = r.FormValue("contact_number")

        h.employeeRepo.AddFamily(fam)
        http.Redirect(w, r, "/employee/card?id="+strconv.Itoa(employeeID), http.StatusSeeOther)
}

// ========== EDUCATION CRUD ==========

// UpdateEducation - Təhsil yenilə
func (h *EmployeeHandler) UpdateEducation(w http.ResponseWriter, r *http.Request) {
        id, _ := strconv.Atoi(r.FormValue("id"))
        
        edu, err := h.employeeRepo.GetEducationByID(id)
        if err != nil {
                http.Error(w, "Təhsil tapılmadı", http.StatusNotFound)
                return
        }

        edu.Institution = r.FormValue("institution")
        edu.Specialty = r.FormValue("specialty")
        edu.Degree = r.FormValue("degree")
        edu.DiplomaNumber = r.FormValue("diploma_number")
        edu.StartYear, _ = strconv.Atoi(r.FormValue("start_year"))
        edu.EndYear, _ = strconv.Atoi(r.FormValue("end_year"))

        h.employeeRepo.UpdateEducation(edu)
        http.Redirect(w, r, "/employee/card?id="+strconv.Itoa(edu.EmployeeID), http.StatusSeeOther)
}

// DeleteEducation - Təhsil sil
func (h *EmployeeHandler) DeleteEducation(w http.ResponseWriter, r *http.Request) {
        id, _ := strconv.Atoi(r.FormValue("id"))
        
        edu, err := h.employeeRepo.GetEducationByID(id)
        if err != nil {
                http.Error(w, "Təhsil tapılmadı", http.StatusNotFound)
                return
        }

        h.employeeRepo.DeleteEducation(id)
        http.Redirect(w, r, "/employee/card?id="+strconv.Itoa(edu.EmployeeID), http.StatusSeeOther)
}

// ========== EXPERIENCE CRUD ==========

// UpdateExperience - İş təcrübəsi yenilə
func (h *EmployeeHandler) UpdateExperience(w http.ResponseWriter, r *http.Request) {
        id, _ := strconv.Atoi(r.FormValue("id"))
        
        exp, err := h.employeeRepo.GetExperienceByID(id)
        if err != nil {
                http.Error(w, "Təcrübə tapılmadı", http.StatusNotFound)
                return
        }

        exp.CompanyName = r.FormValue("company_name")
        exp.Position = r.FormValue("position")
        exp.LeavingReason = r.FormValue("leaving_reason")

        if startDate := r.FormValue("start_date"); startDate != "" {
                t, _ := time.Parse("2006-01-02", startDate)
                exp.StartDate = &t
        }
        if endDate := r.FormValue("end_date"); endDate != "" {
                t, _ := time.Parse("2006-01-02", endDate)
                exp.EndDate = &t
        }

        h.employeeRepo.UpdateExperience(exp)
        http.Redirect(w, r, "/employee/card?id="+strconv.Itoa(exp.EmployeeID), http.StatusSeeOther)
}

// DeleteExperience - İş təcrübəsi sil
func (h *EmployeeHandler) DeleteExperience(w http.ResponseWriter, r *http.Request) {
        id, _ := strconv.Atoi(r.FormValue("id"))
        
        exp, err := h.employeeRepo.GetExperienceByID(id)
        if err != nil {
                http.Error(w, "Təcrübə tapılmadı", http.StatusNotFound)
                return
        }

        h.employeeRepo.DeleteExperience(id)
        http.Redirect(w, r, "/employee/card?id="+strconv.Itoa(exp.EmployeeID), http.StatusSeeOther)
}

// ========== FAMILY CRUD ==========

// UpdateFamily - Ailə üzvü yenilə
func (h *EmployeeHandler) UpdateFamily(w http.ResponseWriter, r *http.Request) {
        id, _ := strconv.Atoi(r.FormValue("id"))
        
        fam, err := h.employeeRepo.GetFamilyByID(id)
        if err != nil {
                http.Error(w, "Ailə üzvü tapılmadı", http.StatusNotFound)
                return
        }

        fam.RelationType = r.FormValue("relation_type")
        fam.FullName = r.FormValue("full_name")
        fam.ContactNumber = r.FormValue("contact_number")

        if birthDate := r.FormValue("birth_date"); birthDate != "" {
                t, _ := time.Parse("2006-01-02", birthDate)
                fam.BirthDate = &t
        } else {
                fam.BirthDate = nil
        }

        h.employeeRepo.UpdateFamily(fam)
        http.Redirect(w, r, "/employee/card?id="+strconv.Itoa(fam.EmployeeID), http.StatusSeeOther)
}

// DeleteFamily - Ailə üzvü sil
func (h *EmployeeHandler) DeleteFamily(w http.ResponseWriter, r *http.Request) {
        id, _ := strconv.Atoi(r.FormValue("id"))
        
        fam, err := h.employeeRepo.GetFamilyByID(id)
        if err != nil {
                http.Error(w, "Ailə üzvü tapılmadı", http.StatusNotFound)
                return
        }

        h.employeeRepo.DeleteFamily(id)
        http.Redirect(w, r, "/employee/card?id="+strconv.Itoa(fam.EmployeeID), http.StatusSeeOther)
}

// DeleteEmployee - İşçi sil
func (h *EmployeeHandler) DeleteEmployee(w http.ResponseWriter, r *http.Request) {
        user := middleware.GetCurrentUser(r)

        id, _ := strconv.Atoi(r.FormValue("id"))
        employee, _ := h.employeeRepo.GetByID(id)

        if !middleware.CanAccessCompany(user, employee.CompanyID) {
                http.Error(w, "Səlahiyyətsiz", http.StatusForbidden)
                return
        }

        // Yalnız Admin silə bilər
        if user.Role != models.RoleAdmin {
                http.Error(w, "Səlahiyyətsiz", http.StatusForbidden)
                return
        }

        h.employeeRepo.Delete(id)
        http.Redirect(w, r, "/employees", http.StatusSeeOther)
}

// API: GetDepartmentsByCompany
func (h *EmployeeHandler) GetDepartmentsByCompany(w http.ResponseWriter, r *http.Request) {
        companyID, _ := strconv.Atoi(r.URL.Query().Get("company_id"))
        departments, _ := h.deptRepo.GetByCompanyID(companyID)

        w.Header().Set("Content-Type", "application/json")
        fmt.Fprintf(w, "[")
        for i, dept := range departments {
                if i > 0 {
                        fmt.Fprintf(w, ",")
                }
                fmt.Fprintf(w, `{"id":%d,"name":"%s"}`, dept.ID, template.HTMLEscapeString(dept.Name))
        }
        fmt.Fprintf(w, "]")
}

// API: GetPositionsByCompany
func (h *EmployeeHandler) GetPositionsByCompany(w http.ResponseWriter, r *http.Request) {
        companyID, _ := strconv.Atoi(r.URL.Query().Get("company_id"))
        positions, _ := h.posRepo.GetByCompanyID(companyID)

        w.Header().Set("Content-Type", "application/json")
        fmt.Fprintf(w, "[")
        for i, pos := range positions {
                if i > 0 {
                        fmt.Fprintf(w, ",")
                }
                fmt.Fprintf(w, `{"id":%d,"name":"%s"}`, pos.ID, template.HTMLEscapeString(pos.Name))
        }
        fmt.Fprintf(w, "]")
}

// AddCertificate - Sertifikat əlavə et
func (h *EmployeeHandler) AddCertificate(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	employeeID, _ := strconv.Atoi(r.FormValue("employee_id"))
	cert := &models.EmployeeCertificate{
		EmployeeID:       employeeID,
		CertificateType:  r.FormValue("certificate_type"),
		CertificateNumber: r.FormValue("certificate_number"),
		IssuedBy:         r.FormValue("issued_by"),
		Notes:            r.FormValue("notes"),
	}
	if d := r.FormValue("issue_date"); d != "" {
		t, _ := time.Parse("2006-01-02", d)
		cert.IssueDate = &t
	}
	if d := r.FormValue("expiry_date"); d != "" {
		t, _ := time.Parse("2006-01-02", d)
		cert.ExpiryDate = &t
	}
	h.employeeRepo.AddCertificate(cert)
	http.Redirect(w, r, fmt.Sprintf("/employee/card?id=%d&tab=certificates", employeeID), http.StatusSeeOther)
}

// UpdateCertificate - Sertifikat yenilə
func (h *EmployeeHandler) UpdateCertificate(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	id, _ := strconv.Atoi(r.FormValue("id"))
	employeeID, _ := strconv.Atoi(r.FormValue("employee_id"))
	cert, _ := h.employeeRepo.GetCertificateByID(id)
	if cert == nil {
		http.Error(w, "Certificate not found", http.StatusNotFound)
		return
	}
	cert.CertificateType = r.FormValue("certificate_type")
	cert.CertificateNumber = r.FormValue("certificate_number")
	cert.IssuedBy = r.FormValue("issued_by")
	cert.Notes = r.FormValue("notes")
	if d := r.FormValue("issue_date"); d != "" {
		t, _ := time.Parse("2006-01-02", d)
		cert.IssueDate = &t
	}
	if d := r.FormValue("expiry_date"); d != "" {
		t, _ := time.Parse("2006-01-02", d)
		cert.ExpiryDate = &t
	}
	h.employeeRepo.UpdateCertificate(cert)
	http.Redirect(w, r, fmt.Sprintf("/employee/card?id=%d&tab=certificates", employeeID), http.StatusSeeOther)
}

// DeleteCertificate - Sertifikat sil
func (h *EmployeeHandler) DeleteCertificate(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	id, _ := strconv.Atoi(r.FormValue("id"))
	employeeID, _ := strconv.Atoi(r.FormValue("employee_id"))
	h.employeeRepo.DeleteCertificate(id)
	http.Redirect(w, r, fmt.Sprintf("/employee/card?id=%d&tab=certificates", employeeID), http.StatusSeeOther)
}

// GetWorkLocations - API: İş yerlərini JSON kimi qaytar
func (h *EmployeeHandler) GetWorkLocations(w http.ResponseWriter, r *http.Request) {
	companyID, _ := strconv.Atoi(r.URL.Query().Get("company_id"))
	locations, err := h.employeeRepo.GetWorkLocationsByCompany(companyID)
	if err != nil {
		http.Error(w, "Error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	result := "["
	for i, loc := range locations {
		if i > 0 {
			result += ","
		}
		result += fmt.Sprintf(`{"id":%d,"name":"%s","type":"%s"}`, loc.ID, loc.Name, loc.Type)
	}
	result += "]"
	w.Write([]byte(result))
}
