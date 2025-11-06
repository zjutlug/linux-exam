package comm

type InfoResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		Username          string `json:"username"`
		TotalScore        int    `json:"total_score"`
		LastSubmit        string `json:"last_submit"`
		CompletedProblems []struct {
			ID    int    `json:"id"`
			Score int    `json:"score"`
			Name  string `json:"name"`
		} `json:"completed_problems"`
		AllProblems []struct {
			ID    int    `json:"id"`
			Score int    `json:"score"`
			Name  string `json:"name"`
		} `json:"all_problems"`
	} `json:"data"`
}
