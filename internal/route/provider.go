package route

import (
	"cyblog/internal/domain/auth"
	"cyblog/internal/route/middlewire"
	"cyblog/pkg/repo"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	NewRegisterMiddleWire,
	NewRegisterFunc,
)

type RegisteredMiddleWire struct {
	JwtAuthMiddleWire func(optional bool) gin.HandlerFunc
}

func (r *RegisteredMiddleWire) Register() {
	middlewire.AuthMiddleWire = r.JwtAuthMiddleWire

	middlewire.IsMiddleWireRegisterFinished = true
}

func NewRegisterMiddleWire(
	userRepo *repo.UserRepo,
	jwtBiz *auth.JwtBiz,
) RegisteredMiddleWire {
	return RegisteredMiddleWire{
		JwtAuthMiddleWire: middlewire.JwtAuthMiddleWire(jwtBiz, userRepo),
	}
}
