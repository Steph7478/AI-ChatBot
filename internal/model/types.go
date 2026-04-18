package model

type ResponseResult struct {
	Text       string
	Type       ResponseType
	Confidence float64
}

type ResponseType int

const (
	ResponseGenerated ResponseType = iota
	ResponseFallback
)
