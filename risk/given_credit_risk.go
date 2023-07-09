package risk

func CalculateCreditRisk(age, numberOfCreditCard int) string {
	sum := age + numberOfCreditCard
	mod := sum % 3
	if mod == 0 {
		return "LOW"
	}
	if mod == 1 {
		return "MEDIUM"
	}
	return "HIGH"
}
