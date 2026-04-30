package handlers

import (
	"net/http"

	"lodge-system/internal/middleware"
	"lodge-system/internal/models"
	"lodge-system/internal/services"
	"lodge-system/pkg/utils"

	"github.com/google/uuid"
)

type MenuHandler struct {
	service *services.MenuService
}

func NewMenuHandler(service *services.MenuService) *MenuHandler {
	return &MenuHandler{service: service}
}

// ── Menu ──────────────────────────────────────────────────────────────────────

func (h *MenuHandler) GetMenu(w http.ResponseWriter, r *http.Request) {
	orgID, _ := middleware.GetOrgIDFromContext(r.Context())
	p := utils.ParsePagination(r)
	category := r.URL.Query().Get("category")

	menu, err := h.service.GetMenu(orgID, category, p.Page, p.PageSize)
	if err != nil {
		utils.RespondError(w, http.StatusNotFound, err.Error())
		return
	}
	utils.RespondJSON(w, http.StatusOK, menu)
}

func (h *MenuHandler) UpsertMenu(w http.ResponseWriter, r *http.Request) {
	orgID, _ := middleware.GetOrgIDFromContext(r.Context())
	p := utils.ParsePagination(r)
	category := r.URL.Query().Get("category")

	var req models.UpdateMenuRequest
	if err := utils.DecodeJson(r, &req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	menu, err := h.service.UpsertMenu(orgID, &req, category, p.Page, p.PageSize)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondJSON(w, http.StatusOK, menu)
}

// ── Menu Items ────────────────────────────────────────────────────────────────

func (h *MenuHandler) CreateMenuItem(w http.ResponseWriter, r *http.Request) {
	orgID, _ := middleware.GetOrgIDFromContext(r.Context())

	var req models.CreateMenuItemRequest
	if err := utils.DecodeJson(r, &req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	item, err := h.service.CreateMenuItem(orgID, &req)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	utils.RespondJSON(w, http.StatusCreated, item)
}

func (h *MenuHandler) UpdateMenuItem(w http.ResponseWriter, r *http.Request) {
	itemID, err := uuid.Parse(r.PathValue("item_id"))
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid item ID")
		return
	}
	orgID, _ := middleware.GetOrgIDFromContext(r.Context())

	var req models.UpdateMenuItemRequest
	if err := utils.DecodeJson(r, &req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	item, err := h.service.UpdateMenuItem(itemID, orgID, &req)
	if err != nil {
		utils.RespondError(w, http.StatusNotFound, err.Error())
		return
	}
	utils.RespondJSON(w, http.StatusOK, item)
}

func (h *MenuHandler) DeleteMenuItem(w http.ResponseWriter, r *http.Request) {
	itemID, err := uuid.Parse(r.PathValue("item_id"))
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid item ID")
		return
	}
	orgID, _ := middleware.GetOrgIDFromContext(r.Context())

	if err := h.service.DeleteMenuItem(itemID, orgID); err != nil {
		utils.RespondError(w, http.StatusNotFound, err.Error())
		return
	}
	utils.RespondJSON(w, http.StatusOK, map[string]string{"message": "Menu item deleted successfully"})
}

// ── Guest (public) ────────────────────────────────────────────────────────────

func (h *MenuHandler) GuestGetMenu(w http.ResponseWriter, r *http.Request) {
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
	p := utils.ParsePagination(r)
	category := r.URL.Query().Get("category")

	menu, err := h.service.GuestGetMenu(orgID, category, p.Page, p.PageSize)
	if err != nil {
		utils.RespondError(w, http.StatusNotFound, "Menu not found")
		return
	}
	utils.RespondJSON(w, http.StatusOK, menu)
}
