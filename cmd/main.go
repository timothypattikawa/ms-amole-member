package main

import (
	"os"

	"github.com/labstack/echo/v4"
	"github.com/timothypattikawa/ms-kamoro-costumer/api"
	"github.com/timothypattikawa/ms-kamoro-costumer/internal/config"
	"github.com/timothypattikawa/ms-kamoro-costumer/internal/handler"
	"github.com/timothypattikawa/ms-kamoro-costumer/internal/repository"
	"github.com/timothypattikawa/ms-kamoro-costumer/internal/service"
	"github.com/timothypattikawa/ms-kamoro-costumer/pkg/utils"
)

func main() {
	getenv := os.Getenv("ENV")
	v := config.LoadViper(getenv)

	newConfig := config.NewConfig(v)

	tokenConfig := &utils.TokenConfig{
		SymmetricKey: []byte(v.GetString("security.jwt.secret")),
		Issuer:       "MemberServiceAmole",
		AccessTTL:    v.GetDuration("security.jwt.access-time"),
		RefreshTTL:   v.GetDuration("security.jwt.refresh-time"),
	}

	pgx := newConfig.DbPostgres.GetConnectionPgx(getenv)
	memberRepository := repository.NewMemberRepository(pgx)
	memberService := service.NewMemberService(v, pgx, memberRepository, *tokenConfig)
	memberHandler := handler.NewMemberHandler(memberService)

	api.RunServer(func(e *echo.Echo) {
		handler.Handler(e, memberHandler)
	}, v)
}
