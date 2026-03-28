-- Migration 05: Company Hierarchy
-- parent_id və company_type əlavə edir (DELIMITER istifadə etmir)
-- Qeyd: Əsas column əlavə cmd/main.go runMigrations() safeAddColumn ilə edilir
-- Bu fayl yalnız additional alter-lar edir

SET NAMES utf8mb4;

-- Mövcud şirkətlərə tip təyin et
UPDATE IGNORE companies SET company_type = 'HOLDING' WHERE name = 'Azman Holding';
UPDATE IGNORE companies SET company_type = 'CONSTRUCTION' WHERE name = 'Azman Construction';
UPDATE IGNORE companies SET company_type = 'LOGISTICS' WHERE name = 'TEZ Logistics';
UPDATE IGNORE companies SET company_type = 'HOTEL' WHERE name = 'Sapphire Hotels Group';
UPDATE IGNORE companies SET company_type = 'RESTAURANT_CHAIN' WHERE name = 'City Service Company';
UPDATE IGNORE companies SET company_type = 'AGRICULTURE' WHERE name = 'EcoProd Azerbaijan';
UPDATE IGNORE companies SET company_type = 'PRODUCTION' WHERE name = 'Mangal MMC';
UPDATE IGNORE companies SET company_type = 'SPORTS' WHERE name = 'Judo Club 2012';
