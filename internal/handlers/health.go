package handlers

import (
	"hr-system/pkg/utils"
	"net/http"
)

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	utils.RespondJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
