package helpers

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"

	"github.com/ilivestrong/rules-engine/models"
)

type (
	FileManager interface {
		LoadRulesFromConfig() ([]models.RuleInfo, error)
		ListApprovedPhones() (ApprovedPhones, error)
		PersistApprovedPhone(phone string) error
	}
	defaultFileManager struct{}
	ApprovedPhones     = map[string]bool
)

func (dfm *defaultFileManager) LoadRulesFromConfig() ([]models.RuleInfo, error) {
	rulesConfig, _ := getJSONPaths()
	data, err := ioutil.ReadFile(rulesConfig)
	if err != nil {
		return nil, err
	}

	var ruleInfos []models.RuleInfo
	if err := json.Unmarshal(data, &ruleInfos); err != nil {
		return nil, err
	}
	return ruleInfos, nil
}

func (dfm *defaultFileManager) ListApprovedPhones() (ApprovedPhones, error) {
	_, approvedPhonesListStore := getJSONPaths()
	fileData, err := ioutil.ReadFile(approvedPhonesListStore)
	if err != nil {
		log.Fatal(err)
		return nil, errors.New("failed to load approved list phones")
	}

	var phoneNumbers ApprovedPhones
	err = json.Unmarshal(fileData, &phoneNumbers)
	if err != nil {
		log.Fatal(err)
		return nil, errors.New("invalid approved phone list")
	}
	return phoneNumbers, nil
}

func (dfm *defaultFileManager) PersistApprovedPhone(phone string) error {
	_, approvedPhonesListStore := getJSONPaths()
	phoneNumbers, err := dfm.ListApprovedPhones()
	if err != nil {
		log.Fatal(err)
		return err
	}

	phoneNumbers[phone] = true

	jsonData, err := json.Marshal(phoneNumbers)
	if err != nil {
		log.Fatal(err)
		return errors.New("failed to persist approved phone number")
	}

	err = ioutil.WriteFile(approvedPhonesListStore, jsonData, 0644)
	if err != nil {
		log.Fatal(err)
		return errors.New("failed to persist approved phone number")
	}
	return nil
}

func getJSONPaths() (rulesConfig string, approvedPhonesList string) {
	rulesConfig = "rules/rules.json"
	approvedPhonesList = "rules/approved-phone-list.json"

	// wd, _ := os.Getwd()
	// fmt.Println(wd)
	// if !strings.HasSuffix(wd, "tech-assessment-backend-engineer") ||
	// 	!strings.HasSuffix(wd, "rules-engine") {
	// 	rulesConfig = "../" + rulesConfig
	// 	approvedPhonesList = "../" + approvedPhonesList
	// }
	return
}

func NewFileManager() *defaultFileManager {
	return &defaultFileManager{}
}
