package model

type GetRequest struct {
	Name    string   `form:"name"`
	Age     int      `form:"age"`
	Friends []string `form:"friends"`
}
