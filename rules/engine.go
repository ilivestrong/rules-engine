package rules

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/ilivestrong/rules-engine/helpers"
	"github.com/ilivestrong/rules-engine/models"
	"github.com/ilivestrong/rules-engine/risk"
)

const (
	RuleIncome             = "Income"
	RuleNoOfCreditCards    = "NoOfCreditCards"
	RuleAge                = "Age"
	RulePoliticallyExposed = "PoliticallyExposed"
	RulePhone              = "PhoneLocation"
	RuleMaster             = "Master"

	minSalaryConstraint          = "minimum_salary"
	minAgeConstraint             = "min_age_allowed"
	maxCreditCardsConstraint     = "max_credit_card_count"
	isExposedConstraint          = "is_pp_exposed"
	allowedAreaCodesConstraint   = "allowed_area_codes"
	checkApprovedPhoneConstraint = "check_approved_phones"

	StatusApproved Status = "approved"
	StatusDeclined Status = "declined"
)

var allRules = []string{RuleMaster, RuleIncome, RuleAge, RuleNoOfCreditCards, RulePhone, RulePoliticallyExposed}

type (
	ApprovalRule interface {
		Execute(ctx context.Context, applicant models.Applicant) bool
	}

	RuleHandler struct {
		name string
		rule ApprovalRule
	}

	PreApprovedRule struct{}
	IncomeRule      struct {
		constraints map[string]any
	}
	AgeRule struct {
		constraints map[string]any
	}
	NoOfCreditCardsRule struct {
		constraints map[string]any
	}
	PoliticallyExposedRule struct {
		constraints map[string]any
	}
	PhoneLocationRule struct {
		constraints map[string]any
	}
	MasterRule struct {
		constraints map[string]any
		fileManager helpers.FileManager
	}

	RulesEngine struct {
		rules      map[string]RuleHandler
		masterRule *RuleHandler
	}

	Status = string
)

func (rh *RuleHandler) Handle(ctx context.Context, applicant *models.Applicant) bool {
	return rh.rule.Execute(ctx, *applicant)
}

func (ir *IncomeRule) Execute(ctx context.Context, applicant models.Applicant) bool {
	minimumSalary := 100000

	if c, ok := ir.constraints[minSalaryConstraint]; ok {
		if v, ok := c.(float64); ok {
			minimumSalary = int(v)
		}
	}
	return applicant.Income > minimumSalary
}

func (ar *AgeRule) Execute(ctx context.Context, applicant models.Applicant) bool {
	minimumAgeRequired := 18

	if c, ok := ar.constraints[minAgeConstraint]; ok {
		if v, ok := c.(float64); ok {
			minimumAgeRequired = int(v)
		}
	}

	return applicant.Age >= minimumAgeRequired
}

func (cr *NoOfCreditCardsRule) Execute(ctx context.Context, applicant models.Applicant) bool {
	maxCreditCardAllowed := 3

	if c, ok := cr.constraints[maxCreditCardsConstraint]; ok {
		if v, ok := c.(float64); ok {
			maxCreditCardAllowed = int(v)
		}
	}

	creditRisk := risk.CalculateCreditRisk(applicant.Age, applicant.NumberOfCreditCards)
	return applicant.NumberOfCreditCards <= maxCreditCardAllowed && creditRisk == "LOW"
}

func (per *PoliticallyExposedRule) Execute(ctx context.Context, applicant models.Applicant) bool {
	politicallyExposed := true

	if c, ok := per.constraints[isExposedConstraint]; ok {
		politicallyExposed = c.(bool)
	}

	return *applicant.PoliticallyExposed == politicallyExposed
}

func (plr *PhoneLocationRule) Execute(ctx context.Context, applicant models.Applicant) bool {
	areaCodes := []string{"0", "2", "5", "8"}

	if c, ok := plr.constraints[allowedAreaCodesConstraint]; ok {
		var codes []any
		if codes, ok = c.([]any); !ok {
			fmt.Println("invaid area codes from config")
			return false
		}

		strCodes := make([]string, len(codes))
		for i, v := range codes {
			strCodes[i] = v.(string)
		}
		areaCodes = strCodes
	}

	pattern := fmt.Sprintf("^[%s]", strings.Join(areaCodes, ""))
	matched, err := regexp.MatchString(pattern, applicant.PhoneNumber)
	if err != nil {
		return false
	}

	return matched
}

func (bpr *MasterRule) Execute(ctx context.Context, applicant models.Applicant) bool {
	bypassIfPhoneIsApproved := true

	if c, ok := bpr.constraints[checkApprovedPhoneConstraint]; ok {
		if v, ok := c.(bool); ok {
			bypassIfPhoneIsApproved = v
		}
	}

	if bypassIfPhoneIsApproved {
		approvedPhones, err := bpr.fileManager.ListApprovedPhones()
		if err != nil {
			fmt.Println("something went wrong with loading approved phones json, will revalidate all rules")
			return false // something went wrong with approved phones checking, let's revalidate all rules again then
		}

		if _, ok := approvedPhones[applicant.PhoneNumber]; ok {
			fmt.Println("applican't phone number is pre-approved, skipping all child rules")
			return true
		}
	}
	return false // don't bypass, execute child rules
}

func (re *RulesEngine) addRuleHandler(rule ApprovalRule, name string) {
	handler := &RuleHandler{
		rule: rule,
		name: name,
	}
	if name == RuleMaster {
		re.masterRule = &RuleHandler{
			rule: rule,
			name: name,
		}
	} else {
		re.rules[name] = *handler
	}
}

func (re *RulesEngine) Verify(ctx context.Context, applicant *models.Applicant) Status {
	if re.masterRule.Handle(ctx, applicant) {
		return StatusApproved
	}

	for _, rule := range re.rules {
		if ok := rule.Handle(ctx, applicant); !ok {
			fmt.Printf("DEBUG:: failed the rule: %s", rule.name)
			return StatusDeclined
		}
	}
	return StatusApproved
}

func createRule(ruleInfo models.RuleInfo, fileMgr helpers.FileManager) (ApprovalRule, error) {
	switch ruleInfo.Name {
	case RuleIncome:
		return &IncomeRule{
			constraints: ruleInfo.Constraints,
		}, nil
	case RuleAge:
		return &AgeRule{
			constraints: ruleInfo.Constraints,
		}, nil
	case RuleNoOfCreditCards:
		return &NoOfCreditCardsRule{
			constraints: ruleInfo.Constraints,
		}, nil
	case RulePoliticallyExposed:
		return &PoliticallyExposedRule{
			constraints: ruleInfo.Constraints,
		}, nil
	case RulePhone:
		return &PhoneLocationRule{
			constraints: ruleInfo.Constraints,
		}, nil
	case RuleMaster:
		return &MasterRule{
			constraints: ruleInfo.Constraints,
			fileManager: fileMgr,
		}, nil
	default:
		fmt.Printf("Unknown rule in config: %s, skipping...\n", ruleInfo.Name)
		return nil, errors.New("unknow rule")
	}
}

func NewRulesEngine(fileManager helpers.FileManager) (*RulesEngine, error) {
	ruleInfos, err := fileManager.LoadRulesFromConfig()
	if err != nil {
		return nil, err
	}

	if len(ruleInfos) == 0 {
		return nil, errors.New("no rules found, please check rules.json")
	}

	rulesEngine := RulesEngine{
		rules: make(map[string]RuleHandler),
	}

	for _, ruleInfo := range ruleInfos {
		rule, err := createRule(ruleInfo, fileManager)
		if err != nil {
			fmt.Println("createRule:: ", err)
			continue
		}
		rulesEngine.addRuleHandler(rule, ruleInfo.Name)
	}

	if !EngineRulesValid(&rulesEngine) {
		return nil, errors.New("missing rules, please check rules.json")
	}
	return &rulesEngine, nil
}

func EngineRulesValid(engine *RulesEngine) bool {
	for _, rule := range allRules {
		if _, ok := engine.rules[rule]; !ok && rule != RuleMaster {
			return false
		}
	}
	return engine.masterRule != nil
}
