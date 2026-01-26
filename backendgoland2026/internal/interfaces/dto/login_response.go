package dto

type LoginResponse struct {
	Email      string                 `json:"email"`
	Challenges []ChallengeDTO          `json:"challenges"`
	NextPage   string                 `json:"next_page,omitempty"`
	TotalPages int                    `json:"total_pages"`
}
