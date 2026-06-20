package main

import (
	"context"
	"errors"
	v1 "github.com/aritxonly/deadlinerserver/kitex_gen/deadliner/v1"
)

// DeadlinerServiceImpl implements the last service interface defined in the IDL.
type DeadlinerServiceImpl struct{}

// Register implements the DeadlinerServiceImpl interface.
func (s *DeadlinerServiceImpl) Register(ctx context.Context, req *v1.RegisterRequest) (resp *v1.RegisterResponse, err error) {
	return nil, errors.New("Register is not implemented yet")
}

// Login implements the DeadlinerServiceImpl interface.
func (s *DeadlinerServiceImpl) Login(ctx context.Context, req *v1.LoginRequest) (resp *v1.LoginResponse, err error) {
	return nil, errors.New("Login is not implemented yet")
}

// RefreshSession implements the DeadlinerServiceImpl interface.
func (s *DeadlinerServiceImpl) RefreshSession(ctx context.Context, req *v1.RefreshSessionRequest) (resp *v1.RefreshSessionResponse, err error) {
	return nil, errors.New("RefreshSession is not implemented yet")
}

// PullChanges implements the DeadlinerServiceImpl interface.
func (s *DeadlinerServiceImpl) PullChanges(ctx context.Context, req *v1.PullChangesRequest) (resp *v1.PullChangesResponse, err error) {
	return nil, errors.New("PullChanges is not implemented yet")
}

// PushChanges implements the DeadlinerServiceImpl interface.
func (s *DeadlinerServiceImpl) PushChanges(ctx context.Context, req *v1.PushChangesRequest) (resp *v1.PushChangesResponse, err error) {
	return nil, errors.New("PushChanges is not implemented yet")
}
