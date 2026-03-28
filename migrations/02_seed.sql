-- Seed Data - Azman Holding System
-- Demo məlumatları (idempotent - INSERT IGNORE)

SET NAMES utf8mb4;

-- --------------------------------------------------------
-- Şirkətlər
-- --------------------------------------------------------
INSERT IGNORE INTO companies (name, is_holding, tax_id, address) VALUES
('Azman Holding', TRUE, '1701234567', 'Səbail r., Hacıbabayevlər küçəsi 4, Bakı, Azərbaycan'),
('Azman Construction', FALSE, '1702345678', 'Baku, Azerbaijan - 16 il təcrübə, ictimai bina və parkların tikintisi'),
('TEZ Logistics', FALSE, '1703456789', 'Baku, Azerbaijan'),
('Sapphire Hotels Group', FALSE, '1704567890', 'Baku, Azerbaijan'),
('City Service Company', FALSE, '1705678901', 'Baku, Azerbaijan'),
('EcoProd Azerbaijan', FALSE, '1706789012', 'Baku, Azerbaijan - Kənd təsərrüfatı və emalat'),
('Mangal MMC', FALSE, '1707890123', 'Baku, Azerbaijan'),
('Judo Club 2012', FALSE, '1708901234', 'Bakı, Azərbaycan');

-- --------------------------------------------------------
-- Departamentlər
-- --------------------------------------------------------
INSERT IGNORE INTO departments (company_id, name) VALUES
-- Azman Construction (2)
(2, 'İnzibati'), (2, 'Tikinti'), (2, 'Layihələndirmə'), (2, 'Təhlükəsizlik'), (2, 'Satınalma'), (2, 'Mühasibat'), (2, 'HR'),
-- TEZ Logistics (3)
(3, 'İnzibati'), (3, 'Əməliyyatlar'), (3, 'Daşıma'), (3, 'Anbar'), (3, 'Xarici Əlaqələr'), (3, 'HR'), (3, 'Mühasibat'),
-- Sapphire Hotels Group (4)
(4, 'İnzibati'), (4, 'Resepsion'), (4, 'Housekeeping'), (4, 'F&B'), (4, 'Münasibətlər'), (4, 'Satış'), (4, 'HR'),
-- City Service Company (5)
(5, 'İnzibati'), (5, 'Mətbəx'), (5, 'Xidmət'), (5, 'Tədarük'), (5, 'Marketinq'), (5, 'HR'), (5, 'Mühasibat'), (5, 'Təhlükəsizlik'),
-- EcoProd Azerbaijan (6)
(6, 'İnzibati'), (6, 'İstehsalat'), (6, 'Kənd Təsərrüfatı'), (6, 'Keyfiyyət'), (6, 'Satış'), (6, 'HR'),
-- Mangal MMC (7)
(7, 'İnzibati'), (7, 'İstehsalat'), (7, 'Texniki'), (7, 'Satış'), (7, 'HR'), (7, 'Logistika'),
-- Judo Club 2012 (8)
(8, 'İdman Sektoru'), (8, 'İnzibati'), (8, 'Təşkilatlar');

-- --------------------------------------------------------
-- Vəzifələr
-- --------------------------------------------------------
INSERT IGNORE INTO positions (company_id, name) VALUES
-- Azman Construction (2)
(2, 'Layihə Meneceri'), (2, 'Tikinti Mühəndisi'), (2, 'Memar'), (2, 'Kənd Təsərrüfatı İşçisi'), (2, 'Təhlükəsizlik Mütəxəssisi'), (2, 'Mühasib'), (2, 'HR Meneceri'),
-- TEZ Logistics (3)
(3, 'Logistika Meneceri'), (3, 'Sürücü'), (3, 'Anbarçı'), (3, 'Ekspeditor'), (3, 'Mühasib'), (3, 'HR Meneceri'), (3, 'Əməliyyat Mütəxəssisi'),
-- Sapphire Hotels Group (4)
(4, 'Otel Meneceri'), (4, 'Resepsionist'), (4, 'Housekeeper'), (4, 'Şef-Pəhriz'), (4, 'Ofisiant'), (4, 'Satış Meneceri'), (4, 'HR Meneceri'),
-- City Service Company (5)
(5, 'Restoran Meneceri'), (5, 'Şef'), (5, 'Sous-şef'), (5, 'Ofisiant'), (5, 'Kassir'), (5, 'Tədarük Mütəxəssisi'), (5, 'HR Meneceri'),
-- EcoProd Azerbaijan (6)
(6, 'İstehsalat Meneceri'), (6, 'Aqronom'), (6, 'Fermer'), (6, 'Keyfiyyət Mütəxəssisi'), (6, 'Satış Meneceri'), (6, 'HR Meneceri'),
-- Mangal MMC (7)
(7, 'İstehsalat Rəhbəri'), (7, 'Operator'), (7, 'Texnik'), (7, 'Satış Meneceri'), (7, 'HR Meneceri'), (7, 'Sürücü'),
-- Judo Club 2012 (8)
(8, 'Baş Trener'), (8, 'Trenyor'), (8, 'İdmançı'), (8, 'Klub Meneceri'), (8, 'Tədbir Təşkilatçısı');

-- --------------------------------------------------------
-- Demo işçilər (14 nəfər)
-- --------------------------------------------------------
INSERT IGNORE INTO employees (company_id, first_name, last_name, father_name, fin_code, birth_date, gender, phone, email, address, status, department_id, position_id, hire_date) VALUES
-- Azman Construction - Cari işçilər (dept IDs 1-7, pos IDs 8-14)
(2, 'Elçin', 'Məmmədov', 'Samir', 'A1B2C3D', '1985-03-15', 'MALE', '+994501234567', 'elcin.mammedov@azmanconstruction.az', 'Baku, Nəsimi r.', 'ACTIVE', 1, 8, '2020-01-15'),
(2, 'Aygün', 'Həsənova', 'Rəşid', 'B2C3D4E', '1990-07-22', 'FEMALE', '+994502345678', 'aygun.hasanova@azmanconstruction.az', 'Baku, Yasamal r.', 'ACTIVE', 6, 13, '2021-03-01'),
(2, 'Rəşad', 'Əliyev', 'Vaqif', 'C3D4E5F', '1988-11-08', 'MALE', '+994503456789', 'resad.eliyev@azmanconstruction.az', 'Baku, Xətai r.', 'ACTIVE', 2, 9, '2019-06-10'),
-- TEZ Logistics - Cari işçilər (dept IDs 8-14, pos IDs 15-21)
(3, 'Fərhad', 'Rzayev', 'Əkbər', 'D4E5F6G', '1987-04-25', 'MALE', '+994507890123', 'ferhad.rzayev@tezlogistics.az', 'Baku, Nərimanov r.', 'ACTIVE', 9, 15, '2020-09-01'),
(3, 'Leyla', 'Qasımlı', 'Nadir', 'E5F6G7H', '1993-12-03', 'FEMALE', '+994508901234', 'leyla.qasimli@tezlogistics.az', 'Baku, Binəqədi r.', 'ACTIVE', 13, 20, '2022-01-15'),
-- Sapphire Hotels Group - Cari işçilər (dept IDs 15-21, pos IDs 22-28)
(4, 'Nigar', 'Abdullayeva', 'Kamal', 'F6G7H8I', '1992-05-18', 'FEMALE', '+994504567890', 'nigar.abdullayeva@sapphirehotels.az', 'Baku, Səbail r.', 'ACTIVE', 16, 23, '2021-07-01'),
(4, 'Orxan', 'Səmədov', 'Elman', 'G7H8I9J', '1990-01-30', 'MALE', '+994505678901', 'orxan.semedov@sapphirehotels.az', 'Baku, Nizami r.', 'ACTIVE', 17, 24, '2022-03-15'),
-- City Service Company - Cari işçilər (dept IDs 22-29, pos IDs 29-35)
(5, 'Kamran', 'İsmayılov', 'Tofiq', 'H8I9J0K', '1985-09-12', 'MALE', '+994506789012', 'kamran.ismayilov@cityservice.az', 'Baku, Suraxanı r.', 'ACTIVE', 26, 29, '2020-04-01'),
(5, 'Günel', 'Məlikova', 'Ədalət', 'I9J0K1L', '1994-06-20', 'FEMALE', '+994509012345', 'gunel.melikova@cityservice.az', 'Baku, 28 May r.', 'ACTIVE', 25, 31, '2021-11-01'),
(5, 'Səbuhi', 'Novruzov', 'Mübariz', 'J0K1L2M', '1991-02-14', 'MALE', '+994501112233', 'sebuhi.novruzov@cityservice.az', 'Baku, Nəsimi r.', 'ACTIVE', 27, 33, '2022-06-01'),
-- EcoProd Azerbaijan - Cari işçilər (dept IDs 30-35, pos IDs 36-41)
(6, 'Elvin', 'Kazımzadə', 'Rafiq', 'K1L2M3N', '1988-08-05', 'MALE', '+994502223344', 'elvin.kazimzade@ecoprod.az', 'Gəncə, Azerbaijan', 'ACTIVE', 31, 36, '2021-01-10'),
-- Mangal MMC - Cari işçilər (dept IDs 36-41, pos IDs 42-47)
(7, 'Tural', 'Babayev', 'Hüseyn', 'L2M3N4P', '1993-03-28', 'MALE', '+994503334455', 'tural.babayev@mangalmmc.az', 'Baku, Balaxanı', 'ACTIVE', 37, 42, '2022-02-15'),
-- Namizədlər
(2, 'Nərgiz', 'Quliyeva', 'Fərhad', 'M3N4P5Q', '1995-02-14', 'FEMALE', '+994504445566', 'nergiz.quliyeva@mail.ru', 'Baku, Sabunçu r.', 'CANDIDATE', NULL, NULL, NULL),
(5, 'Ramin', 'Hüseynov', 'Zaur', 'N4P5Q6R', '1992-10-30', 'MALE', '+994505556677', 'ramin.huseynov@mail.ru', 'Baku, Xəzər r.', 'CANDIDATE', NULL, NULL, NULL),
-- İşdən çıxanlar
(3, 'Səid', 'Aslanov', 'Məmməd', 'P5Q6R7S', '1982-05-18', 'MALE', '+994506667788', 'said.aslanov@mail.ru', 'Baku, Gəncə r.', 'TERMINATED', 8, 15, '2019-01-15');

-- --------------------------------------------------------
-- Demo təhsil məlumatları
-- --------------------------------------------------------
INSERT IGNORE INTO employee_education (employee_id, institution, specialty, degree, start_year, end_year, diploma_number) VALUES
(1, 'BDU', 'İqtisadiyyat', 'BACHELOR', 2003, 2007, 'BDU-2007-001'),
(1, 'ADİU', 'MBA', 'MASTER', 2008, 2010, 'ADİU-2010-045'),
(2, 'UNEC', 'Mühasibat uçotu', 'BACHELOR', 2007, 2011, 'UNEC-2011-234'),
(3, 'ADNA', 'Mühəndislik', 'BACHELOR', 2006, 2010, 'ADNA-2010-567'),
(4, 'BMU', 'Logistika və Nəqliyyat', 'BACHELOR', 2005, 2009, 'BMU-2009-089'),
(7, 'ADTM', 'Otel İşçiliyi', 'BACHELOR', 2009, 2013, 'ADTM-2013-112'),
(9, 'BDU', 'Biznes İdarəetmə', 'MASTER', 2014, 2016, 'BDU-2016-334');

-- --------------------------------------------------------
-- Demo iş təcrübəsi
-- --------------------------------------------------------
INSERT IGNORE INTO employee_experience (employee_id, company_name, position, start_date, end_date, leaving_reason) VALUES
(1, 'Xəzər İstehsalat', 'Mühasib', '2007-06-01', '2015-03-01', 'Karyera inkişafı'),
(1, 'Azər Kimya', 'Maliyyə Meneceri', '2015-04-01', '2020-01-01', 'Daha yüksək maaş'),
(2, 'Bakı Ticarət', 'Mühasib köməkçisi', '2011-06-01', '2015-09-01', 'Ev təhsili'),
(4, 'Alpha Logistika', 'Ekspeditor', '2009-07-01', '2014-12-01', 'Şirkət bağlanması'),
(9, 'Marriott Baku', 'Ofisiant', '2010-08-01', '2016-05-01', 'Xaricə getmə');

-- --------------------------------------------------------
-- Demo ailə məlumatları
-- --------------------------------------------------------
INSERT IGNORE INTO employee_family (employee_id, relation_type, full_name, birth_date, contact_number) VALUES
(1, 'FATHER', 'Məmmədov Samir Ələkbər oğlu', '1955-05-10', '+994501112233'),
(1, 'MOTHER', 'Məmmədova Zərinə Məmməd qızı', '1960-08-22', '+994502223344'),
(1, 'SPOUSE', 'Məmmədova Leyla Rəşid qızı', '1988-03-18', '+994503334455'),
(1, 'CHILD', 'Məmmədov Kənan Elçin oğlu', '2015-07-25', NULL),
(2, 'SPOUSE', 'Həsənov Ramin Rəşid oğlu', '1985-11-12', '+994504445566'),
(9, 'SPOUSE', 'Novruzova Gülay Rəşid qızı', '1993-04-08', '+994505556677');

-- Update terminated employee
UPDATE employees SET termination_date = '2023-06-15', termination_reason = 'Öz istəyi ilə' WHERE fin_code = 'P5Q6R7S';
