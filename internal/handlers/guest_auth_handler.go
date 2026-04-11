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
	userService      *services.UserService
}

func NewGuestAuthHandler(guestAuthService *services.GuestAuthService, userService *services.UserService) *GuestAuthHandler {
	return &GuestAuthHandler{guestAuthService: guestAuthService, userService: userService}
}

// Register handles public guest self-registration.
// POST /api/v1/guest/register
func (h *GuestAuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req services.GuestRegisterRequest
	if err := utils.DecodeJson(r, &req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	user, err := h.guestAuthService.Register(&req)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Issue a JWT immediately so the guest is logged in after registration
	result, err := h.userService.Login(req.Email, req.Password)
	if err != nil {
		// Registration succeeded but auto-login failed — return user without token
		utils.RespondJSON(w, http.StatusCreated, userResponse(user))
		return
	}

	utils.RespondJSON(w, http.StatusCreated, result)
}

// UpdateProfile handles PUT /api/v1/guest/me
func (h *GuestAuthHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(middleware.UserKey).(*models.User)
	if !ok || user == nil {
		utils.RespondError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req services.GuestUpdateProfileRequest
	if err := utils.DecodeJson(r, &req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	profile, err := h.guestAuthService.UpdateProfile(user.UserID, &req)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusOK, profile)
}

// Me returns the logged-in guest's profile.
// GET /api/v1/guest/me
func (h *GuestAuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(middleware.UserKey).(*models.User)
	if !ok || user == nil {
		utils.RespondError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	profile, err := h.guestAuthService.GetProfileByUserID(user.UserID)
	if err != nil {
		utils.RespondError(w, http.StatusNotFound, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"user":    userResponse(user),
		"profile": profile,
	})
}
