package model

type PingPongRequest struct {
	Name    string   `form:"name"`
	Age     int      `form:"age"`
	Friends []string `form:"friends"`
}

type PingPongResponse struct {
	Message string `json:"message"`
}
