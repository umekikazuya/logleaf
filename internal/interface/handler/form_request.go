package handler

type CreateLeafRequest struct {
	Title    string   `json:"title" binding:"required"`
	URL      string   `json:"url" binding:"required"`
	Tags     []string `json:"tags"`
	Platform string   `json:"platform"`
}

type UpdateLeafRequest struct {
	Title    string   `json:"title" binding:"required"`
	URL      string   `json:"url" binding:"required"`
	Platform string   `json:"platform"`
	Tags     []string `json:"tags"`
}
