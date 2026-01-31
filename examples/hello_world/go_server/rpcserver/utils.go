// THIS CODE IS GENERATED

package rpcserver

import (
	"encoding/json"
	"errors"
	"net/http"
)

type rpcError struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func writeError(w http.ResponseWriter, err error) {
	status := http.StatusInternalServerError
	errType := errorTypeCustom
	msg := "error"
	if err != nil {
		msg = err.Error()
	}

	var validationErrPtr *ValidationError
	if errors.As(err, &validationErrPtr) {
		status = http.StatusBadRequest
		errType = errorTypeValidation
		msg = validationErrPtr.Error()
	}
	var validationErr ValidationError
	if errors.As(err, &validationErr) {
		status = http.StatusBadRequest
		errType = errorTypeValidation
		msg = validationErr.Error()
	}
	var inputErrPtr *InputError
	if errors.As(err, &inputErrPtr) {
		status = http.StatusBadRequest
		errType = errorTypeInput
		msg = inputErrPtr.Error()
	}
	var inputErr InputError
	if errors.As(err, &inputErr) {
		status = http.StatusBadRequest
		errType = errorTypeInput
		msg = inputErr.Error()
	}
	var unauthorizedErrPtr *UnauthorizedError
	if errors.As(err, &unauthorizedErrPtr) {
		status = http.StatusUnauthorized
		errType = errorTypeAuth
		msg = unauthorizedErrPtr.Error()
	}
	var unauthorizedErr UnauthorizedError
	if errors.As(err, &unauthorizedErr) {
		status = http.StatusUnauthorized
		errType = errorTypeAuth
		msg = unauthorizedErr.Error()
	}
	var forbiddenErrPtr *ForbiddenError
	if errors.As(err, &forbiddenErrPtr) {
		status = http.StatusForbidden
		errType = errorTypeForbidden
		msg = forbiddenErrPtr.Error()
	}
	var forbiddenErr ForbiddenError
	if errors.As(err, &forbiddenErr) {
		status = http.StatusForbidden
		errType = errorTypeForbidden
		msg = forbiddenErr.Error()
	}
	var notImplErrPtr *NotImplementedError
	if errors.As(err, &notImplErrPtr) {
		status = http.StatusNotImplemented
		errType = errorTypeNotImpl
		msg = notImplErrPtr.Error()
	}
	var notImplErr NotImplementedError
	if errors.As(err, &notImplErr) {
		status = http.StatusNotImplemented
		errType = errorTypeNotImpl
		msg = notImplErr.Error()
	}

	writeJSON(w, status, rpcError{Type: errType, Message: msg})
}

func WriteAuthError(w http.ResponseWriter, message string) {
	writeJSON(w, http.StatusUnauthorized, rpcError{Type: errorTypeAuth, Message: message})
}

func WriteUnauthorizedError(w http.ResponseWriter, message string) {
	WriteAuthError(w, message)
}
