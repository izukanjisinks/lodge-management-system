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

// Submit handles POST /api/v1/web/bookings/corporate?type=accommodation|meals|conference|event
func (h *CorporateBookingRequestHandler) Submit(w http.ResponseWriter, r *http.Request) {
	orgIDStr := r.URL.Query().Get("org_id")
	if orgIDStr == "" {
		utils.RespondError(w, http.StatusBadRequest, "org_id is required")
		return
	}
	orgID, err := uuid.Parse(orgIDStr)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid org_id")
		return
	}

	bookingType := r.URL.Query().Get("type")

	switch bookingType {
	case models.CorporateBookingTypeAccommodation, "":
		var req models.SubmitAccommodationRequest
		if err := utils.DecodeJson(r, &req); err != nil {
			utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
			return
		}
		req.OrgID = orgID
		result, err := h.service.SubmitAccommodation(orgID, &req)
		if err != nil {
			utils.RespondError(w, http.StatusBadRequest, err.Error())
			return
		}
		utils.RespondJSON(w, http.StatusCreated, result)

	case models.CorporateBookingTypeMeals:
		var req models.SubmitMealsRequest
		if err := utils.DecodeJson(r, &req); err != nil {
			utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
			return
		}
		req.OrgID = orgID
		result, err := h.service.SubmitMeals(orgID, &req)
		if err != nil {
			utils.RespondError(w, http.StatusBadRequest, err.Error())
			return
		}
		utils.RespondJSON(w, http.StatusCreated, result)

	case models.CorporateBookingTypeConference:
		var req models.SubmitConferenceRequest
		if err := utils.DecodeJson(r, &req); err != nil {
			utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
			return
		}
		req.OrgID = orgID
		result, err := h.service.SubmitConference(orgID, &req)
		if err != nil {
			utils.RespondError(w, http.StatusBadRequest, err.Error())
			return
		}
		utils.RespondJSON(w, http.StatusCreated, result)

	case models.CorporateBookingTypeEvent:
		var req models.SubmitEventRequest
		if err := utils.DecodeJson(r, &req); err != nil {
			utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
			return
		}
		req.OrgID = orgID
		result, err := h.service.SubmitEvent(orgID, &req)
		if err != nil {
			utils.RespondError(w, http.StatusBadRequest, err.Error())
			return
		}
		utils.RespondJSON(w, http.StatusCreated, result)

	default:
		utils.RespondError(w, http.StatusBadRequest, "type must be accommodation, meals, conference, or event")
	}
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
