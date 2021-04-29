package auth

import (
	"context"
	"errors"

	"github.com/bennycio/bundle/api"
)

type authServer struct {
	api.UnimplementedAuthServiceServer
}

func (a *authServer) NewJwt(ctx context.Context, user *api.User) (*api.Jwt, error) {
	token, err := NewAuthToken(user)
	if err != nil {
		return nil, err
	}
	return &api.Jwt{Jwt: token}, err
}

func (a *authServer) Validate(ctx context.Context, jwt *api.Jwt) (*api.Empty, error) {
	err := validateToken(jwt.Jwt)
	if err != nil {
		return &api.Empty{}, err
	}
	return &api.Empty{}, nil
}

func (a *authServer) Refresh(ctx context.Context, jwt *api.Jwt) (*api.Jwt, error) {
	newToken, err := RefreshToken(jwt.Jwt)
	if err != nil {
		return nil, err
	}
	return &api.Jwt{Jwt: newToken}, nil
}

func (a *authServer) CheckScope(ctx context.Context, claims *api.Claim) (*api.Empty, error) {
	result := checkScope(claims.Jwt.Jwt, claims.Scopes...)
	if !result {
		return &api.Empty{}, errors.New("not authorized")
	}
	return &api.Empty{}, nil
}

func (a *authServer) GetUser(ctx context.Context, jwt *api.Jwt) (*api.User, error) {
	user, err := GetUserFromToken(jwt.Jwt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func NewAuthServer() *authServer {
	s := &authServer{}
	return s
}
