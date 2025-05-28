package handler

type CreateLeafRequest struct {
	Note     string   `json:"note" binding:"required"`
	URL      string   `json:"url" binding:"required" uri:"url"`
	Platform string   `json:"platform"`
	Tags     []string `json:"tags"`
}

type UpdateLeafRequest struct {
	Note     string   `json:"note" binding:"required"`
	URL      string   `json:"url" binding:"required"`
	Platform string   `json:"platform"`
	Tags     []string `json:"tags"`
}
