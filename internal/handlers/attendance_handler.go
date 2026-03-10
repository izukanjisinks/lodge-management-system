package handlers

import (
	"net/http"
	"strconv"
	"time"

	"hr-system/internal/interfaces"
	"hr-system/internal/middleware"
	"hr-system/internal/models"
	"hr-system/internal/services"
	"hr-system/pkg/utils"

	"github.com/google/uuid"
)

type AttendanceHandler struct {
	service    *services.AttendanceService
	empService *services.EmployeeService
}

func NewAttendanceHandler(svc *services.AttendanceService, empSvc *services.EmployeeService) *AttendanceHandler {
	return &AttendanceHandler{service: svc, empService: empSvc}
}

func (h *AttendanceHandler) ClockIn(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.GetUserIDFromContext(r.Context())
	emp := h.getEmployeeByUser(userID)
	if emp == nil {
		utils.RespondError(w, http.StatusBadRequest, "No employee record linked to your account")
		return
	}
	var body struct{ Notes string `json:"notes"` }
	_ = utils.DecodeJson(r, &body)
	a, err := h.service.ClockIn(emp.ID, body.Notes)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	utils.RespondJSON(w, http.StatusOK, a)
}

func (h *AttendanceHandler) ClockOut(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.GetUserIDFromContext(r.Context())
	emp := h.getEmployeeByUser(userID)
	if emp == nil {
		utils.RespondError(w, http.StatusBadRequest, "No employee record linked to your account")
		return
	}
	var body struct{ Notes string `json:"notes"` }
	_ = utils.DecodeJson(r, &body)
	a, err := h.service.ClockOut(emp.ID, body.Notes)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	utils.RespondJSON(w, http.StatusOK, a)
}

func (h *AttendanceHandler) GetMyAttendance(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.GetUserIDFromContext(r.Context())
	emp := h.getEmployeeByUser(userID)
	if emp == nil {
		utils.RespondJSON(w, http.StatusOK, []interface{}{})
		return
	}
	from, to := parseDateRange(r)
	records, err := h.service.ListByEmployee(emp.ID, from, to)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Failed to get attendance")
		return
	}
	utils.RespondJSON(w, http.StatusOK, records)
}

func (h *AttendanceHandler) GetByEmployee(w http.ResponseWriter, r *http.Request) {
	employeeID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid employee ID")
		return
	}
	from, to := parseDateRange(r)
	records, err := h.service.ListByEmployee(employeeID, from, to)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Failed to get attendance")
		return
	}
	utils.RespondJSON(w, http.StatusOK, records)
}

func (h *AttendanceHandler) GetByDepartment(w http.ResponseWriter, r *http.Request) {
	departmentID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid department ID")
		return
	}
	dateStr := r.URL.Query().Get("date")
	var date time.Time
	if dateStr != "" {
		date, _ = time.Parse("2006-01-02", dateStr)
	} else {
		date = time.Now()
	}
	records, err := h.service.ListByDepartmentAndDate(departmentID, date)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Failed to get attendance")
		return
	}
	utils.RespondJSON(w, http.StatusOK, records)
}

func (h *AttendanceHandler) CreateManual(w http.ResponseWriter, r *http.Request) {
	var a models.Attendance
	if err := utils.DecodeJson(r, &a); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	if err := h.service.CreateManual(&a); err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	utils.RespondJSON(w, http.StatusCreated, a)
}

func (h *AttendanceHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid attendance ID")
		return
	}
	var a models.Attendance
	if err := utils.DecodeJson(r, &a); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	a.ID = id
	if err := h.service.Update(&a); err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	utils.RespondJSON(w, http.StatusOK, a)
}

func (h *AttendanceHandler) GetSummary(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	employeeID, err := uuid.Parse(q.Get("employee_id"))
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "employee_id is required")
		return
	}
	month := int(time.Now().Month())
	year := time.Now().Year()
	if m := q.Get("month"); m != "" {
		if v, e := strconv.Atoi(m); e == nil {
			month = v
		}
	}
	if y := q.Get("year"); y != "" {
		if v, e := strconv.Atoi(y); e == nil {
			year = v
		}
	}
	summary, err := h.service.GetMonthlySummary(employeeID, month, year)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Failed to get summary")
		return
	}
	utils.RespondJSON(w, http.StatusOK, summary)
}

func parseDateRange(r *http.Request) (time.Time, time.Time) {
	from := time.Now().AddDate(0, -1, 0)
	to := time.Now()
	if f := r.URL.Query().Get("from"); f != "" {
		if t, err := time.Parse("2006-01-02", f); err == nil {
			from = t
		}
	}
	if t := r.URL.Query().Get("to"); t != "" {
		if parsed, err := time.Parse("2006-01-02", t); err == nil {
			to = parsed
		}
	}
	return from, to
}

func (h *AttendanceHandler) getEmployeeByUser(userID uuid.UUID) *models.Employee {
	emps, _, err := h.empService.List(interfaces.EmployeeFilter{}, 1, 500)
	if err != nil {
		return nil
	}
	for i := range emps {
		e := &emps[i]
		if e.UserID != nil && *e.UserID == userID {
			return e
		}
	}
	return nil
}
