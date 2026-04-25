package handlers

import (
	"net/http"

	"lodge-system/internal/middleware"
	"lodge-system/internal/models"
	"lodge-system/internal/services"
	"lodge-system/pkg/utils"
)

type GuestAuthHandler struct {
	guestAuthService *services.GuestAuthService
}

func NewGuestAuthHandler(guestAuthService *services.GuestAuthService) *GuestAuthHandler {
	return &GuestAuthHandler{guestAuthService: guestAuthService}
}

// Register handles POST /api/v1/guest/auth/register
func (h *GuestAuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req models.GuestRegisterRequest
	if err := utils.DecodeJson(r, &req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	guest, err := h.guestAuthService.Register(&req)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	token, err := utils.GenerateGuestToken(guest.Email, guest.ID)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	utils.RespondJSON(w, http.StatusCreated, map[string]interface{}{
		"guest": guest,
		"token": token,
	})
}

// Login handles POST /api/v1/guest/auth/login
func (h *GuestAuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := utils.DecodeJson(r, &req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	guest, token, err := h.guestAuthService.Login(req.Email, req.Password)
	if err != nil {
		utils.RespondError(w, http.StatusUnauthorized, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"guest": guest,
		"token": token,
	})
}

// Me handles GET /api/v1/guest/me
func (h *GuestAuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	guestID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		utils.RespondError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	guest, err := h.guestAuthService.GetByID(guestID)
	if err != nil {
		utils.RespondError(w, http.StatusNotFound, "Guest not found")
		return
	}

	utils.RespondJSON(w, http.StatusOK, guest)
}

// UpdateProfile handles PUT /api/v1/guest/me
func (h *GuestAuthHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	guestID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		utils.RespondError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req models.GuestUpdateRequest
	if err := utils.DecodeJson(r, &req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	guest, err := h.guestAuthService.UpdateProfile(guestID, &req)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusOK, guest)
}
