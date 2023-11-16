package dto

/* Using 3 out of 7 Loan Origination Process for the reference implementation https://allcloud.in/all-cloud-blog/7-stages-in-loan-origination-process

1. Loan Application
2. Credit Decision
3. Loan Funding

*/
type LoanApplicationInputStep struct {
	AccountNo     string `json:"account_no"`
	ApplicantName string `json:"applicant_name"`
	AadhaarNumber string `json:"aadhaar_number"`
	PanNumber     string `json:"pan_number"`
}

type LoanApplicationOutputStep struct {
	ApplicationNo string `json:"application_no"`
	AccountNo     string `json:"account_no"`
}

type CreditDecisionInputStep struct {
	ApplicationNo string `json:"application_no"`
	AadhaarNumber string `json:"aadhaar_number"`
	PanNumber     string `json:"pan_number"`
	AccountNo     string `json:"account_no"`
}

type CreditDecisionOutputStep struct {
	ApplicationNo string  `json:"application_no"`
	AccountNo     string  `json:"account_no"`
	CreditScore   float32 `json:"credit_score"`
}

type LoanFundingInputStep struct {
	ApplicationNo string `json:"application_no"`
	AccountNo     string `json:"account_no"`
}

type LoanFundingOutputStep struct {
	ApplicationNo      string  `json:"application_no"`
	AccountNo          string  `json:"account_no"`
	DisbursementAmount float64 `json:"disbursement_amount"`
}
