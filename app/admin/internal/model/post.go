package model

type PostFormReq struct {
	Name    string   `form:"name"`
	Age     int      `form:"age"`
	Friends []string `form:"friends"`
}

type PostJsonReq struct {
	Name    string   `json:"name"`
	Age     int      `json:"age"`
	Friends []string `json:"friends"`
}

type PostResponse struct {
	Success string `json:"success"`
}
