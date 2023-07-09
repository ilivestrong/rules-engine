package rules

import (
	"context"
	"fmt"
	"testing"

	"github.com/ilivestrong/rules-engine/helpers"
	"github.com/ilivestrong/rules-engine/helpers/mocks"
	"github.com/ilivestrong/rules-engine/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var mockRules = []models.RuleInfo{
	{
		Name: "Master",
		Constraints: map[string]any{
			"check_approved_phones": true,
		},
	},
	{
		Name: "Income",
		Constraints: map[string]any{
			"minimum_salary": 100000,
		},
	},
	{
		Name: "NoOfCreditCards",
		Constraints: map[string]any{
			"max_credit_card_allowed": 3,
		},
	},
	{
		Name: "Age",
		Constraints: map[string]any{
			"min_age_allowed": 18,
		},
	},
	{
		Name: "PoliticallyExposed",
		Constraints: map[string]any{
			"is_pp_exposed": false,
		},
	},
	{
		Name: "PhoneLocation",
		Constraints: map[string]any{
			"allowed_area_codes": []string{
				"0",
				"2",
				"5",
				"8",
			},
		},
	},
}

func Test_NewRulesEngine(t *testing.T) {
	type args struct {
		config string
	}
	type depFields struct {
		FileManager *mocks.FileManager
	}

	tests := []struct {
		name        string
		args        args
		mocks       func(df *depFields)
		expected    *RulesEngine
		wantErr     bool
		expectedErr error
	}{
		{
			name: "config load error",
			args: args{
				config: "rules_config",
			},
			mocks: func(df *depFields) {
				df.FileManager.On("LoadRulesFromConfig").Return(nil, fmt.Errorf("failed to load config"))
			},
			wantErr:     true,
			expectedErr: fmt.Errorf("failed to load config"),
		},
		{
			name: "empty config",
			args: args{
				config: "rules_config",
			},
			mocks: func(df *depFields) {
				df.FileManager.On("LoadRulesFromConfig").Return([]models.RuleInfo{}, nil)
			},
			wantErr:     true,
			expectedErr: fmt.Errorf("no rules found, please check rules.json"),
		},
		{
			name: "missing rules in config",
			args: args{
				config: "rules_config",
			},
			mocks: func(df *depFields) {
				df.FileManager.On("LoadRulesFromConfig").Return([]models.RuleInfo{
					{
						Name:        RuleIncome,
						Constraints: map[string]any{},
					},
				}, nil)
			},
			wantErr:     true,
			expectedErr: fmt.Errorf("missing rules, please check rules.json"),
		},
		{
			name: "valid rules loaded from config",
			args: args{
				config: "rules_config",
			},
			mocks: func(df *depFields) {
				df.FileManager.On("LoadRulesFromConfig").Return(mockRules, nil)
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			df := &depFields{
				FileManager: mocks.NewFileManager(t),
			}
			tt.mocks(df)

			got, err := NewRulesEngine(df.FileManager)

			if tt.wantErr {
				if assert.True(t, err != nil) {
					assert.Equal(t, tt.expectedErr, err)
					return
				}
				assert.Fail(t, "error not returned")
				return
			}
			assert.NotNil(t, got)
		})
	}
}

func Test_RulesEngine_Verify(t *testing.T) {
	type args struct {
		config    string
		applicant *models.Applicant
	}
	type depFields struct {
		FileManager *mocks.FileManager
	}

	approvedPhoneNumber := "501-324-0507"
	oneApprovedPhone := make(helpers.ApprovedPhones)
	oneApprovedPhone[approvedPhoneNumber] = true
	PPE := false
	PPEYES := true
	fileManager := helpers.NewFileManager()
	rules, _ := fileManager.LoadRulesFromConfig()

	tests := []struct {
		name        string
		args        args
		mocks       func(df *depFields)
		expected    Status
		wantErr     bool
		expectedErr error
	}{
		{
			name: "should be `approved`, pre-approved phone",
			args: args{
				config: "test_config",
				applicant: &models.Applicant{
					Income:              100000,
					NumberOfCreditCards: 3,
					Age:                 1,
					PoliticallyExposed:  &PPE,
					JobIndustryCode:     "2-930 - Exterior Plants",
					PhoneNumber:         approvedPhoneNumber,
				},
			},
			mocks: func(df *depFields) {
				df.FileManager.On("LoadRulesFromConfig").Return(rules, nil)
				df.FileManager.On("ListApprovedPhones").Return(oneApprovedPhone, nil)
			},
			wantErr:  false,
			expected: StatusApproved,
		},
		{
			name: "should be `declined`, age rule fail",
			args: args{
				config: "test_config",
				applicant: &models.Applicant{
					Income:              100000,
					NumberOfCreditCards: 1,
					Age:                 1,
					PoliticallyExposed:  &PPE,
					JobIndustryCode:     "2-930 - Exterior Plants",
					PhoneNumber:         "502-324-0507",
				},
			},
			mocks: func(df *depFields) {
				df.FileManager.On("LoadRulesFromConfig", mock.Anything).Return(rules, nil)
				df.FileManager.On("ListApprovedPhones", mock.Anything).Return(oneApprovedPhone, nil)
			},
			wantErr:  false,
			expected: StatusDeclined,
		},
		{
			name: "should be `declined`, income rule fail",
			args: args{
				config: "test_config",
				applicant: &models.Applicant{
					Income:              90000,
					NumberOfCreditCards: 1,
					Age:                 19,
					PoliticallyExposed:  &PPE,
					JobIndustryCode:     "2-930 - Exterior Plants",
					PhoneNumber:         "502-324-0507",
				},
			},
			mocks: func(df *depFields) {
				df.FileManager.On("LoadRulesFromConfig", mock.Anything).Return(rules, nil)
				df.FileManager.On("ListApprovedPhones", mock.Anything).Return(oneApprovedPhone, nil)
			},
			wantErr:  false,
			expected: StatusDeclined,
		},
		{
			name: "should be `declined`, credit cards rule fail",
			args: args{
				config: "test_config",
				applicant: &models.Applicant{
					Income:              100000,
					NumberOfCreditCards: 4,
					Age:                 19,
					PoliticallyExposed:  &PPE,
					JobIndustryCode:     "2-930 - Exterior Plants",
					PhoneNumber:         "502-324-0507",
				},
			},
			mocks: func(df *depFields) {
				df.FileManager.On("LoadRulesFromConfig", mock.Anything).Return(rules, nil)
				df.FileManager.On("ListApprovedPhones", mock.Anything).Return(oneApprovedPhone, nil)
			},
			wantErr:  false,
			expected: StatusDeclined,
		},
		{
			name: "should be `declined`, credit cards risk > 'LOW'",
			args: args{
				config: "test_config",
				applicant: &models.Applicant{
					Income:              120000,
					NumberOfCreditCards: 1,
					Age:                 19,
					PoliticallyExposed:  &PPE,
					JobIndustryCode:     "2-930 - Exterior Plants",
					PhoneNumber:         "202-324-0507",
				},
			},
			mocks: func(df *depFields) {
				df.FileManager.On("LoadRulesFromConfig", mock.Anything).Return(rules, nil)
				df.FileManager.On("ListApprovedPhones", mock.Anything).Return(oneApprovedPhone, nil)
			},
			wantErr:  false,
			expected: StatusDeclined,
		},
		{
			name: "should be `declined`, phone location rule fail",
			args: args{
				config: "test_config",
				applicant: &models.Applicant{
					Income:              100000,
					NumberOfCreditCards: 1,
					Age:                 19,
					PoliticallyExposed:  &PPE,
					JobIndustryCode:     "2-930 - Exterior Plants",
					PhoneNumber:         "402-324-0507",
				},
			},
			mocks: func(df *depFields) {
				df.FileManager.On("LoadRulesFromConfig", mock.Anything).Return(rules, nil)
				df.FileManager.On("ListApprovedPhones", mock.Anything).Return(oneApprovedPhone, nil)
			},
			wantErr:  false,
			expected: StatusDeclined,
		},
		{
			name: "should be `declined`, politically exposed rule fail",
			args: args{
				config: "test_config",
				applicant: &models.Applicant{
					Income:              100000,
					NumberOfCreditCards: 1,
					Age:                 19,
					PoliticallyExposed:  &PPEYES,
					JobIndustryCode:     "2-930 - Exterior Plants",
					PhoneNumber:         "202-324-0507",
				},
			},
			mocks: func(df *depFields) {
				df.FileManager.On("LoadRulesFromConfig", mock.Anything).Return(rules, nil)
				df.FileManager.On("ListApprovedPhones", mock.Anything).Return(oneApprovedPhone, nil)
			},
			wantErr:  false,
			expected: StatusDeclined,
		},
		{
			name: "1should be `approved`, no pre-approved phone, all other rule pass",
			args: args{
				config: "test_config",
				applicant: &models.Applicant{
					Income:              120000,
					NumberOfCreditCards: 1,
					Age:                 23,
					PoliticallyExposed:  &PPE,
					JobIndustryCode:     "2-930 - Exterior Plants",
					PhoneNumber:         "202-324-0507",
				},
			},
			mocks: func(df *depFields) {
				df.FileManager.On("LoadRulesFromConfig", mock.Anything).Return(rules, nil)
				df.FileManager.On("ListApprovedPhones", mock.Anything).Return(oneApprovedPhone, nil)
			},
			wantErr:  false,
			expected: StatusApproved,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			df := &depFields{
				FileManager: mocks.NewFileManager(t),
			}
			tt.mocks(df)

			engine, _ := NewRulesEngine(df.FileManager)
			status := engine.Verify(context.Background(), tt.args.applicant)

			assert.Equal(t, tt.expected, status)
		})
	}
}
