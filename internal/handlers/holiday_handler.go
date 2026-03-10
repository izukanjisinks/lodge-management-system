package handlers

import (
	"net/http"
	"strconv"

	"hr-system/internal/models"
	"hr-system/internal/services"
	"hr-system/pkg/utils"

	"github.com/google/uuid"
)

type HolidayHandler struct {
	service *services.HolidayService
}

func NewHolidayHandler(service *services.HolidayService) *HolidayHandler {
	return &HolidayHandler{service: service}
}

func (h *HolidayHandler) Create(w http.ResponseWriter, r *http.Request) {
	var holiday models.Holiday
	if err := utils.DecodeJson(r, &holiday); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	if err := h.service.Create(&holiday); err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	utils.RespondJSON(w, http.StatusCreated, holiday)
}

func (h *HolidayHandler) List(w http.ResponseWriter, r *http.Request) {
	year := 0
	if y := r.URL.Query().Get("year"); y != "" {
		if v, err := strconv.Atoi(y); err == nil {
			year = v
		}
	}
	location := r.URL.Query().Get("location")
	holidays, err := h.service.List(year, location)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Failed to list holidays")
		return
	}
	utils.RespondJSON(w, http.StatusOK, holidays)
}

func (h *HolidayHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid holiday ID")
		return
	}
	var holiday models.Holiday
	if err := utils.DecodeJson(r, &holiday); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	holiday.ID = id
	if err := h.service.Update(&holiday); err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	updated, _ := h.service.GetByID(id)
	utils.RespondJSON(w, http.StatusOK, updated)
}

func (h *HolidayHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid holiday ID")
		return
	}
	if err := h.service.Delete(id); err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	utils.RespondJSON(w, http.StatusOK, map[string]string{"message": "Holiday deleted"})
}
