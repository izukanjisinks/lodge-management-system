package handlers

import (
	"net/http"

	"lodge-system/internal/middleware"
	"lodge-system/internal/models"
	"lodge-system/internal/services"
	"lodge-system/pkg/utils"

	"github.com/google/uuid"
)

type IndividualBookingRequestHandler struct {
	service *services.IndividualBookingRequestService
}

func NewIndividualBookingRequestHandler(service *services.IndividualBookingRequestService) *IndividualBookingRequestHandler {
	return &IndividualBookingRequestHandler{service: service}
}

// ─── Web user submission ──────────────────────────────────────────────────────

// Submit handles POST /api/v1/web/bookings
func (h *IndividualBookingRequestHandler) Submit(w http.ResponseWriter, r *http.Request) {
	webUserID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		utils.RespondError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req models.SubmitIndividualBookingRequest
	if err := utils.DecodeJson(r, &req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	result, err := h.service.Submit(webUserID, &req)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusCreated, result)
}

// ListForWebUser handles GET /api/v1/web/bookings
func (h *IndividualBookingRequestHandler) ListForWebUser(w http.ResponseWriter, r *http.Request) {
	webUserID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		utils.RespondError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	p := utils.ParsePagination(r)
	requests, total, err := h.service.ListForWebUser(webUserID, p.Page, p.PageSize)
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

// GetForWebUser handles GET /api/v1/web/bookings/{id}
func (h *IndividualBookingRequestHandler) GetForWebUser(w http.ResponseWriter, r *http.Request) {
	webUserID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		utils.RespondError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request ID")
		return
	}

	req, err := h.service.GetForWebUser(id, webUserID)
	if err != nil {
		utils.RespondError(w, http.StatusNotFound, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusOK, req)
}

// CancelForWebUser handles PATCH /api/v1/web/bookings/{id}/cancel
func (h *IndividualBookingRequestHandler) CancelForWebUser(w http.ResponseWriter, r *http.Request) {
	webUserID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		utils.RespondError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request ID")
		return
	}

	if err := h.service.CancelForWebUser(id, webUserID); err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]string{"message": "Request cancelled"})
}

// ─── Backoffice ───────────────────────────────────────────────────────────────

// List handles GET /api/v1/booking-requests/individual
func (h *IndividualBookingRequestHandler) List(w http.ResponseWriter, r *http.Request) {
	orgID, _ := middleware.GetOrgIDFromContext(r.Context())
	p := utils.ParsePagination(r)
	status := r.URL.Query().Get("status")

	requests, total, err := h.service.List(orgID, status, p.Page, p.PageSize)
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

// GetByID handles GET /api/v1/booking-requests/individual/{id}
func (h *IndividualBookingRequestHandler) GetByID(w http.ResponseWriter, r *http.Request) {
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

// Approve handles PUT /api/v1/booking-requests/individual/{id}/approve
func (h *IndividualBookingRequestHandler) Approve(w http.ResponseWriter, r *http.Request) {
	orgID, _ := middleware.GetOrgIDFromContext(r.Context())
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request ID")
		return
	}

	booking, err := h.service.Approve(id, orgID)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusOK, booking)
}

// Reject handles PUT /api/v1/booking-requests/individual/{id}/reject
func (h *IndividualBookingRequestHandler) Reject(w http.ResponseWriter, r *http.Request) {
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

// Cancel handles PUT /api/v1/booking-requests/individual/{id}/cancel
func (h *IndividualBookingRequestHandler) Cancel(w http.ResponseWriter, r *http.Request) {
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
