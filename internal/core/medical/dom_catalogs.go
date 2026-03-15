package medical

import "time"

// ─── VaccineCatalog ──────────────────────────────────────────────────────────

// VaccineCatalog holds the master list of available vaccines.
// ScheduleDays defines the recommended inter-dose intervals in days from
// dose 1 (e.g. []int{0, 21, 365}).
type VaccineCatalog struct {
	ID            VaccineCatalogID
	Name          string
	Manufacturer  *string
	Species       *string // dog, cat, rabbit…
	DiseaseTarget *string
	TotalDoses    int
	ScheduleDays  []int
	Notes         *string
	IsActive      bool
	CreatedAt     time.Time
}

// ─── Medication ──────────────────────────────────────────────────────────────

// Medication is a catalog entry for a drug or treatment product.
type Medication struct {
	ID                   MedicationID
	Name                 string
	ActiveIngredient     *string
	Presentation         *string // tablet, injectable, syrup…
	Unit                 *string // mg, ml, UI
	RequiresPrescription bool
	SpeciesWarnings      *string
	IsActive             bool
	CreatedAt            time.Time
}

// ─── ServiceCatalog ──────────────────────────────────────────────────────────

// ServiceCatalog defines the services the clinic offers and their base pricing.
type ServiceCatalog struct {
	ID              ServiceCatalogID
	Name            string
	Category        ServiceCategory
	Description     *string
	BasePrice       *float64
	DurationMinutes *int
	RequiresFasting bool
	IsActive        bool
	CreatedAt       time.Time
}
