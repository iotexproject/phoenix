package models

type PutObject struct {
	Name string `json:"name"`
}

// QueryObject query object struct
type QueryObject struct {
	Pod string `json:"pod"`
	Pea string `json:"pea"`
}
