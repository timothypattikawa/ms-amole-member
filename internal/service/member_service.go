package service

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/viper"
	"github.com/timothypattikawa/ms-kamoro-costumer/internal/dto"
	"github.com/timothypattikawa/ms-kamoro-costumer/internal/repository"
	sqlc "github.com/timothypattikawa/ms-kamoro-costumer/internal/repository/postgres"
	"github.com/timothypattikawa/ms-kamoro-costumer/pkg/exception"
	"github.com/timothypattikawa/ms-kamoro-costumer/pkg/utils"
	"log"
)

type MemberService interface {
	RegistrationMember(ctx context.Context, req dto.RegistrationRequest) error
	LoginMember(ctx context.Context, req dto.LoginRequest) (dto.LoginResponse, error)
	GetMemberInfo(ctx context.Context, id int64) (dto.Member, error)
}

type MemberServiceImpl struct {
	v  *viper.Viper
	db *pgxpool.Pool
	mr repository.MemberRepository
	tc utils.TokenConfig
}

func (m MemberServiceImpl) GetMemberInfo(ctx context.Context, id int64) (dto.Member, error) {
	var member dto.Member
	err := m.mr.Exec(ctx, func(q *sqlc.Queries) error {
		memberData, err := q.GetMemberById(ctx, id)
		if err != nil {
			return exception.NewNotFoundError(fmt.Sprintf("member %d not found", id))
		}

		member = dto.Member{
			Name:    memberData.Name,
			Email:   memberData.Email,
			Address: memberData.Address,
		}
		return nil
	})
	if err != nil {
		return dto.Member{}, err
	}

	return member, nil
}

func (m MemberServiceImpl) RegistrationMember(ctx context.Context, req dto.RegistrationRequest) error {
	err := m.mr.ExecTx(ctx, func(q *sqlc.Queries) error {
		_, err := q.GetMemberByEmail(ctx, req.Email)
		if err == nil {
			return exception.NewBadReqeustError("member email already exists")
		}

		password, err := utils.GenerateHashPassword(req.Password)
		if err != nil {
			return exception.NewInternalServerError("Something went wrong!!")
		}

		err = q.InsertMember(ctx, sqlc.InsertMemberParams{
			Email:    req.Email,
			Name:     req.Name,
			Password: password,
			Address:  req.Address,
		})
		if err != nil {
			log.Println(err)
			return err
		}

		return nil
	})

	return err
}

func (m MemberServiceImpl) LoginMember(ctx context.Context, req dto.LoginRequest) (dto.LoginResponse, error) {
	var token string
	var memberData dto.Member
	err := m.mr.Exec(ctx, func(q *sqlc.Queries) error {
		member, err := q.GetMemberByEmail(ctx, req.Email)
		if err != nil {
			log.Printf("not found member by email %s, err{%v}", req.Email, err)
			return exception.NewNotFoundError("failed to login member, email not found")
		}

		password := utils.ValidatePassword(member.Password, req.Password)
		if password == false {
			return exception.NewBadReqeustError("failed to login member, password not match")
		}

		token, err = m.tc.GenerateAccessToken(fmt.Sprint(member.ID), "public")
		if err != nil {
			log.Printf("failed to generate access token: %v", err)
			return exception.NewInternalServerError("Something went wrong")
		}

		memberData = dto.Member{
			Name:    member.Name,
			Email:   member.Email,
			Address: member.Address,
		}
		return nil
	})
	if err != nil {
		return dto.LoginResponse{}, err
	}

	return dto.LoginResponse{
		Token: token, MemberData: memberData,
	}, nil
}

func NewMemberService(v *viper.Viper,
	db *pgxpool.Pool,
	mr repository.MemberRepository,
	tc utils.TokenConfig) MemberService {
	return &MemberServiceImpl{
		db: db,
		v:  v,
		mr: mr,
		tc: tc,
	}
}
