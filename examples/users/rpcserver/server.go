package rpcserver

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type UserModel struct {
	Id       int     `json:"id"`
	Username string  `json:"username"`
	Name     string  `json:"name"`
	Surname  *string `json:"surname"`
}
type GroupModel struct {
	Name  string      `json:"name"`
	Users []UserModel `json:"users"`
}

type GetUserParams struct {
	UserId int `json:"user_id"`
}

type GetUserResult struct {
	User UserModel `json:"user"`
}

type ListUsersParams struct {
}

type ListUsersResult struct {
	Result []UserModel `json:"result"`
}

type CreateUserParams struct {
	Name    string  `json:"name"`
	Surname *string `json:"surname"`
}

type CreateUserResult struct {
	User UserModel `json:"user"`
}

type GetUsernameMapParams struct {
}

type GetUsernameMapResult struct {
	Result map[string]UserModel `json:"result"`
}

type FindGroupByNameParams struct {
	Name string `json:"name"`
}

type FindGroupByNameResult struct {
	Group GroupModel `json:"group"`
}

type RPCHandler interface {
	GetUser(GetUserParams) (GetUserResult, error)
	ListUsers(ListUsersParams) (ListUsersResult, error)
	CreateUser(CreateUserParams) (CreateUserResult, error)
	GetUsernameMap(GetUsernameMapParams) (GetUsernameMapResult, error)
	FindGroupByName(FindGroupByNameParams) (FindGroupByNameResult, error)
}

func CreateHTTPHandler(rpc RPCHandler) http.Handler {
	mux := http.NewServeMux()
	mux.Handle("POST /rpc/get_user", CreateGetUserHandler(rpc))
	mux.Handle("POST /rpc/list_users", CreateListUsersHandler(rpc))
	mux.Handle("POST /rpc/create_user", CreateCreateUserHandler(rpc))
	mux.Handle("POST /rpc/get_username_map", CreateGetUsernameMapHandler(rpc))
	mux.Handle("POST /rpc/find_group_by_name", CreateFindGroupByNameHandler(rpc))
	return mux
}

func CreateGetUserHandler(rpc RPCHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var params GetUserParams
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()
		if err := decoder.Decode(&params); err != nil && err != io.EOF {
			writeError(w, InputError{Message: err.Error()})
			return
		}
		res, err := rpc.GetUser(params)
		if err != nil {
			writeError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, res)
	})
}

func CreateListUsersHandler(rpc RPCHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var params ListUsersParams
		res, err := rpc.ListUsers(params)
		if err != nil {
			writeError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, res)
	})
}

func CreateCreateUserHandler(rpc RPCHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var params CreateUserParams
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()
		if err := decoder.Decode(&params); err != nil && err != io.EOF {
			writeError(w, InputError{Message: err.Error()})
			return
		}
		res, err := rpc.CreateUser(params)
		if err != nil {
			writeError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, res)
	})
}

func CreateGetUsernameMapHandler(rpc RPCHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var params GetUsernameMapParams
		res, err := rpc.GetUsernameMap(params)
		if err != nil {
			writeError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, res)
	})
}

func CreateFindGroupByNameHandler(rpc RPCHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var params FindGroupByNameParams
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()
		if err := decoder.Decode(&params); err != nil && err != io.EOF {
			writeError(w, InputError{Message: err.Error()})
			return
		}
		res, err := rpc.FindGroupByName(params)
		if err != nil {
			writeError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, res)
	})
}

type rpcError struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

type ValidationError struct {
	Message string
}

func (e ValidationError) Error() string {
	return e.Message
}

type InputError struct {
	Message string
}

func (e InputError) Error() string {
	return e.Message
}

type UnauthorizedError struct {
	Message string
}

func (e UnauthorizedError) Error() string {
	return e.Message
}

type ForbiddenError struct {
	Message string
}

func (e ForbiddenError) Error() string {
	return e.Message
}

type NotImplementedError struct {
	Message string
}

func (e NotImplementedError) Error() string {
	return e.Message
}

const (
	errorTypeCustom     = "custom"
	errorTypeValidation = "validation"
	errorTypeInput      = "input"
	errorTypeAuth       = "unauthorized"
	errorTypeForbidden  = "forbidden"
	errorTypeNotImpl    = "not_implemented"
)

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

	var validationErr ValidationError
	if errors.As(err, &validationErr) {
		status = http.StatusBadRequest
		errType = errorTypeValidation
		msg = validationErr.Error()
	}
	var inputErr InputError
	if errors.As(err, &inputErr) {
		status = http.StatusBadRequest
		errType = errorTypeInput
		msg = inputErr.Error()
	}
	var unauthorizedErr UnauthorizedError
	if errors.As(err, &unauthorizedErr) {
		status = http.StatusUnauthorized
		errType = errorTypeAuth
		msg = unauthorizedErr.Error()
	}
	var forbiddenErr ForbiddenError
	if errors.As(err, &forbiddenErr) {
		status = http.StatusForbidden
		errType = errorTypeForbidden
		msg = forbiddenErr.Error()
	}
	var notImplErr NotImplementedError
	if errors.As(err, &notImplErr) {
		status = http.StatusNotImplemented
		errType = errorTypeNotImpl
		msg = notImplErr.Error()
	}

	writeJSON(w, status, rpcError{Type: errType, Message: msg})
}
