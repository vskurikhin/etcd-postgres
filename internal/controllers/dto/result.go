package dto

type KeyValue struct {
	Key   string `json:"key" validate:"required"`
	Value string `json:"value" validate:"required"`
}

type Result struct {
	Value string `json:"value" validate:"required"`
}
