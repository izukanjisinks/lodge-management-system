package handlers

import (
	"net/http"

	"lodge-system/internal/middleware"
	"lodge-system/internal/models"
	"lodge-system/internal/services"
	"lodge-system/pkg/utils"

	"github.com/google/uuid"
)

type MealPlanHandler struct {
	service *services.MealPlanService
}

func NewMealPlanHandler(service *services.MealPlanService) *MealPlanHandler {
	return &MealPlanHandler{service: service}
}

func (h *MealPlanHandler) List(w http.ResponseWriter, r *http.Request) {
	orgID, _ := middleware.GetOrgIDFromContext(r.Context())
	pag := utils.ParsePagination(r)

	var isActive *bool
	if v := r.URL.Query().Get("is_active"); v != "" {
		b := v == "true"
		isActive = &b
	}

	plans, total, err := h.service.List(orgID, isActive, pag.Page, pag.PageSize)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusOK, utils.PaginatedResponse{
		Data:     plans,
		Page:     pag.Page,
		PageSize: pag.PageSize,
		Total:    total,
	})
}

func (h *MealPlanHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid meal plan ID")
		return
	}
	orgID, _ := middleware.GetOrgIDFromContext(r.Context())

	plan, err := h.service.GetByID(id, orgID)
	if err != nil {
		utils.RespondError(w, http.StatusNotFound, "Meal plan not found")
		return
	}

	utils.RespondJSON(w, http.StatusOK, plan)
}

// GuestList handles GET /api/v1/guest/meal-plans — public, org_id is an optional filter
func (h *MealPlanHandler) GuestList(w http.ResponseWriter, r *http.Request) {
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

	var isActive *bool
	if v := r.URL.Query().Get("is_active"); v != "" {
		b := v == "true"
		isActive = &b
	}

	plans, total, err := h.service.GuestList(orgID, isActive, pag.Page, pag.PageSize)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusOK, utils.PaginatedResponse{
		Data:     plans,
		Page:     pag.Page,
		PageSize: pag.PageSize,
		Total:    total,
	})
}

// GuestGetByID handles GET /api/v1/guest/meal-plans/{id} — public, no org required
func (h *MealPlanHandler) GuestGetByID(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid meal plan ID")
		return
	}

	plan, err := h.service.GetByIDUnscoped(id)
	if err != nil {
		utils.RespondError(w, http.StatusNotFound, "Meal plan not found")
		return
	}

	utils.RespondJSON(w, http.StatusOK, plan)
}

func (h *MealPlanHandler) Create(w http.ResponseWriter, r *http.Request) {
	orgID, _ := middleware.GetOrgIDFromContext(r.Context())
	var req models.CreateMealPlanRequest
	if err := utils.DecodeJson(r, &req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	plan, err := h.service.Create(orgID, &req)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusCreated, plan)
}

func (h *MealPlanHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid meal plan ID")
		return
	}
	orgID, _ := middleware.GetOrgIDFromContext(r.Context())

	var req models.UpdateMealPlanRequest
	if err := utils.DecodeJson(r, &req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	plan, err := h.service.Update(id, orgID, &req)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusOK, plan)
}

func (h *MealPlanHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid meal plan ID")
		return
	}
	orgID, _ := middleware.GetOrgIDFromContext(r.Context())

	if err := h.service.Delete(id, orgID); err != nil {
		utils.RespondError(w, http.StatusNotFound, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]string{"message": "Meal plan deleted successfully"})
}
