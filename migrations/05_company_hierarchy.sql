-- Migration 05: Company Hierarchy
-- parent_id əlavə edirik ki, City Service altında restoranlar yaradayaq

SET NAMES utf8mb4;

-- parent_id sütunu əlavə et (stored procedure ilə təhlükəsiz)
DELIMITER //
CREATE PROCEDURE add_parent_id_column() BEGIN
    IF NOT EXISTS (
        SELECT * FROM information_schema.COLUMNS
        WHERE TABLE_SCHEMA = DATABASE()
        AND TABLE_NAME = 'companies'
        AND COLUMN_NAME = 'parent_id'
    ) THEN
        ALTER TABLE companies ADD COLUMN parent_id INT NULL;
    END IF;
END //
DELIMITER ;
CALL add_parent_id_column();
DROP PROCEDURE IF EXISTS add_parent_id_column;

-- company_type sütunu əlavə et (icon üçün)
DELIMITER //
CREATE PROCEDURE add_company_type_column() BEGIN
    IF NOT EXISTS (
        SELECT * FROM information_schema.COLUMNS
        WHERE TABLE_SCHEMA = DATABASE()
        AND TABLE_NAME = 'companies'
        AND COLUMN_NAME = 'company_type'
    ) THEN
        ALTER TABLE companies ADD COLUMN company_type VARCHAR(50) DEFAULT 'OTHER';
    END IF;
END //
DELIMITER ;
CALL add_company_type_column();
DROP PROCEDURE IF EXISTS add_company_type_column;

-- Foreign key
-- MariaDB 10.11: ADD CONSTRAINT IF NOT EXISTS dəstəkləmir, ignore edirik
SET @fk_exists = (SELECT COUNT(*) FROM information_schema.TABLE_CONSTRAINTS
    WHERE CONSTRAINT_SCHEMA = DATABASE()
    AND TABLE_NAME = 'companies'
    AND CONSTRAINT_NAME = 'fk_companies_parent');
SET @sql = IF(@fk_exists = 0,
    'ALTER TABLE companies ADD CONSTRAINT fk_companies_parent FOREIGN KEY (parent_id) REFERENCES companies(id) ON DELETE SET NULL',
    'SELECT 1');
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Mövcud şirkətlərə tip təyin et
UPDATE companies SET company_type = 'HOLDING' WHERE name = 'Azman Holding';
UPDATE companies SET company_type = 'CONSTRUCTION' WHERE name = 'Azman Construction';
UPDATE companies SET company_type = 'LOGISTICS' WHERE name = 'TEZ Logistics';
UPDATE companies SET company_type = 'HOTEL' WHERE name = 'Sapphire Hotels Group';
UPDATE companies SET company_type = 'RESTAURANT_CHAIN' WHERE name = 'City Service Company';
UPDATE companies SET company_type = 'AGRICULTURE' WHERE name = 'EcoProd Azerbaijan';
UPDATE companies SET company_type = 'PRODUCTION' WHERE name = 'Mangal MMC';
UPDATE companies SET company_type = 'SPORTS' WHERE name = 'Judo Club 2012';

-- Index
SET @idx_exists = (SELECT COUNT(*) FROM information_schema.STATISTICS
    WHERE TABLE_SCHEMA = DATABASE()
    AND TABLE_NAME = 'companies'
    AND INDEX_NAME = 'idx_parent_id');
SET @sql = IF(@idx_exists = 0,
    'CREATE INDEX idx_parent_id ON companies(parent_id)',
    'SELECT 1');
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;
