package handlers

import (
	"net/http"

	"hr-system/internal/middleware"
	"hr-system/internal/models"
	"hr-system/internal/services"
	"hr-system/pkg/utils"

	"github.com/google/uuid"
)

type EmployeeDocumentHandler struct {
	service *services.EmployeeDocumentService
}

func NewEmployeeDocumentHandler(service *services.EmployeeDocumentService) *EmployeeDocumentHandler {
	return &EmployeeDocumentHandler{service: service}
}

func (h *EmployeeDocumentHandler) Create(w http.ResponseWriter, r *http.Request) {
	employeeID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid employee ID")
		return
	}

	userID, _ := middleware.GetUserIDFromContext(r.Context())

	var doc models.EmployeeDocument
	if err := utils.DecodeJson(r, &doc); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	doc.EmployeeID = employeeID
	doc.UploadedBy = userID

	if err := h.service.Create(&doc); err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	utils.RespondJSON(w, http.StatusCreated, doc)
}

func (h *EmployeeDocumentHandler) ListByEmployee(w http.ResponseWriter, r *http.Request) {
	employeeID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid employee ID")
		return
	}

	docType := r.URL.Query().Get("type")
	docs, err := h.service.ListByEmployee(employeeID, docType)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Failed to list documents")
		return
	}
	utils.RespondJSON(w, http.StatusOK, docs)
}

func (h *EmployeeDocumentHandler) Verify(w http.ResponseWriter, r *http.Request) {
	docID, err := uuid.Parse(r.PathValue("did"))
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid document ID")
		return
	}

	userID, _ := middleware.GetUserIDFromContext(r.Context())
	if err := h.service.Verify(docID, userID); err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	utils.RespondJSON(w, http.StatusOK, map[string]string{"message": "Document verified"})
}

func (h *EmployeeDocumentHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid document ID")
		return
	}
	if err := h.service.SoftDelete(id); err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	utils.RespondJSON(w, http.StatusOK, map[string]string{"message": "Document deleted"})
}
