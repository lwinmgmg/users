package server

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/golang-jwt/jwt/v5"
	gmodels "github.com/lwinmgmg/gmodels/golang/models"
	"github.com/lwinmgmg/user/middlewares"
	"github.com/lwinmgmg/user/models"
	"github.com/lwinmgmg/user/services"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	TokenErrorCode  codes.Code = 400
	UnauthorizeCode codes.Code = 402
)

var (
	DB = services.PgDb
)

type UserServer struct {
	gmodels.UnimplementedUserServiceServer
}

func (userServer *UserServer) GetProfile(ctx context.Context, req *gmodels.GetProfileRequest) (*gmodels.User, error) {
	var user models.User
	var claim jwt.RegisteredClaims
	if err := middlewares.ValidateToken(req.Token, middlewares.DefaultTokenKey, &claim); err != nil {
		return nil, status.Error(UnauthorizeCode, err.Error())
	}
	if _, err := user.GetPartnerByCode(claim.Subject, DB); err != nil {
		return nil, status.Error(TokenErrorCode, err.Error())
	}
	return &gmodels.User{
		Username:        user.Username,
		Code:            user.Code,
		IsAuthenticator: user.IsAuthenticator,
		Is2FA:           user.OtpUrl != "",
		Partner: &gmodels.Partner{
			FirstName:        user.Partner.FirstName,
			LastName:         user.Partner.LastName,
			Email:            user.Partner.Email,
			IsEmailConfirmed: user.Partner.IsEmailConfirmed,
			Phone:            user.Partner.Phone,
			IsPhoneConfirmed: user.Partner.IsPhoneConfirmed,
			Code:             user.Partner.Code,
		},
	}, nil
}

func (userServer *UserServer) GetUserByCode(ctx context.Context, req *gmodels.GetUserByCodeRequest) (*gmodels.User, error) {
	var user models.User
	if _, err := user.GetPartnerByCode(req.Code, DB); err != nil {
		return nil, status.Errorf(codes.Unknown, err.Error())
	}
	return &gmodels.User{
		Username:        user.Username,
		Code:            user.Code,
		IsAuthenticator: user.IsAuthenticator,
		Is2FA:           user.OtpUrl != "",
		Partner: &gmodels.Partner{
			FirstName:        user.Partner.FirstName,
			LastName:         user.Partner.LastName,
			Email:            user.Partner.Email,
			IsEmailConfirmed: user.Partner.IsEmailConfirmed,
			Phone:            user.Partner.Phone,
			IsPhoneConfirmed: user.Partner.IsPhoneConfirmed,
			Code:             user.Partner.Code,
		},
	}, nil
}

func StartServer(host string, port int) {
	lis, err := net.Listen("tcp", fmt.Sprintf("%v:%v", host, port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	gmodels.RegisterUserServiceServer(s, &UserServer{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
