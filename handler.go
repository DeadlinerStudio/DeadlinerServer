package main

import (
	"context"
	v1 "github.com/aritxonly/deadlinerserver/kitex_gen/deadliner/v1"
)

// DeadlinerServiceImpl implements the last service interface defined in the IDL.
type DeadlinerServiceImpl struct{}

// Register implements the DeadlinerServiceImpl interface.
func (s *DeadlinerServiceImpl) Register(ctx context.Context, req *v1.RegisterRequest) (resp *v1.RegisterResponse, err error) {
	// TODO: Your code here...
	return
}

// Login implements the DeadlinerServiceImpl interface.
func (s *DeadlinerServiceImpl) Login(ctx context.Context, req *v1.LoginRequest) (resp *v1.LoginResponse, err error) {
	// TODO: Your code here...
	return
}

// RefreshSession implements the DeadlinerServiceImpl interface.
func (s *DeadlinerServiceImpl) RefreshSession(ctx context.Context, req *v1.RefreshSessionRequest) (resp *v1.RefreshSessionResponse, err error) {
	// TODO: Your code here...
	return
}

// PullChanges implements the DeadlinerServiceImpl interface.
func (s *DeadlinerServiceImpl) PullChanges(ctx context.Context, req *v1.PullChangesRequest) (resp *v1.PullChangesResponse, err error) {
	// TODO: Your code here...
	return
}

// PushChanges implements the DeadlinerServiceImpl interface.
func (s *DeadlinerServiceImpl) PushChanges(ctx context.Context, req *v1.PushChangesRequest) (resp *v1.PushChangesResponse, err error) {
	// TODO: Your code here...
	return
}
