package auth

import (
	"context"
	"net/http"

	"github.com/pollenjp/gameserver-go/api/entity"
	"github.com/pollenjp/gameserver-go/api/service"
)

//go:generate go run github.com/matryer/moq -out auth_moq_test.go . AuthRepository
type AuthRepository interface {
	GetUserFromToken(ctx context.Context, db service.Queryer, userToken entity.UserTokenType) (*entity.User, error)
}

func NewAuthorizer(db service.Queryer, repo AuthRepository) *Authorizer {
	return &Authorizer{
		DB:   db,
		Repo: repo,
	}
}

type Authorizer struct {
	DB   service.Queryer
	Repo AuthRepository
}

// *http.Request型から認証情報を context に書き込む
func (au *Authorizer) FillContext(r *http.Request) (*http.Request, error) {
	token, err := ExtractBearerToken(r)
	if err != nil {
		return nil, err
	}

	u, err := au.Repo.GetUserFromToken(r.Context(), au.DB, token)
	if err != nil {
		return nil, err
	}

	ctx := service.SetUserId(r.Context(), u.Id)

	clone := r.Clone(ctx)
	return clone, nil
}
