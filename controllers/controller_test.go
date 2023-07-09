package controllers

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ilivestrong/rules-engine/helpers"
	"github.com/ilivestrong/rules-engine/models"
	"github.com/ilivestrong/rules-engine/rules"
	"github.com/stretchr/testify/assert"
)

func Test_Process_Handler(t *testing.T) {
	type args struct {
		applicant *models.Applicant
	}

	PPE := false
	PPEYes := true

	tests := []struct {
		name      string
		args      args
		expectErr bool
		expected  JSONResponse
	}{
		{
			name: "should be `declined`, invalid request",
			args: args{
				applicant: nil,
			},
			expected: JSONResponse{
				Status: rules.StatusDeclined,
			},
			expectErr: true,
		},
		{
			name: "should be `approved`, all rules pass",
			args: args{
				applicant: &models.Applicant{
					Income:              120000,
					NumberOfCreditCards: 1,
					Age:                 29,
					PoliticallyExposed:  &PPE,
					JobIndustryCode:     "15-100 - Plumbing",
					PhoneNumber:         "268-741-8863",
				},
			},
			expected: JSONResponse{
				Status: rules.StatusApproved,
			},
		},
		{
			name: "should be `declined`, age rule fail",
			args: args{
				applicant: &models.Applicant{
					Income:              120000,
					NumberOfCreditCards: 1,
					Age:                 10,
					PoliticallyExposed:  &PPE,
					JobIndustryCode:     "15-100 - Plumbing",
					PhoneNumber:         "269-741-8863",
				},
			},
			expected: JSONResponse{
				Status: rules.StatusDeclined,
			},
		},
		{
			name: "should be `declined`, income rule fail",
			args: args{
				applicant: &models.Applicant{
					Income:              90000,
					NumberOfCreditCards: 1,
					Age:                 29,
					PoliticallyExposed:  &PPE,
					JobIndustryCode:     "15-100 - Plumbing",
					PhoneNumber:         "269-741-8863",
				},
			},
			expected: JSONResponse{
				Status: rules.StatusDeclined,
			},
		},
		{
			name: "should be `declined`, politically exposed rule fail",
			args: args{
				applicant: &models.Applicant{
					Income:              120000,
					NumberOfCreditCards: 1,
					Age:                 29,
					PoliticallyExposed:  &PPEYes,
					JobIndustryCode:     "15-100 - Plumbing",
					PhoneNumber:         "269-741-8863",
				},
			},
			expected: JSONResponse{
				Status: rules.StatusDeclined,
			},
		},
		{
			name: "should be `declined`, phone area code rule fail",
			args: args{
				applicant: &models.Applicant{
					Income:              120000,
					NumberOfCreditCards: 1,
					Age:                 29,
					PoliticallyExposed:  &PPE,
					JobIndustryCode:     "15-100 - Plumbing",
					PhoneNumber:         "369-741-8863",
				},
			},
			expected: JSONResponse{
				Status: rules.StatusDeclined,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fileManager := helpers.NewFileManager()
			rulesEngine, _ := rules.NewRulesEngine(fileManager)
			handler := &CrediCardApprovalHandler{
				RulesEngine: rulesEngine,
				FileManager: fileManager,
			}

			reqBody, _ := json.Marshal(tt.args.applicant)
			req, err := http.NewRequest(http.MethodPost, "/processs", bytes.NewBuffer(reqBody))
			if err != nil {
				t.Fatal(err)
			}
			rr := httptest.NewRecorder()

			handler.ServeHTTP(rr, req)

			if tt.expectErr {
				if status := rr.Result().StatusCode; status != http.StatusBadRequest {
					t.Errorf("Expected status code %d, but got %d", http.StatusBadRequest, status)
				}
			} else {
				if status := rr.Code; status != http.StatusOK {
					t.Errorf("handler returned wrong status code: got %v want %v",
						status, http.StatusOK)
				}
			}

			resp, _ := ioutil.ReadAll(rr.Body)
			var got JSONResponse
			json.Unmarshal(resp, &got)

			assert.Equal(t, tt.expected, got)
		})
	}
}
