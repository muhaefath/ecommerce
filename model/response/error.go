package response

type Error struct {
	StatusCode int    `json:"status_code" example:"400"`
	Message    string `json:"message" example:"strconv.ParseInt: parsing \"a\": invalid syntax"`
}
