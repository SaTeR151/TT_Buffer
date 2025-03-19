package models

type Fact struct {
	PeriodStart         string `json:"period_start"`
	PeriodEnd           string `json:"period_end"`
	PeriodKeep          string `json:"period_key"`
	IndicatorToMOID     string `json:"indicator_to_mo_id"`
	IndicatorToMOFactID string `json:"indicator_to_mo_fact_id"`
	Value               string `json:"value"`
	FactTime            string `json:"fact_time"`
	IsPlan              string `json:"is_plan"`
	Comment             string `json:"comment"`
}

type CheckFact struct {
	PeriodStart     string `json:"period_start"`
	PeriodEnd       string `json:"period_end"`
	PeriodKeep      string `json:"period_key"`
	IndicatorToMOID string `json:"indicator_to_mo_id"`
}
