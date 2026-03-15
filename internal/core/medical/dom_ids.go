package medical

import (
	"clinic-vet-api/internal/shared"
)

// ─── Value-object IDs ────────────────────────────────────────────────────────

// SessionID is the typed identifier for a MedicalSession.
type SessionID struct{ shared.IntegerID }

// VaccinationID is the typed identifier for a SessionVaccination.
type VaccinationID struct{ shared.IntegerID }

// SurgeryID is the typed identifier for a SessionSurgery.
type SurgeryID struct{ shared.IntegerID }

// PrescriptionID is the typed identifier for a SessionPrescription.
type PrescriptionID struct{ shared.IntegerID }

// MedicationID is the typed identifier for a Medication catalog entry.
type MedicationID struct{ shared.IntegerID }

// ServiceCatalogID is the typed identifier for a ServiceCatalog entry.
type ServiceCatalogID struct{ shared.IntegerID }

// VaccineCatalogID is the typed identifier for a VaccineCatalog entry.
type VaccineCatalogID struct{ shared.IntegerID }

type AttachmentID struct{ shared.IntegerID }

// SessionServiceID is the typed identifier for a SessionService.
type SessionServiceID struct{ shared.IntegerID }

func NewSessionServiceID(value uint) SessionServiceID {
	return SessionServiceID{shared.NewBaseID(value)}
}

func NewAttachmentID(value uint) AttachmentID {
	return AttachmentID{shared.NewBaseID(value)}
}

func NewVaccineCatalogID(value uint) VaccineCatalogID {
	return VaccineCatalogID{shared.NewBaseID(value)}
}

func NewSessionID(value uint) SessionID {
	return SessionID{shared.NewBaseID(value)}
}

func NewVaccinationID(value uint) VaccinationID {
	return VaccinationID{shared.NewBaseID(value)}
}

func NewSurgeryID(value uint) SurgeryID {
	return SurgeryID{shared.NewBaseID(value)}
}

func NewPrescriptionID(value uint) PrescriptionID {
	return PrescriptionID{shared.NewBaseID(value)}
}

func NewMedicationID(value uint) MedicationID {
	return MedicationID{shared.NewBaseID(value)}
}

func NewServiceCatalogID(value uint) ServiceCatalogID {
	return ServiceCatalogID{shared.NewBaseID(value)}
}
