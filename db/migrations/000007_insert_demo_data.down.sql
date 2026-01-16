-- 000007_insert_demo_data.down.sql
-- Remove demo data from veterinary clinic management system
-- WARNING: This will delete ALL data from the tables. Use with caution.

-- Delete in reverse order to respect foreign key constraints

-- Delete payments first (references medical_sessions and customers)
DELETE FROM payments WHERE invoice_id LIKE 'INV-2024-%';

-- Delete pet chip implants
DELETE FROM pet_chip_implants WHERE chip_number LIKE '982000123456%';

-- Delete medical sessions (references pets, customers, employees, appointments)
DELETE FROM medical_sessions WHERE visit_date >= '2024-01-01';

-- Delete appointments (references pets, customers, employees)
DELETE FROM appointments WHERE scheduled_date >= '2024-01-01';

-- Delete pet deworming (references pets and employees)
DELETE FROM pet_deworming WHERE administered_date >= '2024-01-01';

-- Delete pet vaccinations (references pets and employees)
DELETE FROM pet_vaccinations WHERE administered_date >= '2024-01-01';

-- Delete pet behavioral notes (references pets)
DELETE FROM pet_behavioral_notes WHERE noted_at >= '2024-01-01';

-- Delete pet feeding instructions (references pets)
DELETE FROM pet_feeding_instructions;

-- Delete pets (references customers)
DELETE FROM pets WHERE name IN ('Max', 'Luna', 'Bobby', 'Mila', 'Rex', 'Coco', 'Simba');

-- Delete employees (references users)
DELETE FROM employees WHERE license_number IN ('VET-2020-001', 'VET-2018-005', 'VET-2022-012');

-- Delete customers (references users)
DELETE FROM customers WHERE first_name IN ('Juan Carlos', 'María Elena', 'Carlos Alberto', 'Ana Patricia', 'Pedro Antonio');

-- Delete demo users
DELETE FROM users WHERE email LIKE '%@vetclinic.com' OR email LIKE '%@email.com';
