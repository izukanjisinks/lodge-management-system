package handlers

import (
	"net/http"

	"hr-system/internal/models"
	"hr-system/internal/services"
	"hr-system/pkg/utils"

	"github.com/google/uuid"
)

type EmergencyContactHandler struct {
	service *services.EmergencyContactService
}

func NewEmergencyContactHandler(service *services.EmergencyContactService) *EmergencyContactHandler {
	return &EmergencyContactHandler{service: service}
}

func (h *EmergencyContactHandler) Create(w http.ResponseWriter, r *http.Request) {
	employeeID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid employee ID")
		return
	}

	var contact models.EmergencyContact
	if err := utils.DecodeJson(r, &contact); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	contact.EmployeeID = employeeID

	if err := h.service.Create(&contact); err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	utils.RespondJSON(w, http.StatusCreated, contact)
}

func (h *EmergencyContactHandler) ListByEmployee(w http.ResponseWriter, r *http.Request) {
	employeeID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid employee ID")
		return
	}

	contacts, err := h.service.ListByEmployee(employeeID)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Failed to list emergency contacts")
		return
	}
	utils.RespondJSON(w, http.StatusOK, contacts)
}

func (h *EmergencyContactHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid contact ID")
		return
	}
	var contact models.EmergencyContact
	if err := utils.DecodeJson(r, &contact); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	contact.ID = id
	if err := h.service.Update(&contact); err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	updated, _ := h.service.GetByID(id)
	utils.RespondJSON(w, http.StatusOK, updated)
}

func (h *EmergencyContactHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid contact ID")
		return
	}
	if err := h.service.Delete(id); err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	utils.RespondJSON(w, http.StatusOK, map[string]string{"message": "Emergency contact deleted"})
}
