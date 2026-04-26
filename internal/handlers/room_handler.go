package handlers

import (
	"net/http"
	"time"

	"lodge-system/internal/middleware"
	"lodge-system/internal/models"
	"lodge-system/internal/services"
	"lodge-system/pkg/utils"

	"github.com/google/uuid"
)

type RoomHandler struct {
	service *services.RoomService
}

func NewRoomHandler(service *services.RoomService) *RoomHandler {
	return &RoomHandler{service: service}
}

func (h *RoomHandler) List(w http.ResponseWriter, r *http.Request) {
	orgID, _ := middleware.GetOrgIDFromContext(r.Context())
	pag := utils.ParsePagination(r)
	roomType := r.URL.Query().Get("type")

	var isAvailable *bool
	if v := r.URL.Query().Get("is_available"); v != "" {
		b := v == "true"
		isAvailable = &b
	}

	rooms, total, err := h.service.List(orgID, roomType, isAvailable, pag.Page, pag.PageSize)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusOK, utils.PaginatedResponse{
		Data:     rooms,
		Page:     pag.Page,
		PageSize: pag.PageSize,
		Total:    total,
	})
}

// GuestList handles GET /api/v1/guest/rooms — public, requires ?org_id= query param
func (h *RoomHandler) GuestList(w http.ResponseWriter, r *http.Request) {
	orgIDStr := r.URL.Query().Get("org_id")
	if orgIDStr == "" {
		utils.RespondError(w, http.StatusBadRequest, "org_id query param is required")
		return
	}
	orgID, err := uuid.Parse(orgIDStr)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid org_id")
		return
	}

	pag := utils.ParsePagination(r)
	roomType := r.URL.Query().Get("type")

	var isAvailable *bool
	if v := r.URL.Query().Get("is_available"); v != "" {
		b := v == "true"
		isAvailable = &b
	}

	rooms, total, err := h.service.List(orgID, roomType, isAvailable, pag.Page, pag.PageSize)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusOK, utils.PaginatedResponse{
		Data:     rooms,
		Page:     pag.Page,
		PageSize: pag.PageSize,
		Total:    total,
	})
}

// GuestGetByID handles GET /api/v1/guest/rooms/{id} — public, no org required (unscoped lookup)
func (h *RoomHandler) GuestGetByID(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid room ID")
		return
	}

	room, err := h.service.GetByIDUnscoped(id)
	if err != nil {
		utils.RespondError(w, http.StatusNotFound, "Room not found")
		return
	}

	utils.RespondJSON(w, http.StatusOK, room)
}

func (h *RoomHandler) ListAvailable(w http.ResponseWriter, r *http.Request) {
	orgID, _ := middleware.GetOrgIDFromContext(r.Context())
	checkInStr := r.URL.Query().Get("check_in")
	checkOutStr := r.URL.Query().Get("check_out")
	roomType := r.URL.Query().Get("type")

	if checkInStr == "" || checkOutStr == "" {
		utils.RespondError(w, http.StatusBadRequest, "check_in and check_out query params are required (YYYY-MM-DD)")
		return
	}

	checkIn, err := time.Parse("2006-01-02", checkInStr)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "invalid check_in date, expected YYYY-MM-DD")
		return
	}
	checkOut, err := time.Parse("2006-01-02", checkOutStr)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "invalid check_out date, expected YYYY-MM-DD")
		return
	}

	rooms, err := h.service.ListAvailable(orgID, checkIn, checkOut, roomType)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusOK, rooms)
}

func (h *RoomHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid room ID")
		return
	}
	orgID, _ := middleware.GetOrgIDFromContext(r.Context())

	room, err := h.service.GetByID(id, orgID)
	if err != nil {
		utils.RespondError(w, http.StatusNotFound, "Room not found")
		return
	}

	utils.RespondJSON(w, http.StatusOK, room)
}

func (h *RoomHandler) Create(w http.ResponseWriter, r *http.Request) {
	orgID, _ := middleware.GetOrgIDFromContext(r.Context())
	var room models.Room
	if err := utils.DecodeJson(r, &room); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.service.Create(&room, orgID); err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusCreated, room)
}

func (h *RoomHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid room ID")
		return
	}
	orgID, _ := middleware.GetOrgIDFromContext(r.Context())

	var updates models.Room
	if err := utils.DecodeJson(r, &updates); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	room, err := h.service.Update(id, orgID, &updates)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusOK, room)
}

func (h *RoomHandler) UpdateImages(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid room ID")
		return
	}
	orgID, _ := middleware.GetOrgIDFromContext(r.Context())

	var req struct {
		Images []string `json:"images"`
	}
	if err := utils.DecodeJson(r, &req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	room, err := h.service.UpdateImages(id, orgID, req.Images)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusOK, room)
}

func (h *RoomHandler) SetAvailability(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid room ID")
		return
	}
	orgID, _ := middleware.GetOrgIDFromContext(r.Context())

	var req struct {
		IsAvailable bool `json:"is_available"`
	}
	if err := utils.DecodeJson(r, &req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.service.SetAvailability(id, orgID, req.IsAvailable); err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"id":           id,
		"is_available": req.IsAvailable,
	})
}

func (h *RoomHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid room ID")
		return
	}
	orgID, _ := middleware.GetOrgIDFromContext(r.Context())

	if err := h.service.Delete(id, orgID); err != nil {
		utils.RespondError(w, http.StatusNotFound, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]string{"message": "Room deleted successfully"})
}
