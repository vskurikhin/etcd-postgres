package dto

import (
	"fmt"
	"github.com/google/uuid"
)

type Status struct {
	Status string
}

type StatusMessage struct {
	Status  string
	Message string
}

type StatusMessageRequestID struct {
	Status    string
	Message   string
	RequestID uuid.UUID
}

type StatusMessageStringRequestID struct {
	Status    string
	Message   string
	RequestID string
}

type StatusRequestID struct {
	Status    string
	RequestID uuid.UUID
}

type StatusResultRequestID struct {
	Status    string
	RequestID uuid.UUID
	Result    any
}

func StatusMessagePathDoesNotExists(path string) StatusMessage {
	return StatusMessage{
		Status:  "fail",
		Message: fmt.Sprintf("Path: %v does not exists on this server", path),
	}
}

func StatusMessageInvalidRequestID(requestID string) StatusMessageStringRequestID {
	return StatusMessageStringRequestID{
		Status:    "fail",
		Message:   "Invalid Request-Id",
		RequestID: requestID,
	}
}
