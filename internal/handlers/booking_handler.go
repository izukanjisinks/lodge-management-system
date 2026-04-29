package handlers

import (
	"net/http"

	"lodge-system/internal/middleware"
	"lodge-system/internal/models"
	"lodge-system/internal/services"
	"lodge-system/pkg/utils"

	"github.com/google/uuid"
)

type BookingHandler struct {
	service *services.BookingService
}

func NewBookingHandler(service *services.BookingService) *BookingHandler {
	return &BookingHandler{service: service}
}

func (h *BookingHandler) List(w http.ResponseWriter, r *http.Request) {
	orgID, _ := middleware.GetOrgIDFromContext(r.Context())
	pag := utils.ParsePagination(r)
	status := r.URL.Query().Get("status")
	clientType := r.URL.Query().Get("client_type")

	var clientID *uuid.UUID
	if v := r.URL.Query().Get("client_id"); v != "" {
		id, err := uuid.Parse(v)
		if err != nil {
			utils.RespondError(w, http.StatusBadRequest, "Invalid client_id")
			return
		}
		clientID = &id
	}

	bookings, total, err := h.service.List(orgID, status, clientType, clientID, pag.Page, pag.PageSize)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusOK, utils.PaginatedResponse{
		Data:     bookings,
		Page:     pag.Page,
		PageSize: pag.PageSize,
		Total:    total,
	})
}

func (h *BookingHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid booking ID")
		return
	}
	orgID, _ := middleware.GetOrgIDFromContext(r.Context())

	booking, err := h.service.GetByID(id, orgID)
	if err != nil {
		utils.RespondError(w, http.StatusNotFound, "Booking not found")
		return
	}

	utils.RespondJSON(w, http.StatusOK, booking)
}

func (h *BookingHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		utils.RespondError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	orgID, _ := middleware.GetOrgIDFromContext(r.Context())

	var req models.CreateBookingRequest
	if err := utils.DecodeJson(r, &req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	booking, err := h.service.Create(userID, orgID, &req)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusCreated, booking)
}

func (h *BookingHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid booking ID")
		return
	}
	orgID, _ := middleware.GetOrgIDFromContext(r.Context())

	var req models.UpdateBookingRequest
	if err := utils.DecodeJson(r, &req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	booking, err := h.service.Update(id, orgID, &req)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusOK, booking)
}

func (h *BookingHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid booking ID")
		return
	}
	orgID, _ := middleware.GetOrgIDFromContext(r.Context())

	var req models.UpdateBookingStatusRequest
	if err := utils.DecodeJson(r, &req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	booking, err := h.service.UpdateStatus(id, orgID, req.Status)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusOK, booking)
}

func (h *BookingHandler) ClearOverstayed(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid booking ID")
		return
	}
	orgID, _ := middleware.GetOrgIDFromContext(r.Context())

	if err := h.service.ClearOverstayed(id, orgID); err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]string{"message": "Overstayed flag cleared"})
}

func (h *BookingHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid booking ID")
		return
	}
	orgID, _ := middleware.GetOrgIDFromContext(r.Context())

	if err := h.service.Delete(id, orgID); err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]string{"message": "Booking deleted successfully"})
}
