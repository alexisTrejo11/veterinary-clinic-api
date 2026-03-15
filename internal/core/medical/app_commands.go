package medical

import "time"

// ─── MedicalSession commands ─────────────────────────────────────────────────

// CreateSessionCommand carries the data needed to open a new medical session.
// Extension commands are submitted separately once the session exists.
type CreateSessionCommand struct {
	PetID         uint
	CustomerID    uint
	EmployeeID    uint
	AppointmentID *uint

	ClinicService ClinicService
	VisitType     string
	VisitDate     time.Time
	IsEmergency   bool

	Symptoms  *string
	Condition *string
	Diagnosis *string
	Treatment *string
	Notes     *string

	Vitals Vitals

	Medications  *string
	FollowUpDate *time.Time
}

// UpdateSessionCommand carries updatable fields for an existing session.
// A nil pointer means "leave unchanged" at the service level.
type UpdateSessionCommand struct {
	ID SessionID

	EmployeeID    *uint
	AppointmentID *uint
	ClinicService *ClinicService
	VisitType     *string
	VisitDate     *time.Time
	IsEmergency   *bool

	Symptoms  *string
	Condition *string
	Diagnosis *string
	Treatment *string
	Notes     *string

	Vitals *Vitals

	Medications  *string
	FollowUpDate *time.Time
}

// ─── SessionVaccination commands ─────────────────────────────────────────────

type AddVaccinationCommand struct {
	SessionID        SessionID
	VaccineCatalogID VaccineCatalogID

	BatchNumber     *string
	DoseNumber      int
	ExpirationDate  *time.Time
	SiteOfInjection *string
	NextDoseDate    *time.Time
	ReactionNotes   *string
	AdministeredBy  *uint
}

type UpdateVaccinationCommand struct {
	ID              VaccinationID
	BatchNumber     *string
	SiteOfInjection *string
	NextDoseDate    *time.Time
	ReactionNotes   *string
}

// ─── SessionSurgery commands ──────────────────────────────────────────────────

type AddSurgeryCommand struct {
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
}

type UpdateSurgeryCommand struct {
	ID              SurgeryID
	AnesthesiaType  *string
	AnesthesiaAgent *string
	PreOpNotes      *string
	IntraOpNotes    *string
	PostOpNotes     *string
	DurationMinutes *int
	Outcome         *SurgeryOutcome
}

// ─── SessionPrescription commands ────────────────────────────────────────────

type AddPrescriptionCommand struct {
	SessionID    SessionID
	MedicationID MedicationID

	Dosage       string
	Frequency    string
	DurationDays *int
	Route        *string
	Instructions *string
	StartDate    time.Time
}

type UpdatePrescriptionCommand struct {
	ID           PrescriptionID
	Dosage       *string
	Frequency    *string
	DurationDays *int
	Route        *string
	Instructions *string
}

// ─── SessionAttachment commands ───────────────────────────────────────────────

type AddAttachmentCommand struct {
	SessionID   SessionID
	FileType    AttachmentFileType
	FileURL     string
	Description *string
	UploadedBy  *uint
}

// ─── SessionService commands ──────────────────────────────────────────────────

type AddSessionServiceCommand struct {
	SessionID        SessionID
	ServiceCatalogID ServiceCatalogID

	Quantity     float64
	PriceApplied *float64
	Notes        *string
}

// ─── Catalog commands ─────────────────────────────────────────────────────────

type CreateVaccineCatalogCommand struct {
	Name          string
	Manufacturer  *string
	Species       *string
	DiseaseTarget *string
	TotalDoses    int
	ScheduleDays  []int
	Notes         *string
}

type CreateMedicationCommand struct {
	Name                 string
	ActiveIngredient     *string
	Presentation         *string
	Unit                 *string
	RequiresPrescription bool
	SpeciesWarnings      *string
}

type CreateServiceCatalogCommand struct {
	Name            string
	Category        ServiceCategory
	Description     *string
	BasePrice       *float64
	DurationMinutes *int
	RequiresFasting bool
}
