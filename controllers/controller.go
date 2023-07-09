package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ilivestrong/rules-engine/helpers"
	"github.com/ilivestrong/rules-engine/models"
	"github.com/ilivestrong/rules-engine/rules"
)

type JSONResponse struct {
	Status string `json:"status"`
}

type CrediCardApprovalHandler struct {
	RulesEngine *rules.RulesEngine
	FileManager helpers.FileManager
	DBManager   helpers.RulesEngineRepo
}

func (handler *CrediCardApprovalHandler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	var applicant models.Applicant
	err := json.NewDecoder(req.Body).Decode(&applicant)

	if err != nil || !validateInput(&applicant) {
		resp.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(resp).Encode(JSONResponse{Status: rules.StatusDeclined})
		return
	}

	switch req.Method {
	case http.MethodPost:
		var response JSONResponse
		if status := handler.RulesEngine.Verify(req.Context(), &applicant); status != rules.StatusApproved {
			response = JSONResponse{Status: rules.StatusDeclined}
		} else {
			// if err := handler.FileManager.PersistApprovedPhone(applicant.PhoneNumber); err != nil {
			// 	fmt.Printf("failed to save approved phone")
			// }

			if err := handler.DBManager.AddApprovedPhone(context.Background(), applicant.PhoneNumber); err != nil {
				fmt.Printf("failed to save approved phone")
			}
			response = JSONResponse{Status: rules.StatusApproved}
		}
		resp.Header().Set("Content-Type", "application/json")
		resp.WriteHeader(http.StatusOK)
		json.NewEncoder(resp).Encode(response)
	default:
		resp.WriteHeader(http.StatusNotImplemented)
		fmt.Fprint(resp, "Not implemented")
	}
}

func validateInput(applicant *models.Applicant) bool {
	return applicant.PoliticallyExposed != nil
}
