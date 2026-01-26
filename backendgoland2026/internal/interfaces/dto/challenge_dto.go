package dto

type ChallengeDTO struct {
	Ticker     string `json:"ticker"`
	TargetFrom string `json:"target_from"`
	TargetTo   string `json:"target_to"`
	Company    string `json:"company"`
	Action     string `json:"action"`
	Brokerage  string `json:"brokerage"`
	RatingFrom string `json:"rating_from"`
	RatingTo   string `json:"rating_to"`
	Time       string `json:"time"`
}

type ChallengeAPIResponse struct {
	Items    []ChallengeDTO `json:"items"`
	NextPage string         `json:"next_page,omitempty"`
}

