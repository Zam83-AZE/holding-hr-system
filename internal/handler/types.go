package handler

import (
        "holding-hr-system/internal/models"
)

// PageData - Template üçün data strukturu
type PageData struct {
        Title            string
        User             *models.Claims
        Stats            *models.DashboardStats
        Employees        []models.Employee
        Employee         *models.Employee
        Education        []models.EmployeeEducation
        Experience       []models.EmployeeExperience
        Family           []models.EmployeeFamily
        Lifecycle        []models.EmployeeLifecycleLog
        Certificates     []models.EmployeeCertificate
        Companies        []models.Company
        Users            []models.User
        Departments      []models.Department
        Positions        []models.Position
        Status              string
        Query               string
        SelectedCompany     int
        SelectedCompanyName string
        IsNew               bool
        Error            string
        Success          string
}
