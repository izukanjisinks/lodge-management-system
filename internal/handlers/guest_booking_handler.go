package handlers

import (
	"net/http"

	"lodge-system/internal/middleware"
	"lodge-system/internal/models"
	"lodge-system/internal/services"
	"lodge-system/pkg/utils"

	"github.com/google/uuid"
)

type GuestBookingHandler struct {
	service *services.GuestBookingService
}

func NewGuestBookingHandler(service *services.GuestBookingService) *GuestBookingHandler {
	return &GuestBookingHandler{service: service}
}

func guestFromContext(r *http.Request) (*models.User, bool) {
	user, ok := r.Context().Value(middleware.UserKey).(*models.User)
	return user, ok && user != nil
}

// Create handles POST /api/v1/guest/bookings
func (h *GuestBookingHandler) Create(w http.ResponseWriter, r *http.Request) {
	user, ok := guestFromContext(r)
	if !ok {
		utils.RespondError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req models.CreateBookingRequest
	if err := utils.DecodeJson(r, &req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	booking, err := h.service.Create(user.UserID, &req)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusCreated, booking)
}

// List handles GET /api/v1/guest/bookings
func (h *GuestBookingHandler) List(w http.ResponseWriter, r *http.Request) {
	user, ok := guestFromContext(r)
	if !ok {
		utils.RespondError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	pag := utils.ParsePagination(r)
	bookings, total, err := h.service.ListForGuest(user.UserID, pag.Page, pag.PageSize)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusOK, utils.PaginatedResponse{
		Data:     bookings,
		Page:     pag.Page,
		PageSize: pag.PageSize,
		Total:    total,
	})
}

// GetByID handles GET /api/v1/guest/bookings/{id}
func (h *GuestBookingHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	user, ok := guestFromContext(r)
	if !ok {
		utils.RespondError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid booking ID")
		return
	}

	booking, err := h.service.GetByID(user.UserID, id)
	if err != nil {
		if err.Error() == "forbidden" {
			utils.RespondError(w, http.StatusForbidden, "Access denied")
			return
		}
		utils.RespondError(w, http.StatusNotFound, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusOK, booking)
}

// Cancel handles PATCH /api/v1/guest/bookings/{id}/cancel
func (h *GuestBookingHandler) Cancel(w http.ResponseWriter, r *http.Request) {
	user, ok := guestFromContext(r)
	if !ok {
		utils.RespondError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid booking ID")
		return
	}

	if err := h.service.Cancel(user.UserID, id); err != nil {
		if err.Error() == "forbidden" {
			utils.RespondError(w, http.StatusForbidden, "Access denied")
			return
		}
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]string{"message": "Booking cancelled successfully"})
}
