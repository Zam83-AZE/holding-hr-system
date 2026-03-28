-- Work Locations table
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

-- Employee Certificates table
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

-- Add work_location_id and uniform_size columns to employees table (idempotent)
ALTER TABLE employees ADD COLUMN IF NOT EXISTS work_location_id INT NULL AFTER position_id;
ALTER TABLE employees ADD COLUMN IF NOT EXISTS uniform_size VARCHAR(20) NULL AFTER work_location_id;

-- Seed Work Locations
INSERT IGNORE INTO work_locations (company_id, name, address, type) VALUES
-- City Service Company restaurants (5)
(5, 'Nar & Sharab', 'Baku, Nizami küç., 23', 'RESTAURANT'),
(5, 'Senator Hall', 'Baku, 28 May küç., 45', 'RESTAURANT'),
(5, 'Maxim Hall', 'Baku, Füzuli küç., 67', 'RESTAURANT'),
(5, 'Mangal Steak House', 'Baku, Bulbul prospekti, 12', 'RESTAURANT'),
(5, 'Zuğulba Saray', 'Baku, Zuğulba yolu, 89', 'RESTAURANT'),
-- Sapphire Hotels (3)
(4, 'Sapphire City Hotel', 'Baku, Azadlıq prospekti, 34', 'HOTEL'),
(4, 'Sapphire Resort', 'Abşeron, Şıxov yolu', 'HOTEL'),
(4, 'Sapphire Business Center', 'Baku, Neftçiler prospekti, 56', 'OFFICE'),
-- TEZ Logistics (2)
(3, 'TEZ Warehouse 1', 'Baku, Sumqayıt yolu, 15', 'WAREHOUSE'),
(3, 'TEZ Terminal', 'Baku, Dənizkənarı bulvar', 'WAREHOUSE'),
-- Azman Construction (2)
(2, 'Ofis Mərkəzi', 'Baku, Səbail rayonu, 78', 'OFFICE'),
(2, 'Tikinti Sahəsi 1', 'Xırdalan, Şəhər qəsəbəsi', 'OTHER'),
-- EcoProd Azerbaijan (2)
(6, 'EcoProd Kənd Təsərrüfatı', 'Samux rayonu', 'FARM'),
(6, 'EcoProd Emalat Zavodu', 'Gəncə, Sənaye zonası', 'FACTORY'),
-- Mangal MMC (2)
(7, 'Mangal MMC Zavodu', 'Baku, Balaxanı sənaye zonası', 'FACTORY'),
(7, 'Mangal MMC Ofis', 'Baku, Nəsimi rayonu, 90', 'OFFICE');
