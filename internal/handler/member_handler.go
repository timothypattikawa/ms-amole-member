package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/timothypattikawa/ms-kamoro-costumer/internal/dto"
	"github.com/timothypattikawa/ms-kamoro-costumer/internal/service"
	"github.com/timothypattikawa/ms-kamoro-costumer/pkg/exception"
	"log"
	"net/http"
	"strconv"
)

type MemberHandler struct {
	ms service.MemberService
}

func NewMemberHandler(ms service.MemberService) *MemberHandler {
	return &MemberHandler{ms: ms}
}

func (h *MemberHandler) CreateMember(c echo.Context) error {
	var registerRequest dto.RegistrationRequest
	if err := c.Bind(&registerRequest); err != nil {
		return exception.NewBadReqeustError("error to parse data for registration")
	}

	err := h.ms.RegistrationMember(c.Request().Context(), registerRequest)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, dto.BaseResponse{Data: "Success Registration new member"})
}

func (h *MemberHandler) GetMemberInfo(c echo.Context) error {

	param := c.Param("id")
	id, err := strconv.Atoi(param)
	if err != nil {
		return err
	}

	memberInfo, err := h.ms.GetMemberInfo(c.Request().Context(), int64(id))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, dto.BaseResponse{
		Data: memberInfo,
	})
}

func (h *MemberHandler) LoginMember(c echo.Context) error {
	var loginRequest dto.LoginRequest
	if err := c.Bind(&loginRequest); err != nil {
		return exception.NewBadReqeustError("error to parse data for login")
	}

	log.Println(loginRequest)
	memberLogin, err := h.ms.LoginMember(c.Request().Context(), loginRequest)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, dto.BaseResponse{
		Data: memberLogin,
	})
}

func Handler(e *echo.Echo, handler *MemberHandler) {
	e.GET("/v1/member/info/:id", handler.GetMemberInfo)
	e.POST("/v1/member/create", handler.CreateMember)
	e.POST("/v1/member/login", handler.LoginMember)
}
