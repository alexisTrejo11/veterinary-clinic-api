package medical

import (
	"fmt"
	"time"
)

// ─── Vitals ──────────────────────────────────────────────────────────────────

// Vitals groups the physical measurements taken during a session.
// All fields are optional since not every visit type requires all of them.
type Vitals struct {
	WeightKg        *float64
	TemperatureC    *float64
	HeartRate       *int
	RespiratoryRate *int
}

// ─── MedicalSession ──────────────────────────────────────────────────────────

// MedicalSession is the central aggregate. Every clinical encounter is
// modelled as a session. Extension records (vaccinations, surgeries, etc.)
// attach to it via SessionID and are loaded/stored independently.
type MedicalSession struct {
	ID            SessionID
	PetID         uint
	CustomerID    uint
	EmployeeID    uint
	AppointmentID *uint

	ClinicService ClinicService
	VisitType     string
	VisitDate     time.Time
	IsEmergency   bool

	// Clinical content
	Symptoms  *string
	Condition *string
	Diagnosis *string
	Treatment *string
	Notes     *string

	// Vitals snapshot
	Vitals Vitals

	// Medications are free-text here; structured prescriptions live in
	// SessionPrescription. This field is kept for quick unstructured notes.
	Medications *string

	FollowUpDate *time.Time

	// Loaded extensions — populated by the application layer as needed.
	Vaccinations  []SessionVaccination
	Surgeries     []SessionSurgery
	Prescriptions []SessionPrescription
	Attachments   []SessionAttachment
	Services      []SessionService

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

// IsDeleted reports whether the session has been soft-deleted.
func (s MedicalSession) IsDeleted() bool { return s.DeletedAt != nil }

// HasExtension returns true when the session carries at least one record
// of the kind matching the current ClinicService.
func (s MedicalSession) HasExtension() bool {
	return len(s.Vaccinations) > 0 ||
		len(s.Surgeries) > 0 ||
		len(s.Prescriptions) > 0 ||
		len(s.Attachments) > 0 ||
		len(s.Services) > 0
}

// Validate checks domain invariants before persisting.
func (s MedicalSession) Validate() error {
	if s.PetID == 0 {
		return fmt.Errorf("medical session: pet_id is required")
	}
	if s.CustomerID == 0 {
		return fmt.Errorf("medical session: customer_id is required")
	}
	if s.EmployeeID == 0 {
		return fmt.Errorf("medical session: employee_id is required")
	}
	if !s.ClinicService.IsValid() {
		return fmt.Errorf("medical session: invalid clinic_service %q", s.ClinicService)
	}
	if s.VisitDate.IsZero() {
		return fmt.Errorf("medical session: visit_date is required")
	}
	return nil
}

// ─── SessionVaccination ───────────────────────────────────────────────────────

// SessionVaccination records a single vaccine dose administered during a session.
type SessionVaccination struct {
	ID               VaccinationID
	SessionID        SessionID
	VaccineCatalogID VaccineCatalogID

	BatchNumber     *string
	DoseNumber      int
	ExpirationDate  *time.Time
	SiteOfInjection *string
	NextDoseDate    *time.Time
	ReactionNotes   *string
	AdministeredBy  *uint // employee_id; may differ from session's vet

	CreatedAt time.Time
}

func (v SessionVaccination) Validate() error {
	if v.SessionID.IsZero() {
		return fmt.Errorf("session vaccination: session_id is required")
	}
	if v.VaccineCatalogID.Value() == 0 {
		return fmt.Errorf("session vaccination: vaccine_catalog_id is required")
	}
	if v.DoseNumber < 1 {
		return fmt.Errorf("session vaccination: dose_number must be >= 1")
	}
	return nil
}

// ─── SessionSurgery ───────────────────────────────────────────────────────────

// SessionSurgery holds the surgical detail for a session whose ClinicService
// is "surgery". More than one surgery may occur in a single session (rare but
// valid, e.g. combined procedures).
type SessionSurgery struct {
	ID        SurgeryID
	SessionID SessionID

	ProcedureName   string
	AnesthesiaType  *string
	AnesthesiaAgent *string
	PreOpNotes      *string
	IntraOpNotes    *string
	PostOpNotes     *string
	DurationMinutes *int
	Outcome         SurgeryOutcome
	SurgeonID       *uint

	CreatedAt time.Time
}

func (s SessionSurgery) Validate() error {
	if s.SessionID.IsZero() {
		return fmt.Errorf("session surgery: session_id is required")
	}
	if s.ProcedureName == "" {
		return fmt.Errorf("session surgery: procedure_name is required")
	}
	return nil
}

// ─── SessionPrescription ─────────────────────────────────────────────────────

// SessionPrescription links a medication from the catalog to a session,
// with dosage, frequency, and duration details.
type SessionPrescription struct {
	ID           PrescriptionID
	SessionID    SessionID
	MedicationID MedicationID

	Dosage       string
	Frequency    string
	DurationDays *int
	Route        *string // oral, topical, IM, SC…
	Instructions *string
	StartDate    time.Time
	EndDate      *time.Time // derived: start + duration

	CreatedAt time.Time
}

func (p SessionPrescription) Validate() error {
	if p.SessionID.IsZero() {
		return fmt.Errorf("session prescription: session_id is required")
	}
	if p.MedicationID.Value() == 0 {
		return fmt.Errorf("session prescription: medication_id is required")
	}
	if p.Dosage == "" {
		return fmt.Errorf("session prescription: dosage is required")
	}
	if p.Frequency == "" {
		return fmt.Errorf("session prescription: frequency is required")
	}
	return nil
}

// ─── SessionAttachment ───────────────────────────────────────────────────────

// SessionAttachment represents a file (image, lab result, PDF, etc.) linked
// to a session.
type SessionAttachment struct {
	ID          AttachmentID
	SessionID   SessionID
	FileType    AttachmentFileType
	FileURL     string
	Description *string
	UploadedBy  *uint

	CreatedAt time.Time
}

func (a SessionAttachment) Validate() error {
	if a.SessionID.IsZero() {
		return fmt.Errorf("session attachment: session_id is required")
	}
	if a.FileURL == "" {
		return fmt.Errorf("session attachment: file_url is required")
	}
	return nil
}

// ─── SessionService ──────────────────────────────────────────────────────────

// SessionService records a billable service applied during a session.
// price_applied may differ from the catalog's base_price (discounts, etc.).
type SessionService struct {
	ID               SessionServiceID
	SessionID        SessionID
	ServiceCatalogID ServiceCatalogID

	Quantity     float64
	PriceApplied *float64
	Notes        *string

	CreatedAt time.Time
}

func (s SessionService) Validate() error {
	if s.SessionID.IsZero() {
		return fmt.Errorf("session service: session_id is required")
	}
	if s.ServiceCatalogID.Value() == 0 {
		return fmt.Errorf("session service: service_catalog_id is required")
	}
	if s.Quantity <= 0 {
		return fmt.Errorf("session service: quantity must be greater than zero")
	}
	return nil
}
