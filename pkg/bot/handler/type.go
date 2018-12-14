package handler

type DialogState struct {
	ResponseURL    string            `json:"responseURL"`
	SubmissionType string            `json:"submissionType"`
	Values         map[string]string `json:"values"`
}
