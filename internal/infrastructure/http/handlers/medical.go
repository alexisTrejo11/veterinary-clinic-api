package handlers

import (
	"errors"

	"clinic-vet-api/internal/core/medical"
	"clinic-vet-api/internal/core/pets"
	"clinic-vet-api/internal/infrastructure/http/handlers/dtos"
	"clinic-vet-api/internal/infrastructure/http/handlers/mappers"
	"clinic-vet-api/internal/shared/http"
	"clinic-vet-api/internal/shared/page"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// MedicalHandler exposes medical session, vaccinations, surgeries, prescriptions,
// attachments, session services, and catalogs. Customers are read-only (their data);
// employees and managers share the same write endpoints.
type MedicalHandler struct {
	sessionSvc       medical.MedicalSessionService
	vaccinationSvc   medical.VaccinationService
	surgerySvc       medical.SurgeryService
	prescriptionSvc  medical.PrescriptionService
	attachmentSvc    medical.AttachmentService
	sessionServiceSvc medical.SessionServiceManager
	vaccineCatalogSvc medical.VaccineCatalogService
	medicationCatalogSvc medical.MedicationCatalogService
	serviceCatalogSvc   medical.ServiceCatalogService
	validator        *validator.Validate
	customerIDResolver CustomerIDResolver
	petService       pets.Service // to verify pet belongs to customer
}

func NewMedicalHandler(
	sessionSvc medical.MedicalSessionService,
	vaccinationSvc medical.VaccinationService,
	surgerySvc medical.SurgeryService,
	prescriptionSvc medical.PrescriptionService,
	attachmentSvc medical.AttachmentService,
	sessionServiceSvc medical.SessionServiceManager,
	vaccineCatalogSvc medical.VaccineCatalogService,
	medicationCatalogSvc medical.MedicationCatalogService,
	serviceCatalogSvc medical.ServiceCatalogService,
	validator *validator.Validate,
	customerIDResolver CustomerIDResolver,
	petService pets.Service,
) *MedicalHandler {
	return &MedicalHandler{
		sessionSvc:          sessionSvc,
		vaccinationSvc:      vaccinationSvc,
		surgerySvc:          surgerySvc,
		prescriptionSvc:     prescriptionSvc,
		attachmentSvc:       attachmentSvc,
		sessionServiceSvc:   sessionServiceSvc,
		vaccineCatalogSvc:   vaccineCatalogSvc,
		medicationCatalogSvc: medicationCatalogSvc,
		serviceCatalogSvc:  serviceCatalogSvc,
		validator:           validator,
		customerIDResolver:  customerIDResolver,
		petService:          petService,
	}
}

func parseSessionID(c *gin.Context) (medical.SessionID, error) {
	id, err := http.ParseParamToUInt(c, "id")
	if err != nil {
		return medical.SessionID{}, err
	}
	return medical.NewSessionID(id), nil
}

func (h *MedicalHandler) getMyCustomerID(c *gin.Context) (uint, error) {
	if h.customerIDResolver == nil {
		return 0, errors.New("customer resolver not configured")
	}
	return CustomerIDFromContext(c, h.customerIDResolver)
}

// ensureSessionBelongsToCustomer loads the session and returns an error if it does not belong to the customer.
func (h *MedicalHandler) ensureSessionBelongsToCustomer(c *gin.Context, sessionID medical.SessionID, customerID uint) (medical.MedicalSession, error) {
	sess, err := h.sessionSvc.GetSessionByID(c.Request.Context(), sessionID)
	if err != nil {
		return medical.MedicalSession{}, err
	}
	if sess.CustomerID != customerID {
		return medical.MedicalSession{}, errors.New("forbidden: session does not belong to customer")
	}
	return sess, nil
}

// ensurePetBelongsToCustomer returns nil if the pet belongs to the customer.
func (h *MedicalHandler) ensurePetBelongsToCustomer(c *gin.Context, petID uint, customerID uint) error {
	_, err := h.petService.GetPetByIDAndCustomerID(c.Request.Context(), pets.NewPetID(petID), customerID)
	return err
}

// ------------------------------------------------------------
// Customer read handlers (only their data / their pets’ data)
// ------------------------------------------------------------

func (h *MedicalHandler) GetMySessions(c *gin.Context) {
	logic := func(ctx *gin.Context) (any, error) {
		customerID, err := h.getMyCustomerID(ctx)
		if err != nil {
			return nil, err
		}
		var pageParams page.PaginationRequest
		if err := http.ShouldBindPageParams(&pageParams, ctx, h.validator); err != nil {
			return nil, err
		}
		return h.sessionSvc.GetSessionsByCustomer(ctx.Request.Context(), customerID, pageParams.ToPagination())
	}
	http.HandleRequestNoBodyWithResponder(h.validator, logic, func(c *gin.Context, res any) {
		if res == nil {
			http.Success(c, nil, "No sessions found")
			return
		}
		p := res.(page.Page[medical.MedicalSession])
		http.Paginated(c, &p, "Sessions")
	})(c)
}

func (h *MedicalHandler) GetMySessionByID(c *gin.Context) {
	logic := func(ctx *gin.Context) (any, error) {
		customerID, err := h.getMyCustomerID(ctx)
		if err != nil {
			return nil, err
		}
		sessionID, err := parseSessionID(ctx)
		if err != nil {
			return nil, err
		}
		return h.ensureSessionBelongsToCustomer(ctx, sessionID, customerID)
	}
	http.HandleGetRequest(h.validator, "Session", logic)(c)
}

func (h *MedicalHandler) GetMySessionFull(c *gin.Context) {
	logic := func(ctx *gin.Context) (any, error) {
		customerID, err := h.getMyCustomerID(ctx)
		if err != nil {
			return nil, err
		}
		sessionID, err := parseSessionID(ctx)
		if err != nil {
			return nil, err
		}
		_, err = h.ensureSessionBelongsToCustomer(ctx, sessionID, customerID)
		if err != nil {
			return nil, err
		}
		return h.sessionSvc.GetSessionFull(ctx.Request.Context(), sessionID)
	}
	http.HandleGetRequest(h.validator, "Session", logic)(c)
}

func (h *MedicalHandler) GetMyPetSessions(c *gin.Context) {
	logic := func(ctx *gin.Context) (any, error) {
		customerID, err := h.getMyCustomerID(ctx)
		if err != nil {
			return nil, err
		}
		petID, err := http.ParseParamToUInt(ctx, "pet_id")
		if err != nil {
			return nil, err
		}
		if err := h.ensurePetBelongsToCustomer(ctx, petID, customerID); err != nil {
			return nil, err
		}
		var pageParams page.PaginationRequest
		if err := http.ShouldBindPageParams(&pageParams, ctx, h.validator); err != nil {
			return nil, err
		}
		return h.sessionSvc.GetSessionsByPet(ctx.Request.Context(), petID, pageParams.ToPagination())
	}
	http.HandleRequestNoBodyWithResponder(h.validator, logic, func(c *gin.Context, res any) {
		if res == nil {
			http.Success(c, nil, "No sessions found")
			return
		}
		p := res.(page.Page[medical.MedicalSession])
		http.Paginated(c, &p, "Sessions")
	})(c)
}

func (h *MedicalHandler) GetMyPetVaccinationSummary(c *gin.Context) {
	logic := func(ctx *gin.Context) (any, error) {
		customerID, err := h.getMyCustomerID(ctx)
		if err != nil {
			return nil, err
		}
		petID, err := http.ParseParamToUInt(ctx, "pet_id")
		if err != nil {
			return nil, err
		}
		if err := h.ensurePetBelongsToCustomer(ctx, petID, customerID); err != nil {
			return nil, err
		}
		return h.vaccinationSvc.GetPetVaccinationSummary(ctx.Request.Context(), petID)
	}
	http.HandleGetRequest(h.validator, "VaccinationSummary", logic)(c)
}

func (h *MedicalHandler) GetMyPetVaccinationHistory(c *gin.Context) {
	logic := func(ctx *gin.Context) (any, error) {
		customerID, err := h.getMyCustomerID(ctx)
		if err != nil {
			return nil, err
		}
		petID, err := http.ParseParamToUInt(ctx, "pet_id")
		if err != nil {
			return nil, err
		}
		if err := h.ensurePetBelongsToCustomer(ctx, petID, customerID); err != nil {
			return nil, err
		}
		var req dtos.VaccinationHistoryRequest
		if err := ctx.ShouldBindQuery(&req); err != nil {
			return nil, err
		}
		req.PetID = petID
		spec := mappers.VaccinationHistoryRequestToSpec(req)
		return h.vaccinationSvc.GetVaccinationHistory(ctx.Request.Context(), *spec)
	}
	http.HandleRequestNoBodyWithResponder(h.validator, logic, func(c *gin.Context, res any) {
		if res == nil {
			http.Success(c, nil, "No history found")
			return
		}
		p := res.(page.Page[medical.SessionVaccination])
		http.Paginated(c, &p, "VaccinationHistory")
	})(c)
}

func (h *MedicalHandler) GetMyPetActivePrescriptions(c *gin.Context) {
	logic := func(ctx *gin.Context) (any, error) {
		customerID, err := h.getMyCustomerID(ctx)
		if err != nil {
			return nil, err
		}
		petID, err := http.ParseParamToUInt(ctx, "pet_id")
		if err != nil {
			return nil, err
		}
		if err := h.ensurePetBelongsToCustomer(ctx, petID, customerID); err != nil {
			return nil, err
		}
		var pageParams page.PaginationRequest
		if err := http.ShouldBindPageParams(&pageParams, ctx, h.validator); err != nil {
			return nil, err
		}
		return h.prescriptionSvc.GetActivePrescriptionsByPet(ctx.Request.Context(), petID, pageParams.ToPagination())
	}
	http.HandleRequestNoBodyWithResponder(h.validator, logic, func(c *gin.Context, res any) {
		if res == nil {
			http.Success(c, nil, "No prescriptions found")
			return
		}
		p := res.(page.Page[medical.SessionPrescription])
		http.Paginated(c, &p, "Prescriptions")
	})(c)
}

// ------------------------------------------------------------
// Staff read handlers (employee + manager; no role split)
// ------------------------------------------------------------

func (h *MedicalHandler) GetSessionByID(c *gin.Context) {
	logic := func(ctx *gin.Context) (any, error) {
		sessionID, err := parseSessionID(ctx)
		if err != nil {
			return nil, err
		}
		return h.sessionSvc.GetSessionByID(ctx.Request.Context(), sessionID)
	}
	http.HandleGetRequest(h.validator, "Session", logic)(c)
}

func (h *MedicalHandler) GetSessionFull(c *gin.Context) {
	logic := func(ctx *gin.Context) (any, error) {
		sessionID, err := parseSessionID(ctx)
		if err != nil {
			return nil, err
		}
		return h.sessionSvc.GetSessionFull(ctx.Request.Context(), sessionID)
	}
	http.HandleGetRequest(h.validator, "Session", logic)(c)
}

func (h *MedicalHandler) GetSessionsBySpecification(c *gin.Context) {
	var req dtos.SessionSearchRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		http.BadRequest(c, err)
		return
	}
	spec, err := mappers.SessionSearchToSpec(req)
	if err != nil {
		http.BadRequest(c, err)
		return
	}
	logic := func(ctx *gin.Context) (any, error) {
		return h.sessionSvc.GetSessionsBySpecification(ctx.Request.Context(), *spec)
	}
	http.HandleRequestNoBodyWithResponder(h.validator, logic, func(c *gin.Context, res any) {
		if res == nil {
			http.Success(c, nil, "No sessions found")
			return
		}
		p := res.(page.Page[medical.MedicalSession])
		http.Paginated(c, &p, "Sessions")
	})(c)
}

func (h *MedicalHandler) GetSessionsByPet(c *gin.Context) {
	logic := func(ctx *gin.Context) (any, error) {
		petID, err := http.ParseParamToUInt(ctx, "pet_id")
		if err != nil {
			return nil, err
		}
		var pageParams page.PaginationRequest
		if err := http.ShouldBindPageParams(&pageParams, ctx, h.validator); err != nil {
			return nil, err
		}
		return h.sessionSvc.GetSessionsByPet(ctx.Request.Context(), petID, pageParams.ToPagination())
	}
	http.HandleRequestNoBodyWithResponder(h.validator, logic, func(c *gin.Context, res any) {
		if res == nil {
			http.Success(c, nil, "No sessions found")
			return
		}
		p := res.(page.Page[medical.MedicalSession])
		http.Paginated(c, &p, "Sessions")
	})(c)
}

func (h *MedicalHandler) GetSessionsByCustomer(c *gin.Context) {
	logic := func(ctx *gin.Context) (any, error) {
		customerID, err := CustomerIDFromParam(ctx)
		if err != nil {
			return nil, err
		}
		var pageParams page.PaginationRequest
		if err := http.ShouldBindPageParams(&pageParams, ctx, h.validator); err != nil {
			return nil, err
		}
		return h.sessionSvc.GetSessionsByCustomer(ctx.Request.Context(), customerID, pageParams.ToPagination())
	}
	http.HandleRequestNoBodyWithResponder(h.validator, logic, func(c *gin.Context, res any) {
		if res == nil {
			http.Success(c, nil, "No sessions found")
			return
		}
		p := res.(page.Page[medical.MedicalSession])
		http.Paginated(c, &p, "Sessions")
	})(c)
}

func (h *MedicalHandler) GetSessionStats(c *gin.Context) {
	logic := func(ctx *gin.Context) (any, error) {
		return h.sessionSvc.GetSessionStats(ctx.Request.Context())
	}
	http.HandleGetRequest(h.validator, "SessionStats", logic)(c)
}

func (h *MedicalHandler) GetVaccinationsBySession(c *gin.Context) {
	logic := func(ctx *gin.Context) (any, error) {
		sessionID, err := parseSessionID(ctx)
		if err != nil {
			return nil, err
		}
		return h.vaccinationSvc.GetVaccinationsBySession(ctx.Request.Context(), sessionID)
	}
	http.HandleGetRequest(h.validator, "Vaccinations", logic)(c)
}

func (h *MedicalHandler) GetSurgeriesBySession(c *gin.Context) {
	logic := func(ctx *gin.Context) (any, error) {
		sessionID, err := parseSessionID(ctx)
		if err != nil {
			return nil, err
		}
		return h.surgerySvc.GetSurgeriesBySession(ctx.Request.Context(), sessionID)
	}
	http.HandleGetRequest(h.validator, "Surgeries", logic)(c)
}

func (h *MedicalHandler) GetPrescriptionsBySession(c *gin.Context) {
	logic := func(ctx *gin.Context) (any, error) {
		sessionID, err := parseSessionID(ctx)
		if err != nil {
			return nil, err
		}
		return h.prescriptionSvc.GetPrescriptionsBySession(ctx.Request.Context(), sessionID)
	}
	http.HandleGetRequest(h.validator, "Prescriptions", logic)(c)
}

func (h *MedicalHandler) GetAttachmentsBySession(c *gin.Context) {
	logic := func(ctx *gin.Context) (any, error) {
		sessionID, err := parseSessionID(ctx)
		if err != nil {
			return nil, err
		}
		return h.attachmentSvc.GetAttachmentsBySession(ctx.Request.Context(), sessionID)
	}
	http.HandleGetRequest(h.validator, "Attachments", logic)(c)
}

func (h *MedicalHandler) GetServicesBySession(c *gin.Context) {
	logic := func(ctx *gin.Context) (any, error) {
		sessionID, err := parseSessionID(ctx)
		if err != nil {
			return nil, err
		}
		return h.sessionServiceSvc.GetServicesBySession(ctx.Request.Context(), sessionID)
	}
	http.HandleGetRequest(h.validator, "SessionServices", logic)(c)
}

func (h *MedicalHandler) GetVaccinationHistory(c *gin.Context) {
	var req dtos.VaccinationHistoryRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		http.BadRequest(c, err)
		return
	}
	spec := mappers.VaccinationHistoryRequestToSpec(req)
	logic := func(ctx *gin.Context) (any, error) {
		return h.vaccinationSvc.GetVaccinationHistory(ctx.Request.Context(), *spec)
	}
	http.HandleRequestNoBodyWithResponder(h.validator, logic, func(c *gin.Context, res any) {
		if res == nil {
			http.Success(c, nil, "No history found")
			return
		}
		p := res.(page.Page[medical.SessionVaccination])
		http.Paginated(c, &p, "VaccinationHistory")
	})(c)
}

func (h *MedicalHandler) GetPetVaccinationSummary(c *gin.Context) {
	logic := func(ctx *gin.Context) (any, error) {
		petID, err := http.ParseParamToUInt(ctx, "pet_id")
		if err != nil {
			return nil, err
		}
		return h.vaccinationSvc.GetPetVaccinationSummary(ctx.Request.Context(), petID)
	}
	http.HandleGetRequest(h.validator, "VaccinationSummary", logic)(c)
}

func (h *MedicalHandler) GetActivePrescriptionsByPet(c *gin.Context) {
	logic := func(ctx *gin.Context) (any, error) {
		petID, err := http.ParseParamToUInt(ctx, "pet_id")
		if err != nil {
			return nil, err
		}
		var pageParams page.PaginationRequest
		if err := http.ShouldBindPageParams(&pageParams, ctx, h.validator); err != nil {
			return nil, err
		}
		return h.prescriptionSvc.GetActivePrescriptionsByPet(ctx.Request.Context(), petID, pageParams.ToPagination())
	}
	http.HandleRequestNoBodyWithResponder(h.validator, logic, func(c *gin.Context, res any) {
		if res == nil {
			http.Success(c, nil, "No prescriptions found")
			return
		}
		p := res.(page.Page[medical.SessionPrescription])
		http.Paginated(c, &p, "Prescriptions")
	})(c)
}

// Catalog read
func (h *MedicalHandler) ListVaccines(c *gin.Context) {
	logic := func(ctx *gin.Context) (any, error) {
		var pageParams page.PaginationRequest
		if err := http.ShouldBindPageParams(&pageParams, ctx, h.validator); err != nil {
			return nil, err
		}
		return h.vaccineCatalogSvc.ListVaccines(ctx.Request.Context(), pageParams.ToPagination())
	}
	http.HandleRequestNoBodyWithResponder(h.validator, logic, func(c *gin.Context, res any) {
		if res == nil {
			http.Success(c, nil, "No vaccines found")
			return
		}
		p := res.(page.Page[medical.VaccineCatalog])
		http.Paginated(c, &p, "Vaccines")
	})(c)
}

func (h *MedicalHandler) GetVaccineByID(c *gin.Context) {
	logic := func(ctx *gin.Context) (any, error) {
		id, err := http.ParseParamToUInt(ctx, "id")
		if err != nil {
			return nil, err
		}
		return h.vaccineCatalogSvc.GetVaccineByID(ctx.Request.Context(), medical.NewVaccineCatalogID(id))
	}
	http.HandleGetRequest(h.validator, "Vaccine", logic)(c)
}

func (h *MedicalHandler) ListVaccinesBySpecies(c *gin.Context) {
	logic := func(ctx *gin.Context) (any, error) {
		species := c.Param("species")
		var pageParams page.PaginationRequest
		if err := http.ShouldBindPageParams(&pageParams, ctx, h.validator); err != nil {
			return nil, err
		}
		return h.vaccineCatalogSvc.ListVaccinesBySpecies(ctx.Request.Context(), species, pageParams.ToPagination())
	}
	http.HandleRequestNoBodyWithResponder(h.validator, logic, func(c *gin.Context, res any) {
		if res == nil {
			http.Success(c, nil, "No vaccines found")
			return
		}
		p := res.(page.Page[medical.VaccineCatalog])
		http.Paginated(c, &p, "Vaccines")
	})(c)
}

func (h *MedicalHandler) ListMedications(c *gin.Context) {
	logic := func(ctx *gin.Context) (any, error) {
		var pageParams page.PaginationRequest
		if err := http.ShouldBindPageParams(&pageParams, ctx, h.validator); err != nil {
			return nil, err
		}
		return h.medicationCatalogSvc.ListMedications(ctx.Request.Context(), pageParams.ToPagination())
	}
	http.HandleRequestNoBodyWithResponder(h.validator, logic, func(c *gin.Context, res any) {
		if res == nil {
			http.Success(c, nil, "No medications found")
			return
		}
		p := res.(page.Page[medical.Medication])
		http.Paginated(c, &p, "Medications")
	})(c)
}

func (h *MedicalHandler) SearchMedications(c *gin.Context) {
	logic := func(ctx *gin.Context) (any, error) {
		term := c.Query("term")
		var pageParams page.PaginationRequest
		if err := http.ShouldBindPageParams(&pageParams, ctx, h.validator); err != nil {
			return nil, err
		}
		return h.medicationCatalogSvc.SearchMedications(ctx.Request.Context(), term, pageParams.ToPagination())
	}
	http.HandleRequestNoBodyWithResponder(h.validator, logic, func(c *gin.Context, res any) {
		if res == nil {
			http.Success(c, nil, "No medications found")
			return
		}
		p := res.(page.Page[medical.Medication])
		http.Paginated(c, &p, "Medications")
	})(c)
}

func (h *MedicalHandler) GetMedicationByID(c *gin.Context) {
	logic := func(ctx *gin.Context) (any, error) {
		id, err := http.ParseParamToUInt(ctx, "id")
		if err != nil {
			return nil, err
		}
		return h.medicationCatalogSvc.GetMedicationByID(ctx.Request.Context(), medical.NewMedicationID(id))
	}
	http.HandleGetRequest(h.validator, "Medication", logic)(c)
}

func (h *MedicalHandler) ListServices(c *gin.Context) {
	logic := func(ctx *gin.Context) (any, error) {
		var pageParams page.PaginationRequest
		if err := http.ShouldBindPageParams(&pageParams, ctx, h.validator); err != nil {
			return nil, err
		}
		return h.serviceCatalogSvc.ListServices(ctx.Request.Context(), pageParams.ToPagination())
	}
	http.HandleRequestNoBodyWithResponder(h.validator, logic, func(c *gin.Context, res any) {
		if res == nil {
			http.Success(c, nil, "No services found")
			return
		}
		p := res.(page.Page[medical.ServiceCatalog])
		http.Paginated(c, &p, "Services")
	})(c)
}

func (h *MedicalHandler) GetServiceByID(c *gin.Context) {
	logic := func(ctx *gin.Context) (any, error) {
		id, err := http.ParseParamToUInt(ctx, "id")
		if err != nil {
			return nil, err
		}
		return h.serviceCatalogSvc.GetServiceByID(ctx.Request.Context(), medical.NewServiceCatalogID(id))
	}
	http.HandleGetRequest(h.validator, "Service", logic)(c)
}

func (h *MedicalHandler) ListServicesByCategory(c *gin.Context) {
	logic := func(ctx *gin.Context) (any, error) {
		catStr := c.Param("category")
		cat, err := medical.ParseServiceCategory(catStr)
		if err != nil {
			return nil, err
		}
		var pageParams page.PaginationRequest
		if err := http.ShouldBindPageParams(&pageParams, ctx, h.validator); err != nil {
			return nil, err
		}
		return h.serviceCatalogSvc.ListServicesByCategory(ctx.Request.Context(), cat, pageParams.ToPagination())
	}
	http.HandleRequestNoBodyWithResponder(h.validator, logic, func(c *gin.Context, res any) {
		if res == nil {
			http.Success(c, nil, "No services found")
			return
		}
		p := res.(page.Page[medical.ServiceCatalog])
		http.Paginated(c, &p, "Services")
	})(c)
}

// ------------------------------------------------------------
// Staff write handlers (shared for employee + manager)
// ------------------------------------------------------------

func (h *MedicalHandler) CreateSession(c *gin.Context) {
	logic := func(ctx *gin.Context, req dtos.CreateSessionRequest) (any, error) {
		cmd, err := mappers.CreateSessionRequestToCommand(req)
		if err != nil {
			return nil, err
		}
		sess, err := h.sessionSvc.CreateSession(ctx.Request.Context(), cmd)
		if err != nil {
			return nil, err
		}
		return sess.ID.Value(), nil
	}
	http.HandleCreateRequest(h.validator, "Session", logic)(c)
}

func (h *MedicalHandler) UpdateSession(c *gin.Context) {
	logic := func(ctx *gin.Context, req dtos.UpdateSessionRequest) (any, error) {
		sessionID, err := parseSessionID(ctx)
		if err != nil {
			return nil, err
		}
		cmd, err := mappers.UpdateSessionRequestToCommand(sessionID.Value(), req)
		if err != nil {
			return nil, err
		}
		_, err = h.sessionSvc.UpdateSession(ctx.Request.Context(), cmd)
		if err != nil {
			return nil, err
		}
		return nil, nil
	}
	http.HandleUpdateRequest(h.validator, "Session", logic)(c)
}

func (h *MedicalHandler) SoftDeleteSession(c *gin.Context) {
	logic := func(ctx *gin.Context) (any, error) {
		sessionID, err := parseSessionID(ctx)
		if err != nil {
			return nil, err
		}
		return nil, h.sessionSvc.SoftDeleteSession(ctx.Request.Context(), sessionID)
	}
	http.HandleDeleteRequest(h.validator, "Session", logic)(c)
}

func (h *MedicalHandler) HardDeleteSession(c *gin.Context) {
	logic := func(ctx *gin.Context) (any, error) {
		sessionID, err := parseSessionID(ctx)
		if err != nil {
			return nil, err
		}
		return nil, h.sessionSvc.HardDeleteSession(ctx.Request.Context(), sessionID)
	}
	http.HandleDeleteRequest(h.validator, "Session", logic)(c)
}

func (h *MedicalHandler) RestoreSession(c *gin.Context) {
	logic := func(ctx *gin.Context) (any, error) {
		sessionID, err := parseSessionID(ctx)
		if err != nil {
			return nil, err
		}
		return nil, h.sessionSvc.RestoreSession(ctx.Request.Context(), sessionID)
	}
	http.HandleRequestNoBodyWithResponder(h.validator, logic, func(c *gin.Context, _ any) {
		http.Updated(c, nil, "Session")
	})(c)
}

func (h *MedicalHandler) AddVaccination(c *gin.Context) {
	logic := func(ctx *gin.Context, req dtos.AddVaccinationRequest) (any, error) {
		cmd, err := mappers.AddVaccinationRequestToCommand(req)
		if err != nil {
			return nil, err
		}
		v, err := h.vaccinationSvc.AddVaccination(ctx.Request.Context(), cmd)
		if err != nil {
			return nil, err
		}
		return v.ID.Value(), nil
	}
	http.HandleCreateRequest(h.validator, "Vaccination", logic)(c)
}

func (h *MedicalHandler) UpdateVaccination(c *gin.Context) {
	logic := func(ctx *gin.Context, req dtos.UpdateVaccinationRequest) (any, error) {
		id, err := http.ParseParamToUInt(ctx, "vaccination_id")
		if err != nil {
			return nil, err
		}
		cmd, err := mappers.UpdateVaccinationRequestToCommand(id, req)
		if err != nil {
			return nil, err
		}
		_, err = h.vaccinationSvc.UpdateVaccination(ctx.Request.Context(), cmd)
		if err != nil {
			return nil, err
		}
		return nil, nil
	}
	http.HandleUpdateRequest(h.validator, "Vaccination", logic)(c)
}

func (h *MedicalHandler) RemoveVaccination(c *gin.Context) {
	logic := func(ctx *gin.Context) (any, error) {
		id, err := http.ParseParamToUInt(ctx, "vaccination_id")
		if err != nil {
			return nil, err
		}
		return nil, h.vaccinationSvc.RemoveVaccination(ctx.Request.Context(), medical.NewVaccinationID(id))
	}
	http.HandleDeleteRequest(h.validator, "Vaccination", logic)(c)
}

func (h *MedicalHandler) AddSurgery(c *gin.Context) {
	logic := func(ctx *gin.Context, req dtos.AddSurgeryRequest) (any, error) {
		cmd, err := mappers.AddSurgeryRequestToCommand(req)
		if err != nil {
			return nil, err
		}
		s, err := h.surgerySvc.AddSurgery(ctx.Request.Context(), cmd)
		if err != nil {
			return nil, err
		}
		return s.ID.Value(), nil
	}
	http.HandleCreateRequest(h.validator, "Surgery", logic)(c)
}

func (h *MedicalHandler) UpdateSurgery(c *gin.Context) {
	logic := func(ctx *gin.Context, req dtos.UpdateSurgeryRequest) (any, error) {
		id, err := http.ParseParamToUInt(ctx, "surgery_id")
		if err != nil {
			return nil, err
		}
		cmd, err := mappers.UpdateSurgeryRequestToCommand(id, req)
		if err != nil {
			return nil, err
		}
		_, err = h.surgerySvc.UpdateSurgery(ctx.Request.Context(), cmd)
		if err != nil {
			return nil, err
		}
		return nil, nil
	}
	http.HandleUpdateRequest(h.validator, "Surgery", logic)(c)
}

func (h *MedicalHandler) RemoveSurgery(c *gin.Context) {
	logic := func(ctx *gin.Context) (any, error) {
		id, err := http.ParseParamToUInt(ctx, "surgery_id")
		if err != nil {
			return nil, err
		}
		return nil, h.surgerySvc.RemoveSurgery(ctx.Request.Context(), medical.NewSurgeryID(id))
	}
	http.HandleDeleteRequest(h.validator, "Surgery", logic)(c)
}

func (h *MedicalHandler) AddPrescription(c *gin.Context) {
	logic := func(ctx *gin.Context, req dtos.AddPrescriptionRequest) (any, error) {
		cmd, err := mappers.AddPrescriptionRequestToCommand(req)
		if err != nil {
			return nil, err
		}
		p, err := h.prescriptionSvc.AddPrescription(ctx.Request.Context(), cmd)
		if err != nil {
			return nil, err
		}
		return p.ID.Value(), nil
	}
	http.HandleCreateRequest(h.validator, "Prescription", logic)(c)
}

func (h *MedicalHandler) UpdatePrescription(c *gin.Context) {
	logic := func(ctx *gin.Context, req dtos.UpdatePrescriptionRequest) (any, error) {
		id, err := http.ParseParamToUInt(ctx, "prescription_id")
		if err != nil {
			return nil, err
		}
		cmd := mappers.UpdatePrescriptionRequestToCommand(id, req)
		_, err = h.prescriptionSvc.UpdatePrescription(ctx.Request.Context(), cmd)
		if err != nil {
			return nil, err
		}
		return nil, nil
	}
	http.HandleUpdateRequest(h.validator, "Prescription", logic)(c)
}

func (h *MedicalHandler) RemovePrescription(c *gin.Context) {
	logic := func(ctx *gin.Context) (any, error) {
		id, err := http.ParseParamToUInt(ctx, "prescription_id")
		if err != nil {
			return nil, err
		}
		return nil, h.prescriptionSvc.RemovePrescription(ctx.Request.Context(), medical.NewPrescriptionID(id))
	}
	http.HandleDeleteRequest(h.validator, "Prescription", logic)(c)
}

func (h *MedicalHandler) AddAttachment(c *gin.Context) {
	logic := func(ctx *gin.Context, req dtos.AddAttachmentRequest) (any, error) {
		cmd, err := mappers.AddAttachmentRequestToCommand(req)
		if err != nil {
			return nil, err
		}
		a, err := h.attachmentSvc.AddAttachment(ctx.Request.Context(), cmd)
		if err != nil {
			return nil, err
		}
		return a.ID.Value(), nil
	}
	http.HandleCreateRequest(h.validator, "Attachment", logic)(c)
}

func (h *MedicalHandler) RemoveAttachment(c *gin.Context) {
	logic := func(ctx *gin.Context) (any, error) {
		id, err := http.ParseParamToUInt(ctx, "attachment_id")
		if err != nil {
			return nil, err
		}
		return nil, h.attachmentSvc.RemoveAttachment(ctx.Request.Context(), medical.NewAttachmentID(id))
	}
	http.HandleDeleteRequest(h.validator, "Attachment", logic)(c)
}

func (h *MedicalHandler) AddSessionService(c *gin.Context) {
	logic := func(ctx *gin.Context, req dtos.AddSessionServiceRequest) (any, error) {
		cmd := mappers.AddSessionServiceRequestToCommand(req)
		s, err := h.sessionServiceSvc.AddService(ctx.Request.Context(), cmd)
		if err != nil {
			return nil, err
		}
		return s.ID.Value(), nil
	}
	http.HandleCreateRequest(h.validator, "SessionService", logic)(c)
}

func (h *MedicalHandler) RemoveSessionService(c *gin.Context) {
	logic := func(ctx *gin.Context) (any, error) {
		id, err := http.ParseParamToUInt(ctx, "session_service_id")
		if err != nil {
			return nil, err
		}
		return nil, h.sessionServiceSvc.RemoveService(ctx.Request.Context(), medical.NewSessionServiceID(id))
	}
	http.HandleDeleteRequest(h.validator, "SessionService", logic)(c)
}

func (h *MedicalHandler) CreateVaccineCatalog(c *gin.Context) {
	logic := func(ctx *gin.Context, req dtos.CreateVaccineCatalogRequest) (any, error) {
		cmd := mappers.CreateVaccineCatalogRequestToCommand(req)
		v, err := h.vaccineCatalogSvc.CreateVaccine(ctx.Request.Context(), cmd)
		if err != nil {
			return nil, err
		}
		return v.ID.Value(), nil
	}
	http.HandleCreateRequest(h.validator, "Vaccine", logic)(c)
}

func (h *MedicalHandler) DeactivateVaccine(c *gin.Context) {
	logic := func(ctx *gin.Context) (any, error) {
		id, err := http.ParseParamToUInt(ctx, "id")
		if err != nil {
			return nil, err
		}
		return nil, h.vaccineCatalogSvc.DeactivateVaccine(ctx.Request.Context(), medical.NewVaccineCatalogID(id))
	}
	http.HandleDeleteRequest(h.validator, "Vaccine", logic)(c)
}

func (h *MedicalHandler) CreateMedicationCatalog(c *gin.Context) {
	logic := func(ctx *gin.Context, req dtos.CreateMedicationCatalogRequest) (any, error) {
		cmd := mappers.CreateMedicationCatalogRequestToCommand(req)
		m, err := h.medicationCatalogSvc.CreateMedication(ctx.Request.Context(), cmd)
		if err != nil {
			return nil, err
		}
		return m.ID.Value(), nil
	}
	http.HandleCreateRequest(h.validator, "Medication", logic)(c)
}

func (h *MedicalHandler) DeactivateMedication(c *gin.Context) {
	logic := func(ctx *gin.Context) (any, error) {
		id, err := http.ParseParamToUInt(ctx, "id")
		if err != nil {
			return nil, err
		}
		return nil, h.medicationCatalogSvc.DeactivateMedication(ctx.Request.Context(), medical.NewMedicationID(id))
	}
	http.HandleDeleteRequest(h.validator, "Medication", logic)(c)
}

func (h *MedicalHandler) CreateServiceCatalog(c *gin.Context) {
	logic := func(ctx *gin.Context, req dtos.CreateServiceCatalogRequest) (any, error) {
		cmd, err := mappers.CreateServiceCatalogRequestToCommand(req)
		if err != nil {
			return nil, err
		}
		s, err := h.serviceCatalogSvc.CreateService(ctx.Request.Context(), cmd)
		if err != nil {
			return nil, err
		}
		return s.ID.Value(), nil
	}
	http.HandleCreateRequest(h.validator, "Service", logic)(c)
}

func (h *MedicalHandler) DeactivateServiceCatalog(c *gin.Context) {
	logic := func(ctx *gin.Context) (any, error) {
		id, err := http.ParseParamToUInt(ctx, "id")
		if err != nil {
			return nil, err
		}
		return nil, h.serviceCatalogSvc.DeactivateService(ctx.Request.Context(), medical.NewServiceCatalogID(id))
	}
	http.HandleDeleteRequest(h.validator, "Service", logic)(c)
}
