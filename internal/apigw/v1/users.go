package v1

import (
	"encoding/json"
	"net/http"

	"gitlab.com/robotomize/gb-golang/homework/03-02-umanager/pkg/pb"
)

func newUsersHandler(usersClient usersClient) *usersHandler {
	return &usersHandler{client: usersClient}
}

type usersHandler struct {
	client usersClient
}

func (h *usersHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	resp, err := h.client.ListUsers(ctx, &pb.Empty{})
	if err != nil {
		writeErrorResponse(w, err)
		return
	}

	writeResponse(w, resp.Users)
}

func (h *usersHandler) PostUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req pb.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err := h.client.CreateUser(ctx, &req)
	if err != nil {
		writeErrorResponse(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *usersHandler) DeleteUsersId(w http.ResponseWriter, r *http.Request, id string) {
	ctx := r.Context()

	_, err := h.client.DeleteUser(ctx, &pb.DeleteUserRequest{Id: id})
	if err != nil {
		writeErrorResponse(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *usersHandler) GetUsersId(w http.ResponseWriter, r *http.Request, id string) {
	ctx := r.Context()

	resp, err := h.client.GetUser(ctx, &pb.GetUserRequest{Id: id})
	if err != nil {
		writeErrorResponse(w, err)
		return
	}

	writeResponse(w, resp)
}

func (h *usersHandler) PutUsersId(w http.ResponseWriter, r *http.Request, id string) {
	ctx := r.Context()
	var req pb.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	req.Id = id

	_, err := h.client.UpdateUser(ctx, &req)
	if err != nil {
		writeErrorResponse(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
