-- Migration 06: City Service Restaurant Sub-Companies
-- City Service Company (ID 5) altında 7 restoran/kafe/banket zalı

SET NAMES utf8mb4;

-- --------------------------------------------------------
-- Alt-restoranlar (parent_id = 5 → City Service Company)
-- --------------------------------------------------------
INSERT IGNORE INTO companies (name, parent_id, is_holding, company_type, tax_id, address) VALUES
('Nar&Sharab', 5, FALSE, 'RESTAURANT', '1705678902', 'Bakı, Azərbaycan - Restoran & Şərab Zalı'),
('Senator Hall', 5, FALSE, 'RESTAURANT', '1705678903', 'Bakı, Azərbaycan - Banket Zalı'),
('Maxim Hall', 5, FALSE, 'RESTAURANT', '1705678904', 'Bakı, Azərbaycan - Banket Zalı'),
('Mangal Steak House Zagulba', 5, FALSE, 'RESTAURANT', '1705678905', 'Bakı, Zagulba - Restoran'),
('Zuğulba Saray', 5, FALSE, 'RESTAURANT', '1705678906', 'Bakı, Zagulba - Banket & Əyləncə Məkanı'),
('Beer Hall', 5, FALSE, 'RESTAURANT', '1705678907', 'Bakı, Azərbaycan - Restoran & Bar'),
('JazzClub by BeerHall', 5, FALSE, 'RESTAURANT', '1705678908', 'Bakı, Azərbaycan - Jazz Klub & Bar');

-- --------------------------------------------------------
-- Restoran departamentləri (hər restoran üçün universal)
-- --------------------------------------------------------
-- Nar&Sharab departamentləri
INSERT IGNORE INTO departments (company_id, name) VALUES
-- Nar&Sharab
(9, 'İnzibati'), (9, 'Mətbəx'), (9, 'Xidmət'), (9, 'Bar'), (9, 'Tədarük'), (9, 'Marketinq');

-- Senator Hall
(10, 'İnzibati'), (10, 'Mətbəx'), (10, 'Xidmət'), (10, 'Tədarük'), (10, 'Banket Xidməti');

-- Maxim Hall
(11, 'İnzibati'), (11, 'Mətbəx'), (11, 'Xidmət'), (11, 'Tədarük'), (11, 'Banket Xidməti');

-- Mangal Steak House Zagulba
(12, 'İnzibati'), (12, 'Mətbəx'), (12, 'Xidmət'), (12, 'Barbekyu Sahəsi'), (12, 'Tədarük');

-- Zuğulba Saray
(13, 'İnzibati'), (13, 'Mətbəx'), (13, 'Xidmət'), (13, 'Təşkilat'), (13, 'Tədarük'), (13, 'Əyləncə');

-- Beer Hall
(14, 'İnzibati'), (14, 'Mətbəx'), (14, 'Xidmət'), (14, 'Brewery'), (14, 'Bar');

-- JazzClub by BeerHall
(15, 'İnzibati'), (15, 'Xidmət'), (15, 'Bar'), (15, 'Musiqi & Əyləncə');

-- --------------------------------------------------------
-- Restoran vəzifələri
-- --------------------------------------------------------
INSERT IGNORE INTO positions (company_id, name) VALUES
-- Nar&Sharab (9)
(9, 'Restoran Meneceri'), (9, 'Baş Şef'), (9, 'Sous-şef'), (9, 'Ofisiant'), (9, 'Barmen'), (9, 'Somelye'), (9, 'Kassir'), (9, 'Təmizlik İşçisi'),
-- Senator Hall (10)
(10, 'Zal Meneceri'), (10, 'Baş Şef'), (10, 'Şef-pəhriz'), (10, 'Ofisiant'), (10, 'Banket Koordinatoru'), (10, 'Kassir'), (10, 'Təmizlik İşçisi'),
-- Maxim Hall (11)
(11, 'Zal Meneceri'), (11, 'Baş Şef'), (11, 'Şef-pəhriz'), (11, 'Ofisiant'), (11, 'Banket Koordinatoru'), (11, 'Kassir'), (11, 'Təmizlik İşçisi'),
-- Mangal Steak House Zagulba (12)
(12, 'Restoran Meneceri'), (12, 'Baş Şef'), (12, 'Barbekyu Usta'), (12, 'Ofisiant'), (12, 'Kassir'), (12, 'Təmizlik İşçisi'),
-- Zuğulba Saray (13)
(13, 'Saray Meneceri'), (13, 'Baş Şef'), (13, 'Şef-pəhriz'), (13, 'Ofisiant'), (13, 'Tədbir Koordinatoru'), (13, 'Dekorator'), (13, 'Kassir'), (13, 'Təhlükəsizlik'),
-- Beer Hall (14)
(14, 'Restoran Meneceri'), (14, 'Baş Şef'), (14, 'Brew Master'), (14, 'Ofisiant'), (14, 'Barmen'), (14, 'Kassir'), (14, 'Təmizlik İşçisi'),
-- JazzClub by BeerHall (15)
(15, 'Klub Meneceri'), (15, 'Barmen'), (15, 'Ofisiant'), (15, 'Səs Mühəndisi'), (15, 'Təhlükəsizlik'), (15, 'Təmizlik İşçisi');

-- --------------------------------------------------------
-- Demo işçilər - Restoranlar üçün (bəzi restoranlara nümunə işçilər)
-- --------------------------------------------------------
INSERT IGNORE INTO employees (company_id, first_name, last_name, father_name, fin_code, birth_date, gender, phone, email, address, status, department_id, position_id, hire_date) VALUES
-- Nar&Sharab
(9, 'Emil', 'Əlizadə', 'Rafiq', 'Q1R2S3T', '1990-01-20', 'MALE', '+994507778899', 'emil.alizade@narsharab.az', 'Bakı', 'ACTIVE', 36, 55, '2018-05-01'),
(9, 'Ləman', 'Hüseynova', 'Vüsal', 'R2S3T4U', '1992-06-15', 'FEMALE', '+994508889900', 'leman.huseynova@narsharab.az', 'Bakı', 'ACTIVE', 37, 58, '2019-03-15'),
-- Beer Hall
(14, 'Tural', 'Məlikov', 'Elşən', 'S3T4U5V', '1988-09-10', 'MALE', '+994509990011', 'tural.melikov@beerhall.az', 'Bakı', 'ACTIVE', 76, 79, '2017-08-20'),
(14, 'Ayan', 'Kərimova', 'Namiq', 'T4U5V6W', '1994-03-25', 'FEMALE', '+994501112233', 'ayan.kerimova@beerhall.az', 'Bakı', 'ACTIVE', 78, 80, '2020-11-01'),
-- JazzClub by BeerHall
(15, 'Rüfət', 'Quliyev', 'Kamran', 'U5V6W7X', '1991-11-30', 'MALE', '+994502233445', 'rufet.guliyev@jazzclub.az', 'Bakı', 'ACTIVE', 84, 85, '2021-06-01'),
-- Zuğulba Saray
(13, 'Günay', 'Əhmədova', 'Sərvər', 'V6W7X8Y', '1989-07-08', 'FEMALE', '+994503344556', 'gunay.ahmedova@zugulbasaray.az', 'Bakı', 'ACTIVE', 67, 68, '2019-09-10'),
-- Mangal Steak House
(12, 'Elbrus', 'Nəsirov', 'Akif', 'W7X8Y9Z', '1986-04-12', 'MALE', '+994504455667', 'elbrus.nasirov@mangalsteak.az', 'Bakı, Zagulba', 'ACTIVE', 63, 66, '2020-01-15'),
-- Senator Hall
(10, 'Nurlana', 'Aslanova', 'Fikrət', 'X8Y9Z0A', '1993-12-18', 'FEMALE', '+994505566778', 'nurlana.aslanova@senatorhall.az', 'Bakı', 'ACTIVE', 42, 47, '2021-04-01'),
-- Maxim Hall
(11, 'Samir', 'Bədəlov', 'Məmməd', 'Y9Z0A1B', '1987-08-22', 'MALE', '+994506677889', 'samir.bedelov@maximhall.az', 'Bakı', 'ACTIVE', 48, 53, '2018-11-15');
