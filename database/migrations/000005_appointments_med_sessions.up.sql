-- 000005_appointments_med_sessions.up.sql
-- Appointments and medical sessions

CREATE TABLE IF NOT EXISTS appointments(
    id SERIAL PRIMARY KEY,
    clinic_service VARCHAR(50) NOT NULL CHECK (clinic_service IN (
        'general_consultation', 'vaccination', 'surgery', 'dental_care', 'emergency_care',
        'grooming', 'nutrition_consult', 'behavior_consult', 'wellness_exam', 'other'
    )),
    scheduled_date TIMESTAMP WITH TIME ZONE NOT NULL,
    status appointment_status NOT NULL,
    notes TEXT DEFAULT '',
    customer_id INT NOT NULL,
    pet_id INT NOT NULL,
    employee_id INT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE NULL,
    FOREIGN KEY (pet_id) REFERENCES pets(id) ON DELETE CASCADE,
    FOREIGN KEY (customer_id) REFERENCES customers(id) ON DELETE CASCADE,
    FOREIGN KEY (employee_id) REFERENCES employees(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS medical_sessions (
    id SERIAL PRIMARY KEY,
    pet_id INT NOT NULL,
    customer_id INT NOT NULL,
    employee_id INT NOT NULL,
    appointment_id INT,
    clinic_service VARCHAR(50) NOT NULL CHECK (clinic_service IN (
        'general_consultation', 'vaccination', 'surgery', 'dental_care', 'emergency_care',
        'grooming', 'nutrition_consult', 'behavior_consult', 'wellness_exam', 'other'
    )),
    visit_date TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    visit_type VARCHAR(50) NOT NULL DEFAULT '',
    diagnosis TEXT,
    notes TEXT,
    treatment TEXT,
    condition VARCHAR(254),
    weight DECIMAL(4,2),
    temperature DECIMAL(3,1),
    heart_rate INT,
    respiratory_rate INT,
    symptoms TEXT,
    medications TEXT,
    follow_up_date TIMESTAMP WITH TIME ZONE,
    is_emergency BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE NULL,
    FOREIGN KEY (pet_id) REFERENCES pets(id) ON DELETE CASCADE,
    FOREIGN KEY (customer_id) REFERENCES customers(id) ON DELETE CASCADE,
    FOREIGN KEY (employee_id) REFERENCES employees(id) ON DELETE CASCADE,
    FOREIGN KEY (appointment_id) REFERENCES appointments(id) ON DELETE SET NULL
);


-- ─────────────────────────────────────────
-- REFERENCE CATALOGS
-- ─────────────────────────────────────────

CREATE TABLE IF NOT EXISTS vaccine_catalog (
    id              SERIAL PRIMARY KEY,
    name            VARCHAR(120) NOT NULL,
    manufacturer    VARCHAR(120),
    species         VARCHAR(50),                       -- dog, cat, rabbit, etc.
    disease_target  VARCHAR(120),                      -- Rabies, Parvovirus, etc.
    total_doses     INT DEFAULT 1,
    schedule_days   INT[],                             -- e.g. {0, 21, 365} days from dose 1
    notes           TEXT,
    is_active       BOOLEAN DEFAULT TRUE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS medications (
    id                    SERIAL PRIMARY KEY,
    name                  VARCHAR(120) NOT NULL,
    active_ingredient     VARCHAR(120),
    presentation          VARCHAR(80),                 -- tablet, injectable, syrup…
    unit                  VARCHAR(30),                 -- mg, ml, IU
    requires_prescription BOOLEAN DEFAULT FALSE,
    species_warnings      TEXT,
    is_active             BOOLEAN DEFAULT TRUE,
    created_at            TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS service_catalog (
    id               SERIAL PRIMARY KEY,
    name             VARCHAR(120) NOT NULL,
    category         VARCHAR(50) NOT NULL              -- consultation, surgical, grooming…
                     CHECK (category IN (
                         'consultation','vaccination','surgery','dental',
                         'grooming','laboratory','imaging','emergency',
                         'nutrition','behavior','wellness','other'
                     )),
    description      TEXT,
    base_price       NUMERIC(10,2),
    duration_min     INT,                              -- estimated duration
    requires_fasting BOOLEAN DEFAULT FALSE,
    is_active        BOOLEAN DEFAULT TRUE,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW()
);


-- ─────────────────────────────────────────
-- SESSION EXTENSIONS (all optional)
-- ─────────────────────────────────────────

-- Record of vaccines administered in the session
CREATE TABLE IF NOT EXISTS session_vaccinations (
    id                 SERIAL PRIMARY KEY,
    session_id         INT NOT NULL REFERENCES medical_sessions(id) ON DELETE CASCADE,
    vaccine_catalog_id INT NOT NULL REFERENCES vaccine_catalog(id),
    batch_number       VARCHAR(60),
    dose_number        INT NOT NULL DEFAULT 1,
    expiration_date    DATE,
    site_of_injection  VARCHAR(60),                    -- subcutaneous/IM, anatomical site
    next_dose_date     DATE,
    reaction_notes     TEXT,
    administered_by    INT REFERENCES employees(id),   -- may differ from primary vet
    created_at         TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Surgical detail
CREATE TABLE IF NOT EXISTS session_surgeries (
    id                    SERIAL PRIMARY KEY,
    session_id            INT NOT NULL REFERENCES medical_sessions(id) ON DELETE CASCADE,
    procedure_name        VARCHAR(150) NOT NULL,
    anesthesia_type       VARCHAR(80),
    anesthesia_agent      VARCHAR(120),
    pre_op_notes          TEXT,
    intra_op_notes        TEXT,
    post_op_notes         TEXT,
    duration_minutes      INT,
    outcome               VARCHAR(50) CHECK (outcome IN ('successful','complicated','aborted','pending')),
    surgeon_id            INT REFERENCES employees(id),
    created_at            TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Prescriptions issued in the session (may be N per session)
CREATE TABLE IF NOT EXISTS session_prescriptions (
    id              SERIAL PRIMARY KEY,
    session_id      INT NOT NULL REFERENCES medical_sessions(id) ON DELETE CASCADE,
    medication_id   INT NOT NULL REFERENCES medications(id),
    dosage          VARCHAR(80) NOT NULL,              -- e.g. "5 mg/kg"
    frequency       VARCHAR(80) NOT NULL,              -- e.g. "every 12h"
    duration_days   INT,
    route           VARCHAR(50),                       -- oral, topical, IM, SC…
    instructions    TEXT,
    start_date      DATE DEFAULT CURRENT_DATE,
    end_date        DATE GENERATED ALWAYS AS
                        (start_date + duration_days) STORED,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Attachments: X-rays, lab results, photos, PDFs
CREATE TABLE IF NOT EXISTS session_attachments (
    id          SERIAL PRIMARY KEY,
    session_id  INT NOT NULL REFERENCES medical_sessions(id) ON DELETE CASCADE,
    file_type   VARCHAR(40) NOT NULL
                CHECK (file_type IN ('image','xray','lab_result','ecg','pdf','other')),
    file_url    TEXT NOT NULL,
    description VARCHAR(255),
    uploaded_by INT REFERENCES employees(id),
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Billable services charged in the session (bridge to service_catalog)
CREATE TABLE IF NOT EXISTS session_services (
    id                 SERIAL PRIMARY KEY,
    session_id         INT NOT NULL REFERENCES medical_sessions(id) ON DELETE CASCADE,
    service_catalog_id INT NOT NULL REFERENCES service_catalog(id),
    quantity           NUMERIC(6,2) NOT NULL DEFAULT 1,
    price_applied      NUMERIC(10,2),                 -- may differ from base_price
    notes              TEXT,
    created_at         TIMESTAMPTZ NOT NULL DEFAULT NOW()
);


-- ─────────────────────────────────────────
-- VIEW: vaccination schedule per pet
-- ─────────────────────────────────────────

CREATE OR REPLACE VIEW pet_vaccination_history AS
SELECT
    ms.pet_id,
    sv.id                    AS vaccination_id,
    ms.id                    AS session_id,
    ms.visit_date,
    vc.name                  AS vaccine_name,
    vc.disease_target,
    sv.dose_number,
    sv.batch_number,
    sv.next_dose_date,
    sv.reaction_notes,
    ms.employee_id           AS vet_id
FROM session_vaccinations sv
JOIN medical_sessions  ms ON ms.id = sv.session_id
JOIN vaccine_catalog   vc ON vc.id = sv.vaccine_catalog_id
ORDER BY ms.pet_id, ms.visit_date;