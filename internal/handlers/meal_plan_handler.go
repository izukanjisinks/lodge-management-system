package handlers

import (
	"net/http"

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
	pag := utils.ParsePagination(r)

	var isActive *bool
	if v := r.URL.Query().Get("is_active"); v != "" {
		b := v == "true"
		isActive = &b
	}

	plans, total, err := h.service.List(isActive, pag.Page, pag.PageSize)
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

	plan, err := h.service.GetByID(id)
	if err != nil {
		utils.RespondError(w, http.StatusNotFound, "Meal plan not found")
		return
	}

	utils.RespondJSON(w, http.StatusOK, plan)
}

func (h *MealPlanHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req models.CreateMealPlanRequest
	if err := utils.DecodeJson(r, &req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	plan, err := h.service.Create(&req)
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

	var req models.UpdateMealPlanRequest
	if err := utils.DecodeJson(r, &req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	plan, err := h.service.Update(id, &req)
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

	if err := h.service.Delete(id); err != nil {
		utils.RespondError(w, http.StatusNotFound, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]string{"message": "Meal plan deleted successfully"})
}
