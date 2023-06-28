package utils

type MeRequest struct {
	BaseAPIResponse
	Data struct {
		User     SelfUser  `json:"user"`
		Projects []Project `json:"projects"`
	} `json:"data"`
}

type SelfUser struct {
	ID       string       `json:"id"`
	Username string       `json:"username"`
	Email    string       `json:"email"`
	GitHub   *interface{} `json:"github"`
	APIToken string       `json:"api_token"`
	Icon     *string      `json:"icon"`
}

type User struct {
	ID       string  `json:"id"`
	Username string  `json:"username"`
	Icon     *string `json:"icon"`
}

type Project struct {
	ID        string  `json:"id"`
	Namespace string  `json:"namespace"`
	Name      string  `json:"name"`
	Tier      int     `json:"tier"`
	Icon      *string `json:"icon"`
	Feeds     []Feed  `json:"feeds"`
	Members   []User  `json:"members"`
}

type Feed struct {
	ID          string `json:"id"`
	ProjectID   string `json:"project_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Emoji       string `json:"emoji"`
	Type        int    `json:"type"`
}

type BaseAPIResponse struct {
	Success bool `json:"success"`
}

type ErrorAPIResponse struct {
	BaseAPIResponse
	Error struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}
