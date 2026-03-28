-- Migration 05: Safety net - ensure all tables and columns exist
-- This migration is designed to be safe on both fresh and existing databases

SET NAMES utf8mb4;

-- Ensure work_locations table exists
CREATE TABLE IF NOT EXISTS work_locations (
    id INT AUTO_INCREMENT PRIMARY KEY,
    company_id INT NOT NULL,
    name VARCHAR(255) NOT NULL,
    address VARCHAR(500),
    type ENUM('RESTAURANT', 'HOTEL', 'OFFICE', 'WAREHOUSE', 'FACTORY', 'FARM', 'OTHER') DEFAULT 'OTHER',
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (company_id) REFERENCES companies(id) ON DELETE CASCADE,
    INDEX idx_company (company_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Ensure employee_certificates table exists
CREATE TABLE IF NOT EXISTS employee_certificates (
    id INT AUTO_INCREMENT PRIMARY KEY,
    employee_id INT NOT NULL,
    certificate_type VARCHAR(100) NOT NULL,
    certificate_number VARCHAR(255),
    issued_by VARCHAR(255),
    issue_date DATE,
    expiry_date DATE,
    status ENUM('ACTIVE', 'EXPIRED', 'REVOKED') DEFAULT 'ACTIVE',
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (employee_id) REFERENCES employees(id) ON DELETE CASCADE,
    INDEX idx_employee (employee_id),
    INDEX idx_expiry (expiry_date)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
