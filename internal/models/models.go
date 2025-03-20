package models

type Fact struct {
	PeriodStart         string `json:"period_start"`
	PeriodEnd           string `json:"period_end"`
	PeriodKey           string `json:"period_key"`
	IndicatorToMOID     string `json:"indicator_to_mo_id"`
	IndicatorToMOFactID string `json:"indicator_to_mo_fact_id"`
	Value               string `json:"value"`
	FactTime            string `json:"fact_time"`
	IsPlan              string `json:"is_plan"`
	AuthuserId          string `json:"auth_user_id"`
	Comment             string `json:"comment"`
}

type CheckFact struct {
	PeriodStart     string `json:"period_start"`
	PeriodEnd       string `json:"period_end"`
	PeriodKeep      string `json:"period_key"`
	IndicatorToMOID string `json:"indicator_to_mo_id"`
}
