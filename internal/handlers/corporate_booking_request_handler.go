package handlers

import (
	"net/http"

	"lodge-system/internal/middleware"
	"lodge-system/internal/models"
	"lodge-system/internal/services"
	"lodge-system/pkg/utils"

	"github.com/google/uuid"
)

type CorporateBookingRequestHandler struct {
	service *services.CorporateBookingRequestService
}

func NewCorporateBookingRequestHandler(service *services.CorporateBookingRequestService) *CorporateBookingRequestHandler {
	return &CorporateBookingRequestHandler{service: service}
}

// ─── Guest submission ─────────────────────────────────────────────────────────

// SubmitAccommodation handles POST /api/v1/guest/bookings/corporate-event
// The org_id is expected in the request body (from frontend), not as a query param.
func (h *CorporateBookingRequestHandler) SubmitAccommodation(w http.ResponseWriter, r *http.Request) {
	var req models.SubmitAccommodationRequest
	if err := utils.DecodeJson(r, &req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate org_id from body
	if req.OrgID == uuid.Nil {
		utils.RespondError(w, http.StatusBadRequest, "org_id is required")
		return
	}

	result, err := h.service.SubmitAccommodation(req.OrgID, &req)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	utils.RespondJSON(w, http.StatusCreated, result)
}

// ─── Backoffice ───────────────────────────────────────────────────────────────

// List handles GET /api/v1/bookings/requests
func (h *CorporateBookingRequestHandler) List(w http.ResponseWriter, r *http.Request) {
	orgID, _ := middleware.GetOrgIDFromContext(r.Context())
	p := utils.ParsePagination(r)
	bookingType := r.URL.Query().Get("type")
	status := r.URL.Query().Get("status")

	requests, total, err := h.service.List(orgID, bookingType, status, p.Page, p.PageSize)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusOK, utils.PaginatedResponse{
		Data:     requests,
		Page:     p.Page,
		PageSize: p.PageSize,
		Total:    total,
	})
}

// GetByID handles GET /api/v1/bookings/requests/{id}
func (h *CorporateBookingRequestHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	orgID, _ := middleware.GetOrgIDFromContext(r.Context())
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request ID")
		return
	}

	req, err := h.service.GetByID(id, orgID)
	if err != nil {
		utils.RespondError(w, http.StatusNotFound, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusOK, req)
}

// Approve handles PUT /api/v1/bookings/requests/{id}/approve
func (h *CorporateBookingRequestHandler) Approve(w http.ResponseWriter, r *http.Request) {
	orgID, _ := middleware.GetOrgIDFromContext(r.Context())
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request ID")
		return
	}

	if err := h.service.Approve(id, orgID); err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]string{"message": "Request approved"})
}

// Reject handles PUT /api/v1/bookings/requests/{id}/reject
func (h *CorporateBookingRequestHandler) Reject(w http.ResponseWriter, r *http.Request) {
	orgID, _ := middleware.GetOrgIDFromContext(r.Context())
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request ID")
		return
	}

	if err := h.service.Reject(id, orgID); err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]string{"message": "Request rejected"})
}

// Cancel handles PUT /api/v1/bookings/requests/{id}/cancel
func (h *CorporateBookingRequestHandler) Cancel(w http.ResponseWriter, r *http.Request) {
	orgID, _ := middleware.GetOrgIDFromContext(r.Context())
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request ID")
		return
	}

	if err := h.service.Cancel(id, orgID); err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]string{"message": "Request cancelled"})
}
