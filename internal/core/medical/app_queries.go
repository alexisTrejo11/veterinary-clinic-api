package medical

// ─── Stats / query results ───────────────────────────────────────────────────

// SessionStats holds aggregate counts for dashboards and reporting.
type SessionStats struct {
	Total           int64
	Active          int64
	Emergency       int64
	ByClinicService map[ClinicService]int64
}

// PetVaccinationSummary is a read-model representing the vaccination state of
// a single pet — which vaccines are up to date and which are overdue.
type PetVaccinationSummary struct {
	PetID        uint
	Vaccinations []VaccinationStatus
}

// VaccinationStatus shows whether a specific vaccine is current for a pet.
type VaccinationStatus struct {
	VaccineName   string
	DiseaseTarget string
	LastDoseDate  *interface{} // kept as any to avoid circular time import
	NextDoseDate  *interface{}
	IsOverdue     bool
}

// SessionWithExtensions is a read-model that bundles a session together with
// all its extension records. Used by GetSessionFull on the service layer.
type SessionWithExtensions struct {
	Session       MedicalSession
	Vaccinations  []SessionVaccination
	Surgeries     []SessionSurgery
	Prescriptions []SessionPrescription
	Attachments   []SessionAttachment
	Services      []SessionService
}
