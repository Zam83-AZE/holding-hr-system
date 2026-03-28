-- Migration 04: Update from website data
-- Judo Club 2012, new work locations, address updates

SET NAMES utf8mb4;

-- --------------------------------------------------------
-- 1. Judo Club 2012 - New company (company_id will be 8)
-- --------------------------------------------------------
INSERT INTO companies (name, is_holding, tax_id, address) VALUES
('Judo Club 2012', FALSE, '1708901234', 'Bakı, Azərbaycan');

-- --------------------------------------------------------
-- 2. Judo Club departments
-- --------------------------------------------------------
INSERT INTO departments (company_id, name) VALUES
(8, 'İdman Sektoru'),
(8, 'İnzibati'),
(8, 'Təşkilatlar');

-- --------------------------------------------------------
-- 3. Judo Club positions
-- --------------------------------------------------------
INSERT INTO positions (company_id, name) VALUES
(8, 'Baş Trener'),
(8, 'Trenyor'),
(8, 'İdmançı'),
(8, 'Klub Meneceri'),
(8, 'Tədbir Təşkilatçısı');

-- --------------------------------------------------------
-- 4. Judo Club work location
-- --------------------------------------------------------
INSERT INTO work_locations (company_id, name, address, type) VALUES
(8, 'Judo Club 2012 Arenası', 'Bakı, Azərbaycan', 'OTHER');

-- --------------------------------------------------------
-- 5. Missing City Service restaurants
-- --------------------------------------------------------
INSERT INTO work_locations (company_id, name, address, type) VALUES
(5, 'Beer Hall', 'Bakı, Azərbaycan', 'RESTAURANT'),
(5, 'JazzClub by BeerHall', 'Bakı, Azərbaycan', 'RESTAURANT');

-- --------------------------------------------------------
-- 6. Mangal MMC restaurants
-- --------------------------------------------------------
INSERT INTO work_locations (company_id, name, address, type) VALUES
(7, 'Steak House', 'Bakı, Azərbaycan', 'RESTAURANT'),
(7, 'Mangal Milli', 'Bakı, Azərbaycan', 'RESTAURANT'),
(7, 'Ocakbaşı Mangal', 'Bakı, Azərbaycan', 'RESTAURANT'),
(7, 'Ballroom', 'Bakı, Azərbaycan', 'RESTAURANT');

-- --------------------------------------------------------
-- 7. Sapphire Hotels - add missing hotels
-- --------------------------------------------------------
INSERT INTO work_locations (company_id, name, address, type) VALUES
(4, 'Sapphire INN', 'Bakı, Azərbaycan', 'HOTEL'),
(4, 'Sapphire City', 'Bakı, Azərbaycan', 'HOTEL'),
(4, 'Sapphire Baku', 'Bakı, Azərbaycan', 'HOTEL'),
(4, 'Sapphire Bayıl', 'Bakı, Bayıl', 'HOTEL'),
(4, 'Sapphire Marine', 'Bakı, Azərbaycan', 'HOTEL');

-- --------------------------------------------------------
-- 8. Update existing Sapphire work locations
-- --------------------------------------------------------
UPDATE work_locations SET type = 'OFFICE', name = 'Sapphire Ofis' WHERE company_id = 4 AND name = 'Sapphire Business Center';

-- --------------------------------------------------------
-- 9. Update EcoProd address
-- --------------------------------------------------------
UPDATE work_locations SET address = 'Samux rayonu, Azərbaycan' WHERE company_id = 6 AND type = 'FARM';
