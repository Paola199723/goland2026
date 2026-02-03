package dto

type LoginResponse struct {
	Email      string                 `json:"email"`
	Challenges []ChallengeDTO          `json:"challenges"`
	NextPage   string                 `json:"next_page"`
	TotalPages int                    `json:"total_pages"`
}
