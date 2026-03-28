-- Holding HR System Database Schema
-- MariaDB 10.11+

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- --------------------------------------------------------
-- Şirkətlər (Holding və Alt Şirkətlər)
-- --------------------------------------------------------
CREATE TABLE IF NOT EXISTS companies (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    is_holding BOOLEAN DEFAULT FALSE,
    tax_id VARCHAR(50) UNIQUE,
    address TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- --------------------------------------------------------
-- İstifadəçilər (HR-lar)
-- --------------------------------------------------------
CREATE TABLE IF NOT EXISTS users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    company_id INT NULL,
    full_name VARCHAR(150) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role ENUM('HOLDING_HR', 'SUBSIDIARY_HR', 'ADMIN') NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (company_id) REFERENCES companies(id) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- --------------------------------------------------------
-- Departamentlər
-- --------------------------------------------------------
CREATE TABLE IF NOT EXISTS departments (
    id INT AUTO_INCREMENT PRIMARY KEY,
    company_id INT NOT NULL,
    name VARCHAR(150) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (company_id) REFERENCES companies(id) ON DELETE CASCADE,
    UNIQUE KEY unique_dept_company (company_id, name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- --------------------------------------------------------
-- Vəzifələr
-- --------------------------------------------------------
CREATE TABLE IF NOT EXISTS positions (
    id INT AUTO_INCREMENT PRIMARY KEY,
    company_id INT NOT NULL,
    name VARCHAR(150) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (company_id) REFERENCES companies(id) ON DELETE CASCADE,
    UNIQUE KEY unique_pos_company (company_id, name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- --------------------------------------------------------
-- İşçilər (Əsas cədvəl)
-- --------------------------------------------------------
CREATE TABLE IF NOT EXISTS employees (
    id INT AUTO_INCREMENT PRIMARY KEY,
    company_id INT NOT NULL,

    -- Şəxsi Məlumatlar (Blok A)
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    father_name VARCHAR(100),
    fin_code VARCHAR(7) UNIQUE NOT NULL,
    birth_date DATE,
    gender ENUM('MALE', 'FEMALE'),
    photo_path VARCHAR(255),

    -- Əlaqə məlumatları
    phone VARCHAR(20),
    email VARCHAR(100),
    address TEXT,

    -- 3 Əsas Qrup (Status)
    status ENUM('CANDIDATE', 'ACTIVE', 'TERMINATED') DEFAULT 'CANDIDATE',

    -- Korporativ Məlumatlar (Blok E)
    department_id INT NULL,
    position_id INT NULL,
    hire_date DATE NULL,
    termination_date DATE NULL,
    termination_reason TEXT,

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    FOREIGN KEY (company_id) REFERENCES companies(id) ON DELETE CASCADE,
    FOREIGN KEY (department_id) REFERENCES departments(id) ON DELETE SET NULL,
    FOREIGN KEY (position_id) REFERENCES positions(id) ON DELETE SET NULL,

    INDEX idx_company_status (company_id, status),
    INDEX idx_fin_code (fin_code),
    INDEX idx_name (last_name, first_name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- --------------------------------------------------------
-- Təhsil (Blok B)
-- --------------------------------------------------------
CREATE TABLE IF NOT EXISTS employee_education (
    id INT AUTO_INCREMENT PRIMARY KEY,
    employee_id INT NOT NULL,
    institution VARCHAR(255) NOT NULL,
    specialty VARCHAR(255),
    degree ENUM('SECONDARY', 'VOCATIONAL', 'BACHELOR', 'MASTER', 'PHD') DEFAULT 'BACHELOR',
    start_year INT,
    end_year INT,
    diploma_number VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (employee_id) REFERENCES employees(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- --------------------------------------------------------
-- Əvvəlki İş Yerləri (Blok C)
-- --------------------------------------------------------
CREATE TABLE IF NOT EXISTS employee_experience (
    id INT AUTO_INCREMENT PRIMARY KEY,
    employee_id INT NOT NULL,
    company_name VARCHAR(255) NOT NULL,
    position VARCHAR(150),
    start_date DATE,
    end_date DATE,
    leaving_reason TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (employee_id) REFERENCES employees(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- --------------------------------------------------------
-- Ailə Məlumatları (Blok D)
-- --------------------------------------------------------
CREATE TABLE IF NOT EXISTS employee_family (
    id INT AUTO_INCREMENT PRIMARY KEY,
    employee_id INT NOT NULL,
    relation_type ENUM('FATHER', 'MOTHER', 'SPOUSE', 'CHILD', 'OTHER') NOT NULL,
    full_name VARCHAR(255) NOT NULL,
    birth_date DATE,
    contact_number VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (employee_id) REFERENCES employees(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- --------------------------------------------------------
-- İşçi Yaşam Dövrü Tarixçəsi
-- --------------------------------------------------------
CREATE TABLE IF NOT EXISTS employee_lifecycle_logs (
    id INT AUTO_INCREMENT PRIMARY KEY,
    employee_id INT NOT NULL,
    user_id INT NOT NULL,
    old_status VARCHAR(50),
    new_status VARCHAR(50),
    event_description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (employee_id) REFERENCES employees(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

SET FOREIGN_KEY_CHECKS = 1;
