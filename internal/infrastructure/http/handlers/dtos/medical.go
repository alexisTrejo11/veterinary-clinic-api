package dtos

import (
	"clinic-vet-api/internal/shared/page"
)

// ─── Session ─────────────────────────────────────────────────────────────────

type CreateSessionRequest struct {
	PetID         uint     `json:"pet_id" binding:"required"`
	CustomerID    uint     `json:"customer_id" binding:"required"`
	EmployeeID    uint     `json:"employee_id" binding:"required"`
	AppointmentID *uint    `json:"appointment_id"`
	ClinicService string   `json:"clinic_service" binding:"required"`
	VisitType     string   `json:"visit_type" binding:"required"`
	VisitDate     string   `json:"visit_date" binding:"required"` // RFC3339
	IsEmergency   bool     `json:"is_emergency"`
	Symptoms      *string  `json:"symptoms"`
	Condition     *string  `json:"condition"`
	Diagnosis     *string  `json:"diagnosis"`
	Treatment     *string  `json:"treatment"`
	Notes         *string  `json:"notes"`
	Medications   *string  `json:"medications"`
	FollowUpDate  *string  `json:"follow_up_date"` // RFC3339
	WeightKg      *float64 `json:"weight_kg"`
	TemperatureC  *float64 `json:"temperature_c"`
	HeartRate     *int     `json:"heart_rate"`
	RespiratoryRate *int   `json:"respiratory_rate"`
}

type UpdateSessionRequest struct {
	EmployeeID    *uint    `json:"employee_id"`
	AppointmentID *uint    `json:"appointment_id"`
	ClinicService *string  `json:"clinic_service"`
	VisitType     *string  `json:"visit_type"`
	VisitDate     *string  `json:"visit_date"`
	IsEmergency   *bool    `json:"is_emergency"`
	Symptoms      *string  `json:"symptoms"`
	Condition     *string  `json:"condition"`
	Diagnosis     *string  `json:"diagnosis"`
	Treatment     *string  `json:"treatment"`
	Notes         *string  `json:"notes"`
	Medications   *string  `json:"medications"`
	FollowUpDate  *string  `json:"follow_up_date"`
	WeightKg      *float64 `json:"weight_kg"`
	TemperatureC  *float64 `json:"temperature_c"`
	HeartRate     *int     `json:"heart_rate"`
	RespiratoryRate *int   `json:"respiratory_rate"`
}

type SessionSearchRequest struct {
	PetID          uint    `form:"pet_id"`
	CustomerID     uint    `form:"customer_id"`
	EmployeeID     uint    `form:"employee_id"`
	ClinicService  string  `form:"clinic_service"`
	IsEmergency    *bool   `form:"is_emergency"`
	VisitDateFrom  string  `form:"visit_date_from"`
	VisitDateTo    string  `form:"visit_date_to"`
	page.PaginationRequest
}

// ─── Vaccination ────────────────────────────────────────────────────────────

type AddVaccinationRequest struct {
	SessionID        uint    `json:"session_id" binding:"required"`
	VaccineCatalogID uint    `json:"vaccine_catalog_id" binding:"required"`
	BatchNumber      *string `json:"batch_number"`
	DoseNumber       int     `json:"dose_number" binding:"required,min=1"`
	ExpirationDate   *string `json:"expiration_date"`
	SiteOfInjection  *string `json:"site_of_injection"`
	NextDoseDate     *string `json:"next_dose_date"`
	ReactionNotes    *string `json:"reaction_notes"`
	AdministeredBy   *uint   `json:"administered_by"`
}

type UpdateVaccinationRequest struct {
	BatchNumber     *string `json:"batch_number"`
	SiteOfInjection *string `json:"site_of_injection"`
	NextDoseDate    *string `json:"next_dose_date"`
	ReactionNotes   *string `json:"reaction_notes"`
}

// ─── Surgery ────────────────────────────────────────────────────────────────

type AddSurgeryRequest struct {
	SessionID        uint    `json:"session_id" binding:"required"`
	ProcedureName    string  `json:"procedure_name" binding:"required"`
	AnesthesiaType   *string `json:"anesthesia_type"`
	AnesthesiaAgent  *string `json:"anesthesia_agent"`
	PreOpNotes       *string `json:"pre_op_notes"`
	IntraOpNotes     *string `json:"intra_op_notes"`
	PostOpNotes      *string `json:"post_op_notes"`
	DurationMinutes  *int    `json:"duration_minutes"`
	Outcome          string  `json:"outcome" binding:"required"`
	SurgeonID        *uint   `json:"surgeon_id"`
}

type UpdateSurgeryRequest struct {
	AnesthesiaType   *string `json:"anesthesia_type"`
	AnesthesiaAgent  *string `json:"anesthesia_agent"`
	PreOpNotes       *string `json:"pre_op_notes"`
	IntraOpNotes     *string `json:"intra_op_notes"`
	PostOpNotes      *string `json:"post_op_notes"`
	DurationMinutes  *int    `json:"duration_minutes"`
	Outcome          *string `json:"outcome"`
}

// ─── Prescription ───────────────────────────────────────────────────────────

type AddPrescriptionRequest struct {
	SessionID    uint   `json:"session_id" binding:"required"`
	MedicationID uint   `json:"medication_id" binding:"required"`
	Dosage       string `json:"dosage" binding:"required"`
	Frequency    string `json:"frequency" binding:"required"`
	DurationDays *int   `json:"duration_days"`
	Route        *string `json:"route"`
	Instructions *string `json:"instructions"`
	StartDate    string `json:"start_date" binding:"required"` // RFC3339
}

type UpdatePrescriptionRequest struct {
	Dosage       *string `json:"dosage"`
	Frequency    *string `json:"frequency"`
	DurationDays *int    `json:"duration_days"`
	Route        *string `json:"route"`
	Instructions *string `json:"instructions"`
}

// ─── Attachment ─────────────────────────────────────────────────────────────

type AddAttachmentRequest struct {
	SessionID   uint   `json:"session_id" binding:"required"`
	FileType    string `json:"file_type" binding:"required"`
	FileURL     string `json:"file_url" binding:"required"`
	Description *string `json:"description"`
	UploadedBy  *uint   `json:"uploaded_by"`
}

// ─── Session service ────────────────────────────────────────────────────────

type AddSessionServiceRequest struct {
	SessionID        uint     `json:"session_id" binding:"required"`
	ServiceCatalogID uint     `json:"service_catalog_id" binding:"required"`
	Quantity         float64  `json:"quantity" binding:"required,gt=0"`
	PriceApplied     *float64 `json:"price_applied"`
	Notes            *string  `json:"notes"`
}

// ─── Catalogs ───────────────────────────────────────────────────────────────

type CreateVaccineCatalogRequest struct {
	Name          string  `json:"name" binding:"required"`
	Manufacturer  *string `json:"manufacturer"`
	Species       *string `json:"species"`
	DiseaseTarget *string `json:"disease_target"`
	TotalDoses    int     `json:"total_doses" binding:"required,min=1"`
	ScheduleDays  []int   `json:"schedule_days"`
	Notes         *string `json:"notes"`
}

type CreateMedicationCatalogRequest struct {
	Name                 string `json:"name" binding:"required"`
	ActiveIngredient     *string `json:"active_ingredient"`
	Presentation         *string `json:"presentation"`
	Unit                 *string `json:"unit"`
	RequiresPrescription bool   `json:"requires_prescription"`
	SpeciesWarnings      *string `json:"species_warnings"`
}

type CreateServiceCatalogRequest struct {
	Name            string   `json:"name" binding:"required"`
	Category        string   `json:"category" binding:"required"`
	Description     *string  `json:"description"`
	BasePrice       *float64 `json:"base_price"`
	DurationMinutes *int     `json:"duration_minutes"`
	RequiresFasting bool     `json:"requires_fasting"`
}

// VaccinationHistoryRequest for query params (pet_id, vaccine_id, date_from, date_to, page, size)
type VaccinationHistoryRequest struct {
	PetID     uint   `form:"pet_id"`
	VaccineID uint   `form:"vaccine_id"`
	DateFrom  string `form:"date_from"`
	DateTo    string `form:"date_to"`
	page.PaginationRequest
}
