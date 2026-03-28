-- Seed Data - Holding HR System
-- Demo məlumatları

SET NAMES utf8mb4;

-- --------------------------------------------------------
-- Şirkətlər
-- --------------------------------------------------------
INSERT INTO companies (name, is_holding, tax_id, address) VALUES
('ABC Holding', TRUE, '1234567891', 'Bakı şəh., Nizami r., 123'),
('ABC Tekstil MMC', FALSE, '1234567892', 'Bakı şəh., Suraxanı r., 45'),
('ABC Logistika MMC', FALSE, '1234567893', 'Bakı şəh., Xətai r., 78'),
('ABC İstehsalat MMC', FALSE, '1234567894', 'Sumqayıt şəh., Səməd Vurğun küç., 12');

-- --------------------------------------------------------
-- Departamentlər
-- --------------------------------------------------------
INSERT INTO departments (company_id, name) VALUES
-- ABC Tekstil (company_id = 2)
(2, 'İnsan Resursları'),
(2, 'Maliyyə'),
(2, 'İstehsalat'),
(2, 'Satış'),
(2, 'Logistika'),
-- ABC Logistika (company_id = 3)
(3, 'İnsan Resursları'),
(3, 'Nəqliyyat'),
(3, 'Anbar'),
(3, 'Maliyyə'),
-- ABC İstehsalat (company_id = 4)
(4, 'İnsan Resursları'),
(4, 'İstehsalat'),
(4, 'Texniki Xidmət'),
(4, 'Keyfiyyətə Nəzarət');

-- --------------------------------------------------------
-- Vəzifələr
-- --------------------------------------------------------
INSERT INTO positions (company_id, name) VALUES
-- ABC Tekstil
(2, 'HR Menecer'),
(2, 'HR Mütəxəssis'),
(2, 'Maliyyə Meneceri'),
(2, 'Mühasib'),
(2, 'İstehsalat Rəhbəri'),
(2, 'Tikmə ustası'),
(2, 'Satış Meneceri'),
(2, 'Satış nümayəndəsi'),
-- ABC Logistika
(3, 'HR Menecer'),
(3, 'Nəqliyyat Rəhbəri'),
(3, 'Sürücü'),
(3, 'Anbar Rəhbəri'),
(3, 'Anbar işçisi'),
-- ABC İstehsalat
(4, 'HR Menecer'),
(4, 'İstehsalat Rəhbəri'),
(4, 'Mühəndis'),
(4, 'Texnik'),
(4, 'Keyfiyyət Nəzarətçisi');

-- --------------------------------------------------------
-- Admin istifadəçiləri application start olarkən yaradılır
-- Şifrə: admin123
-- --------------------------------------------------------

-- --------------------------------------------------------
-- Demo işçilər
-- --------------------------------------------------------
INSERT INTO employees (company_id, first_name, last_name, father_name, fin_code, birth_date, gender, phone, email, address, status, department_id, position_id, hire_date) VALUES
-- ABC Tekstil - Cari işçilər
(2, 'Əli', 'Məmmədov', 'İbrahim', '1A2B3CD', '1985-03-15', 'MALE', '+994501234567', 'eli.memmedov@abctekstil.az', 'Bakı şəh., Nəsimi r.', 'ACTIVE', 1, 5),
(2, 'Aygün', 'Həsənova', 'Rəşid', '2B3C4DE', '1990-07-22', 'FEMALE', '+994502345678', 'aygun.hasanova@abctekstil.az', 'Bakı şəh., Nizami r.', 'ACTIVE', 2, 4),
(2, 'Rəşad', 'Əliyev', 'Vaqif', '3C4D5EF', '1988-11-08', 'MALE', '+994503456789', 'resad.eliyev@abctekstil.az', 'Bakı şəh., Xətai r.', 'ACTIVE', 3, 6),
-- ABC Tekstil - Namizədlər
(2, 'Nərgiz', 'Quliyeva', 'Fərhad', '4D5E6FG', '1995-02-14', 'FEMALE', '+994504567890', 'nergiz.quliyeva@mail.ru', 'Bakı şəh., Suraxanı r.', 'CANDIDATE', NULL, NULL, NULL),
(2, 'Kamran', 'İsmayılov', 'Tofiq', '5E6F7GH', '1992-09-30', 'MALE', '+994505678901', 'kamran.ismayilov@mail.ru', 'Bakı şəh., Sabunçu r.', 'CANDIDATE', NULL, NULL, NULL),
-- ABC Tekstil - İşdən çıxanlar
(2, 'Səid', 'Hüseynov', 'Məmməd', '6F7G8HI', '1982-05-18', 'MALE', '+994506789012', 'said.huseynov@mail.ru', 'Bakı şəh., Yasamal r.', 'TERMINATED', 1, 5, '2020-01-15'),
-- ABC Logistika - Cari işçilər
(3, 'Fərhad', 'Rzayev', 'Əkbər', '7G8H9IJ', '1987-04-25', 'MALE', '+994507890123', 'ferhad.rzayev@abclogistika.az', 'Bakı şəh., Nərimanov r.', 'ACTIVE', 7, 10),
(3, 'Leyla', 'Qasımlı', 'Nadir', '8H9I0JK', '1993-12-03', 'FEMALE', '+994508901234', 'leyla.qasimli@abclogistika.az', 'Bakı şəh., Binəqədi r.', 'ACTIVE', 6, 9);

-- --------------------------------------------------------
-- Demo təhsil məlumatları
-- --------------------------------------------------------
INSERT INTO employee_education (employee_id, institution, specialty, degree, start_year, end_year, diploma_number) VALUES
(1, 'BDU', 'İqtisadiyyat', 'BACHELOR', 2003, 2007, 'BDU-2007-001'),
(1, 'ADİU', 'MBA', 'MASTER', 2008, 2010, 'ADİU-2010-045'),
(2, 'UNEC', 'Mühasibat uçotu', 'BACHELOR', 2007, 2011, 'UNEC-2011-234'),
(3, 'ADNA', 'Mühəndislik', 'BACHELOR', 2006, 2010, 'ADNA-2010-567');

-- --------------------------------------------------------
-- Demo iş təcrübəsi
-- --------------------------------------------------------
INSERT INTO employee_experience (employee_id, company_name, position, start_date, end_date, leaving_reason) VALUES
(1, 'Xəzər İstehsalat', 'Mühasib', '2007-06-01', '2015-03-01', 'Karyera inkişafı'),
(1, 'Azər Kimya', 'Maliyyə Meneceri', '2015-04-01', '2020-01-01', 'Daha yüksək maaş'),
(2, 'Bakı Ticarət', 'Mühasib köməkçisi', '2011-06-01', '2015-09-01', 'Ev təhsili');

-- --------------------------------------------------------
-- Demo ailə məlumatları
-- --------------------------------------------------------
INSERT INTO employee_family (employee_id, relation_type, full_name, birth_date, contact_number) VALUES
(1, 'FATHER', 'Məmmədov İbrahim Ələkbər oğlu', '1955-05-10', '+994501112233'),
(1, 'MOTHER', 'Məmmədov Zərinə Məmməd qızı', '1960-08-22', '+994502223344'),
(1, 'SPOUSE', 'Məmmədov Leyla Rəşid qızı', '1988-03-18', '+994503334455'),
(1, 'CHILD', 'Məmmədov Kənan Əli oğlu', '2015-07-25', NULL),
(2, 'SPOUSE', 'Həsənov Ramin Rəşid oğlu', '1985-11-12', '+994504445566');

-- Update terminated employee
UPDATE employees SET termination_date = '2023-06-15', termination_reason = 'Öz istəyi ilə' WHERE id = 6;
