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

// ── Menus ─────────────────────────────────────────────────────────────────────

func (h *MenuHandler) ListMenus(w http.ResponseWriter, r *http.Request) {
	orgID, _ := middleware.GetOrgIDFromContext(r.Context())
	pag := utils.ParsePagination(r)

	menus, total, err := h.service.ListMenus(orgID, pag.Page, pag.PageSize)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondJSON(w, http.StatusOK, utils.PaginatedResponse{
		Data:     menus,
		Page:     pag.Page,
		PageSize: pag.PageSize,
		Total:    total,
	})
}

func (h *MenuHandler) GetMenu(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid menu ID")
		return
	}
	orgID, _ := middleware.GetOrgIDFromContext(r.Context())

	menu, err := h.service.GetMenuByID(id, orgID)
	if err != nil {
		utils.RespondError(w, http.StatusNotFound, err.Error())
		return
	}
	utils.RespondJSON(w, http.StatusOK, menu)
}

func (h *MenuHandler) CreateMenu(w http.ResponseWriter, r *http.Request) {
	orgID, _ := middleware.GetOrgIDFromContext(r.Context())

	var req models.CreateMenuRequest
	if err := utils.DecodeJson(r, &req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	menu, err := h.service.CreateMenu(orgID, &req)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	utils.RespondJSON(w, http.StatusCreated, menu)
}

func (h *MenuHandler) UpdateMenu(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid menu ID")
		return
	}
	orgID, _ := middleware.GetOrgIDFromContext(r.Context())

	var req models.UpdateMenuRequest
	if err := utils.DecodeJson(r, &req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	menu, err := h.service.UpdateMenu(id, orgID, &req)
	if err != nil {
		utils.RespondError(w, http.StatusNotFound, err.Error())
		return
	}
	utils.RespondJSON(w, http.StatusOK, menu)
}

func (h *MenuHandler) DeleteMenu(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid menu ID")
		return
	}
	orgID, _ := middleware.GetOrgIDFromContext(r.Context())

	if err := h.service.DeleteMenu(id, orgID); err != nil {
		utils.RespondError(w, http.StatusNotFound, err.Error())
		return
	}
	utils.RespondJSON(w, http.StatusOK, map[string]string{"message": "Menu deleted successfully"})
}

// ── Menu Items ────────────────────────────────────────────────────────────────

func (h *MenuHandler) CreateMenuItem(w http.ResponseWriter, r *http.Request) {
	menuID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid menu ID")
		return
	}
	orgID, _ := middleware.GetOrgIDFromContext(r.Context())

	var req models.CreateMenuItemRequest
	if err := utils.DecodeJson(r, &req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	item, err := h.service.CreateMenuItem(menuID, orgID, &req)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	utils.RespondJSON(w, http.StatusCreated, item)
}

func (h *MenuHandler) UpdateMenuItem(w http.ResponseWriter, r *http.Request) {
	_, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid menu ID")
		return
	}
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
	_, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid menu ID")
		return
	}
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

// ── Guest ─────────────────────────────────────────────────────────────────────

func (h *MenuHandler) GuestListMenus(w http.ResponseWriter, r *http.Request) {
	var orgID *uuid.UUID
	if orgIDStr := r.URL.Query().Get("org_id"); orgIDStr != "" {
		parsed, err := uuid.Parse(orgIDStr)
		if err != nil {
			utils.RespondError(w, http.StatusBadRequest, "Invalid org_id")
			return
		}
		orgID = &parsed
	}

	pag := utils.ParsePagination(r)
	menus, total, err := h.service.GuestListMenus(orgID, pag.Page, pag.PageSize)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondJSON(w, http.StatusOK, utils.PaginatedResponse{
		Data:     menus,
		Page:     pag.Page,
		PageSize: pag.PageSize,
		Total:    total,
	})
}
