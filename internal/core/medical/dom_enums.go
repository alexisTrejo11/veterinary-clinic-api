package medical

import "fmt"

// ─── ClinicService ───────────────────────────────────────────────────────────

// ClinicService represents the high-level category of a medical session,
// mirroring the CHECK constraint in the medical_sessions table.
type ClinicService string

const (
	ClinicServiceGeneralConsultation ClinicService = "general_consultation"
	ClinicServiceVaccination         ClinicService = "vaccination"
	ClinicServiceSurgery             ClinicService = "surgery"
	ClinicServiceDentalCare          ClinicService = "dental_care"
	ClinicServiceEmergencyCare       ClinicService = "emergency_care"
	ClinicServiceGrooming            ClinicService = "grooming"
	ClinicServiceNutritionConsult    ClinicService = "nutrition_consult"
	ClinicServiceBehaviorConsult     ClinicService = "behavior_consult"
	ClinicServiceWellnessExam        ClinicService = "wellness_exam"
	ClinicServiceOther               ClinicService = "other"
)

var validClinicServices = map[ClinicService]struct{}{
	ClinicServiceGeneralConsultation: {},
	ClinicServiceVaccination:         {},
	ClinicServiceSurgery:             {},
	ClinicServiceDentalCare:          {},
	ClinicServiceEmergencyCare:       {},
	ClinicServiceGrooming:            {},
	ClinicServiceNutritionConsult:    {},
	ClinicServiceBehaviorConsult:     {},
	ClinicServiceWellnessExam:        {},
	ClinicServiceOther:               {},
}

func ParseClinicService(s string) (ClinicService, error) {
	cs := ClinicService(s)
	if _, ok := validClinicServices[cs]; !ok {
		return "", fmt.Errorf("invalid clinic service: %q", s)
	}
	return cs, nil
}

func (c ClinicService) IsValid() bool {
	_, ok := validClinicServices[c]
	return ok
}

// ─── SurgeryOutcome ──────────────────────────────────────────────────────────

type SurgeryOutcome string

const (
	SurgeryOutcomeSuccessful  SurgeryOutcome = "successful"
	SurgeryOutcomeComplicated SurgeryOutcome = "complicated"
	SurgeryOutcomeAborted     SurgeryOutcome = "aborted"
	SurgeryOutcomePending     SurgeryOutcome = "pending"
)

func ParseSurgeryOutcome(s string) (SurgeryOutcome, error) {
	switch SurgeryOutcome(s) {
	case SurgeryOutcomeSuccessful, SurgeryOutcomeComplicated,
		SurgeryOutcomeAborted, SurgeryOutcomePending:
		return SurgeryOutcome(s), nil
	}
	return "", fmt.Errorf("invalid surgery outcome: %q", s)
}

// ─── AttachmentFileType ──────────────────────────────────────────────────────

type AttachmentFileType string

const (
	AttachmentFileTypeImage     AttachmentFileType = "image"
	AttachmentFileTypeXRay      AttachmentFileType = "xray"
	AttachmentFileTypeLabResult AttachmentFileType = "lab_result"
	AttachmentFileTypeECG       AttachmentFileType = "ecg"
	AttachmentFileTypePDF       AttachmentFileType = "pdf"
	AttachmentFileTypeOther     AttachmentFileType = "other"
)

func ParseAttachmentFileType(s string) (AttachmentFileType, error) {
	switch AttachmentFileType(s) {
	case AttachmentFileTypeImage, AttachmentFileTypeXRay, AttachmentFileTypeLabResult,
		AttachmentFileTypeECG, AttachmentFileTypePDF, AttachmentFileTypeOther:
		return AttachmentFileType(s), nil
	}
	return "", fmt.Errorf("invalid attachment file type: %q", s)
}

// ─── ServiceCategory ─────────────────────────────────────────────────────────

type ServiceCategory string

const (
	ServiceCategoryConsultation ServiceCategory = "consultation"
	ServiceCategoryVaccination  ServiceCategory = "vaccination"
	ServiceCategorySurgery      ServiceCategory = "surgery"
	ServiceCategoryDental       ServiceCategory = "dental"
	ServiceCategoryGrooming     ServiceCategory = "grooming"
	ServiceCategoryLaboratory   ServiceCategory = "laboratory"
	ServiceCategoryImaging      ServiceCategory = "imaging"
	ServiceCategoryEmergency    ServiceCategory = "emergency"
	ServiceCategoryNutrition    ServiceCategory = "nutrition"
	ServiceCategoryBehavior     ServiceCategory = "behavior"
	ServiceCategoryWellness     ServiceCategory = "wellness"
	ServiceCategoryOther        ServiceCategory = "other"
)

var validServiceCategories = map[ServiceCategory]struct{}{
	ServiceCategoryConsultation: {},
	ServiceCategoryVaccination:  {},
	ServiceCategorySurgery:      {},
	ServiceCategoryDental:       {},
	ServiceCategoryGrooming:     {},
	ServiceCategoryLaboratory:   {},
	ServiceCategoryImaging:     {},
	ServiceCategoryEmergency:   {},
	ServiceCategoryNutrition:    {},
	ServiceCategoryBehavior:   {},
	ServiceCategoryWellness:   {},
	ServiceCategoryOther:      {},
}

func ParseServiceCategory(s string) (ServiceCategory, error) {
	sc := ServiceCategory(s)
	if _, ok := validServiceCategories[sc]; !ok {
		return "", fmt.Errorf("invalid service category: %q", s)
	}
	return sc, nil
}
