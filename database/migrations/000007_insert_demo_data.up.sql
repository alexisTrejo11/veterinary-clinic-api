-- 000007_insert_demo_data.up.sql
-- Insert demo data for veterinary clinic management system
-- WARNING: This script is for demonstration purposes only. Do not use in production.

-- Insert demo users (password is 'password123' hashed with bcrypt)
INSERT INTO users (email, phone_number, hashed_password, status, role, last_login, name) VALUES
-- Veterinarians
('dr.martinez@vetclinic.com', '+52-555-0001', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'active', 'veterinarian', '2024-01-15 08:30:00', 'Dr. Martinez'),
('dr.garcia@vetclinic.com', '+52-555-0002', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'active', 'veterinarian', '2024-01-15 09:15:00', 'Dr. Garcia'),
('dr.rodriguez@vetclinic.com', '+52-555-0003', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'active', 'veterinarian', '2024-01-14 14:20:00', 'Dr. Rodriguez'),
-- Receptionists and Admin
('admin@vetclinic.com', '+52-555-0004', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'active', 'admin', '2024-01-15 07:45:00', 'Admin User'),
('recepcion@vetclinic.com', '+52-555-0005', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'active', 'receptionist', '2024-01-15 08:00:00', 'Recepcion'),
-- Customers
('juan.lopez@email.com', '+52-555-1001', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'active', 'customer', '2024-01-14 18:30:00', 'Juan Lopez'),
('maria.gonzalez@email.com', '+52-555-1002', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'active', 'customer', '2024-01-13 20:15:00', 'Maria Gonzalez'),
('carlos.hernandez@email.com', '+52-555-1003', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'active', 'customer', '2024-01-12 16:45:00', 'Carlos Hernandez'),
('ana.flores@email.com', '+52-555-1004', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'active', 'customer', '2024-01-15 12:20:00', 'Ana Flores'),
('pedro.silva@email.com', '+52-555-1005', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'active', 'customer', NULL, 'Pedro Silva');

-- Insert demo customers
INSERT INTO customers (first_name, last_name, photo, date_of_birth, gender, user_id, is_active) VALUES
('Juan Carlos', 'López Méndez', 'https://www.gravatar.com/avatar/juan_lopez', '1985-03-15', 'male', 6, TRUE),
('María Elena', 'González Ruiz', 'https://www.gravatar.com/avatar/maria_gonzalez', '1990-07-22', 'female', 7, TRUE),
('Carlos Alberto', 'Hernández Vega', 'https://www.gravatar.com/avatar/carlos_hernandez', '1982-11-08', 'male', 8, TRUE),
('Ana Patricia', 'Flores Morales', 'https://www.gravatar.com/avatar/ana_flores', '1988-05-12', 'female', 9, TRUE),
('Pedro Antonio', 'Silva Jiménez', 'https://www.gravatar.com/avatar/pedro_silva', '1975-09-30', 'male', 10, TRUE);

-- Insert demo employees (veterinarians)
INSERT INTO employees (first_name, last_name, gender, date_of_birth, photo, license_number, speciality, years_of_experience, is_active, user_id, schedule_json) VALUES
('Dr. Eduardo', 'Martínez Sánchez', 'male', '1980-04-20', 'https://www.gravatar.com/avatar/dr_martinez', 'VET-2020-001', 'general_practice', 8, TRUE, 1,
'{"monday": {"start": "08:00", "end": "17:00"}, "tuesday": {"start": "08:00", "end": "17:00"}, "wednesday": {"start": "08:00", "end": "17:00"}, "thursday": {"start": "08:00", "end": "17:00"}, "friday": {"start": "08:00", "end": "16:00"}, "saturday": {"start": "09:00", "end": "13:00"}, "sunday": "off"}'),

('Dra. Carmen', 'García López', 'female', '1978-09-14', 'https://www.gravatar.com/avatar/dra_garcia', 'VET-2018-005', 'surgery', 12, TRUE, 2,
'{"monday": {"start": "07:00", "end": "15:00"}, "tuesday": {"start": "07:00", "end": "15:00"}, "wednesday": {"start": "07:00", "end": "15:00"}, "thursday": {"start": "07:00", "end": "15:00"}, "friday": {"start": "07:00", "end": "15:00"}, "saturday": "off", "sunday": "on_call"}'),

('Dr. Roberto', 'Rodríguez Moreno', 'male', '1985-12-03', 'https://www.gravatar.com/avatar/dr_rodriguez', 'VET-2022-012', 'emergency_critical_care', 5, TRUE, 3,
'{"monday": {"start": "14:00", "end": "22:00"}, "tuesday": {"start": "14:00", "end": "22:00"}, "wednesday": {"start": "14:00", "end": "22:00"}, "thursday": {"start": "14:00", "end": "22:00"}, "friday": {"start": "14:00", "end": "22:00"}, "saturday": {"start": "10:00", "end": "18:00"}, "sunday": {"start": "10:00", "end": "18:00"}}');

-- Insert demo pets
INSERT INTO pets (name, photo, species, breed, age, gender, color, microchip, blood_type, is_neutered, customer_id, is_active, allergies, current_medications, special_needs, feeding_instructions, behavioral_notes, veterinary_contact, emergency_contact_name, emergency_contact_phone) VALUES
-- Pets for Juan Carlos López
('Max', 'https://example.com/photos/max.jpg', 'Dog', 'Golden Retriever', 5, 'male', 'Golden', '982000123456789', 'A+', TRUE, 1, TRUE,
'Ninguna conocida', NULL, NULL,
'2 tazas de alimento seco dos veces al día',
'Muy amigable, le gusta jugar con otros perros',
'Dr. Martínez - Clínica Veterinaria Central',
'Juan Carlos López', '+52-555-1001'),

('Luna', 'https://example.com/photos/luna.jpg', 'Cat', 'Siamese', 3, 'female', 'Seal Point', '982000123456790', 'B+', TRUE, 1, TRUE,
'Alergia al pollo', 'Antihistamínicos según necesidad', NULL,
'Alimento húmedo sin pollo, 3 veces al día',
'Tímida con extraños, prefiere lugares altos',
'Dr. Martínez - Clínica Veterinaria Central',
'Juan Carlos López', '+52-555-1001'),

-- Pets for María Elena González
('Bobby', 'https://example.com/photos/bobby.jpg', 'Dog', 'Bulldog Francés', 2, 'male', 'Brindle', '982000123456791', 'A+', FALSE, 2, TRUE,
'Sensibilidad alimentaria', 'Dieta hipoalergénica', 'Problemas respiratorios leves',
'Alimento hipoalergénico 3 veces al día, porciones pequeñas',
'Energético pero se cansa fácilmente debido a su respiración',
'Dra. García - Clínica Veterinaria Central',
'María Elena González', '+52-555-1002'),

-- Pets for Carlos Alberto Hernández
('Mila', 'https://example.com/photos/mila.jpg', 'Cat', 'Persian', 4, 'female', 'White', '982000123456792', 'AB+', TRUE, 3, TRUE,
'Ninguna conocida', NULL, 'Requiere cepillado diario',
'Alimento premium para gatos persas, 2 veces al día',
'Muy tranquila, disfruta ser cepillada',
'Dr. Rodríguez - Clínica Veterinaria Central',
'Carlos Alberto Hernández', '+52-555-1003'),

('Rex', 'https://example.com/photos/rex.jpg', 'Dog', 'Pastor Alemán', 7, 'male', 'Black and Tan', '982000123456793', 'A+', TRUE, 3, TRUE,
'Ninguna conocida', 'Suplemento para articulaciones', 'Displasia de cadera leve',
'3 tazas de alimento senior divididas en 2 comidas',
'Protector, muy leal, entrenado como perro guardián',
'Dr. Rodríguez - Clínica Veterinaria Central',
'Carlos Alberto Hernández', '+52-555-1003'),

-- Pets for Ana Patricia Flores
('Coco', 'https://example.com/photos/coco.jpg', 'Bird', 'Cockatiel', 2, 'unknown', 'Gray', NULL, 'A+', FALSE, 4, TRUE,
'Ninguna conocida', NULL, 'Ave sensible al humo y aerosoles',
'Semillas premium para cacatúas y frutas frescas diariamente',
'Muy social, imita sonidos, disfruta de la música',
'Dra. García - Clínica Veterinaria Central',
'Ana Patricia Flores', '+52-555-1004'),

-- Pets for Pedro Antonio Silva
('Simba', 'https://example.com/photos/simba.jpg', 'Cat', 'Maine Coon', 6, 'male', 'Brown Tabby', '982000123456794', 'A+', TRUE, 5, TRUE,
'Ninguna conocida', NULL, NULL,
'Alimento seco de alta calidad, 2 veces al día, gran cantidad por su tamaño',
'Muy grande y gentil, le gusta el agua',
'Dr. Martínez - Clínica Veterinaria Central',
'Pedro Antonio Silva', '+52-555-1005');

-- Insert demo pet feeding instructions
INSERT INTO pet_feeding_instructions (pet_id, food_brand, food_type, amount_per_serving, frequency, special_instructions) VALUES
(1, 'Royal Canin', 'Adult Golden Retriever', '2 tazas', '2 veces al día', 'Mantener horarios regulares de alimentación'),
(2, 'Hill''s Prescription Diet', 'Feline Z/D', '1/2 lata', '3 veces al día', 'Sin pollo ni subproductos de pollo'),
(3, 'Purina Pro Plan', 'Sensitive Skin & Stomach', '3/4 taza', '3 veces al día', 'Porciones pequeñas para evitar problemas digestivos'),
(4, 'Royal Canin', 'Persian Adult', '1/2 taza', '2 veces al día', 'Usar platos poco profundos para evitar problemas con su cara plana'),
(5, 'Hill''s Science Diet', 'Senior Large Breed', '1.5 tazas', '2 veces al día', 'Agregar suplemento para articulaciones'),
(6, 'ZuPreem', 'Cockatiel Mix', '2 cucharadas', 'Ad libitum', 'Suplementar con frutas y verduras frescas'),
(7, 'Blue Buffalo', 'Maine Coon Formula', '1 taza', '2 veces al día', 'Puede requerir más cantidad debido a su gran tamaño');

-- Insert demo pet behavioral notes
INSERT INTO pet_behavioral_notes (pet_id, note, noted_at, noted_by) VALUES
(1, 'Muy sociable con otros perros en el parque. Excelente temperamento para terapia.', '2024-01-10 10:30:00', 1),
(2, 'Prefiere observar desde lugares altos. Se esconde cuando hay visitas nuevas.', '2024-01-08 14:15:00', 1),
(3, 'Respiración laboriosa durante ejercicio intenso. Limitar actividad en clima caluroso.', '2024-01-12 09:45:00', 2),
(4, 'Muy cooperativa durante el aseo. Disfruta ser cepillada y examinada.', '2024-01-05 11:20:00', 3),
(5, 'Protector del territorio pero obediente a comandos. Bien socializado con familia.', '2024-01-07 16:30:00', 3),
(6, 'Imita el timbre del teléfono y palabras simples. Muy activo en las mañanas.', '2024-01-09 08:00:00', 2),
(7, 'Extremadamente gentil a pesar de su tamaño. Le gusta jugar con agua.', '2024-01-11 13:45:00', 1);

-- Insert demo pet vaccinations
INSERT INTO pet_vaccinations (pet_id, vaccine_name, administered_date, next_due_date, administered_by, batch_number, vaccine_type, notes) VALUES
-- Max's vaccinations
(1, 'DHPP (Distemper, Hepatitis, Parvovirus, Parainfluenza)', '2024-01-10', '2025-01-10', 1, 'VAC2024-001', 'core', 'Sin reacciones adversas'),
(1, 'Rabies Vaccine', '2024-01-10', '2027-01-10', 1, 'RAB2024-001', 'core', 'Vacuna de 3 años'),
(1, 'Bordetella', '2024-01-10', '2024-07-10', 1, 'BOR2024-001', 'non-core', 'Para socialización en parque'),

-- Luna's vaccinations
(2, 'FVRCP (Feline Viral Rhinotracheitis, Calicivirus, Panleukopenia)', '2024-01-08', '2025-01-08', 1, 'FVRCP2024-001', 'core', 'Sin reacciones adversas'),
(2, 'Feline Rabies', '2024-01-08', '2025-01-08', 1, 'FRAB2024-001', 'core', 'Vacuna anual para gatos'),

-- Bobby's vaccinations
(3, 'DHPP', '2024-01-12', '2025-01-12', 2, 'VAC2024-002', 'core', 'Monitoreado por problemas respiratorios'),
(3, 'Rabies Vaccine', '2024-01-12', '2027-01-12', 2, 'RAB2024-002', 'core', 'Sin complicaciones'),

-- Mila's vaccinations
(4, 'FVRCP', '2024-01-05', '2025-01-05', 3, 'FVRCP2024-002', 'core', 'Excelente tolerancia'),
(4, 'Feline Rabies', '2024-01-05', '2025-01-05', 3, 'FRAB2024-002', 'core', 'Sin reacciones'),

-- Rex's vaccinations
(5, 'DHPP', '2024-01-07', '2025-01-07', 3, 'VAC2024-003', 'core', 'Vacuna de rutina para perro senior'),
(5, 'Rabies Vaccine', '2024-01-07', '2027-01-07', 3, 'RAB2024-003', 'core', 'Tolera bien las vacunas'),

-- Simba's vaccinations
(7, 'FVRCP', '2024-01-11', '2025-01-11', 1, 'FVRCP2024-003', 'core', 'Gato grande, dosis estándar');

-- Insert demo pet deworming
INSERT INTO pet_deworming (pet_id, medication_name, administered_date, next_due_date, administered_by, notes) VALUES
(1, 'Drontal Plus', '2024-01-10', '2024-04-10', 1, 'Desparasitación de rutina'),
(2, 'Profender', '2024-01-08', '2024-04-08', 1, 'Tratamiento tópico, bien tolerado'),
(3, 'Panacur', '2024-01-12', '2024-04-12', 2, 'Debido a sensibilidad digestiva'),
(4, 'Revolution Plus', '2024-01-05', '2024-04-05', 3, 'Tratamiento mensual combinado'),
(5, 'Drontal Plus', '2024-01-07', '2024-04-07', 3, 'Dosis ajustada por peso y edad'),
(7, 'Strongid', '2024-01-11', '2024-04-11', 1, 'Bien tolerado por gato grande');

-- Insert demo appointments
INSERT INTO appointments (clinic_service, scheduled_date, status, notes, customer_id, pet_id, employee_id) VALUES
-- Future appointments
('wellness_exam', '2024-01-20 10:00:00+00:00', 'confirmed', 'Chequeo de rutina anual', 1, 1, 1),
('vaccination', '2024-01-18 14:30:00+00:00', 'confirmed', 'Refuerzo de vacunas', 2, 3, 2),
('general_consultation', '2024-01-22 09:15:00+00:00', 'pending', 'Revisión de comportamiento', 4, 6, 2),
('dental_care', '2024-01-25 11:00:00+00:00', 'confirmed', 'Limpieza dental programada', 5, 7, 1),

-- Past appointments (completed)
('general_consultation', '2024-01-15 15:30:00+00:00', 'completed', 'Consulta por problemas respiratorios', 2, 3, 2),
('wellness_exam', '2024-01-12 10:00:00+00:00', 'completed', 'Chequeo senior', 3, 5, 3),
('grooming', '2024-01-10 13:00:00+00:00', 'completed', 'Corte de uñas y baño', 3, 4, 1),
('emergency_care', '2024-01-08 18:45:00+00:00', 'completed', 'Emergencia: ingesta de objeto extraño', 1, 1, 3),

-- Cancelled appointments
('surgery', '2024-01-16 08:00:00+00:00', 'cancelled', 'Esterilización pospuesta por enfermedad leve', 4, 6, 2);

-- Insert demo medical sessions
INSERT INTO medical_sessions (pet_id, customer_id, employee_id, appointment_id, clinic_service, visit_date, visit_type, diagnosis, notes, treatment, condition, weight, temperature, heart_rate, respiratory_rate, symptoms, medications, follow_up_date, is_emergency) VALUES
-- Max's emergency session
(1, 1, 3, 8, 'emergency_care', '2024-01-08 18:45:00+00:00', 'Emergency',
'Ingesta de objeto extraño (pelota de tenis pequeña)',
'Paciente presentado por vómito y letargia. Radiografía revela objeto en estómago.',
'Inducción de vómito exitosa. Objeto expulsado completamente.',
'Recuperado completamente', 28.5, 38.2, 90, 24,
'Vómito, letargia, pérdida de apetito',
'Metoclopramida para náuseas, dieta blanda por 24 horas',
'2024-01-10 10:00:00+00:00', TRUE),

-- Bobby's respiratory consultation
(3, 2, 2, 5, 'general_consultation', '2024-01-15 15:30:00+00:00', 'Follow-up',
'Síndrome braquicefálico leve',
'Evaluación de problemas respiratorios. Examen físico normal dentro de parámetros de la raza.',
'Recomendaciones de manejo ambiental. Evitar ejercicio intenso en calor.',
'Estable', 12.3, 38.0, 110, 28,
'Respiración laboriosa después de ejercicio',
'Ninguna. Manejo conservador.',
'2024-03-15 15:30:00+00:00', FALSE),

-- Rex's senior wellness exam
(5, 3, 3, 6, 'wellness_exam', '2024-01-12 10:00:00+00:00', 'Routine',
'Salud general buena para edad. Displasia de cadera leve estable.',
'Examen físico completo. Análisis de sangre dentro de parámetros normales para edad.',
'Continuar con suplemento articular. Mantener peso ideal.',
'Condición corporal 4/5', 32.1, 38.5, 85, 20,
'Rigidez matutina leve',
'Glucosamina/Condroitina, continuar régimen actual',
'2024-07-12 10:00:00+00:00', FALSE),

-- Mila's grooming session
(4, 3, 1, 7, 'grooming', '2024-01-10 13:00:00+00:00', 'Routine',
'Sesión de grooming completa',
'Corte de uñas, limpieza de oídos, cepillado intensivo. Paciente muy cooperativa.',
'Mantenimiento de pelaje. Educación al propietario sobre cepillado diario.',
'Excelente', 4.8, 38.1, 140, 30,
'Ninguno',
'Ninguna',
NULL, FALSE);

-- Insert demo payments
INSERT INTO payments (amount, currency, status, method, med_session_id, transaction_id, description, due_date, paid_at, paid_by_customer_id, invoice_id) VALUES
-- Completed payments
(850.00, 'MXN', 'completed', 'credit_card', 1, 'TXN-2024-001',
'Consulta de emergencia - Ingesta de objeto extraño',
'2024-01-08 18:45:00+00:00', '2024-01-08 19:30:00+00:00', 1, 'INV-2024-001'),

(450.00, 'MXN', 'completed', 'cash', 2, 'TXN-2024-002',
'Consulta respiratoria - Bulldog Francés',
'2024-01-15 15:30:00+00:00', '2024-01-15 16:00:00+00:00', 2, 'INV-2024-002'),

(650.00, 'MXN', 'completed', 'debit_card', 3, 'TXN-2024-003',
'Examen de bienestar senior + análisis de sangre',
'2024-01-12 10:00:00+00:00', '2024-01-12 11:15:00+00:00', 3, 'INV-2024-003'),

(280.00, 'MXN', 'completed', 'cash', 4, 'TXN-2024-004',
'Sesión de grooming - Gato Persa',
'2024-01-10 13:00:00+00:00', '2024-01-10 13:45:00+00:00', 3, 'INV-2024-004'),

-- Pending payments (for future appointments)
(320.00, 'MXN', 'pending', 'cash', NULL, NULL,
'Examen de bienestar - Max (Golden Retriever)',
'2024-01-20 10:00:00+00:00', NULL, 1, 'INV-2024-005'),

(180.00, 'MXN', 'pending', 'credit_card', NULL, NULL,
'Vacunación - Bobby (Bulldog Francés)',
'2024-01-18 14:30:00+00:00', NULL, 2, 'INV-2024-006'),

(1200.00, 'MXN', 'pending', 'bank_transfer', NULL, NULL,
'Limpieza dental - Simba (Maine Coon)',
'2024-01-25 11:00:00+00:00', NULL, 5, 'INV-2024-007');

-- Insert demo chip implants (for pets that have microchips)
INSERT INTO pet_chip_implants (pet_id, implant_date, implant_location, chip_number, implanted_by, notes) VALUES
(1, '2019-06-15', 'Entre omóplatos', '982000123456789', 1, 'Implante realizado cuando era cachorro'),
(2, '2021-08-20', 'Entre omóplatos', '982000123456790', 2, 'Implante de rutina en gata joven'),
(3, '2022-11-10', 'Entre omóplatos', '982000123456791', 1, 'Implante antes de esterilización'),
(4, '2020-03-25', 'Entre omóplatos', '982000123456792', 3, 'Implante en gata adulta rescatada'),
(5, '2017-05-12', 'Entre omóplatos', '982000123456793', 2, 'Implante en perro joven adoptado'),
(7, '2018-09-30', 'Entre omóplatos', '982000123456794', 1, 'Implante en Maine Coon de 6 meses');

-- Update timestamps to current time for more realistic data
UPDATE users SET created_at = CURRENT_TIMESTAMP - INTERVAL '30 days' + (RANDOM() * INTERVAL '30 days');
UPDATE customers SET created_at = CURRENT_TIMESTAMP - INTERVAL '60 days' + (RANDOM() * INTERVAL '60 days');
UPDATE employees SET created_at = CURRENT_TIMESTAMP - INTERVAL '365 days' + (RANDOM() * INTERVAL '365 days');
UPDATE pets SET created_at = CURRENT_TIMESTAMP - INTERVAL '180 days' + (RANDOM() * INTERVAL '180 days');

-- Create a sample notification log entry (if you have notifications table)
-- This would typically be in MongoDB but showing structure for reference
/*
Sample MongoDB notification document:
{
  "_id": ObjectId(),
  "type": "appointment_reminder",
  "recipient": {
    "customerId": 1,
    "email": "juan.lopez@email.com",
    "phone": "+52-555-1001"
  },
  "content": {
    "subject": "Recordatorio: Cita veterinaria para Max",
    "message": "Le recordamos que tiene una cita programada para Max el 20 de enero a las 10:00 AM con Dr. Martínez."
  },
  "status": "sent",
  "sentAt": ISODate("2024-01-19T08:00:00Z"),
  "createdAt": ISODate("2024-01-19T08:00:00Z")
}
*/
