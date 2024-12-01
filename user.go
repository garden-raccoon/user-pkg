package user

import (
	"context"
	"fmt"
	"github.com/garden-raccoon/user-pkg/models"
	"time"

	proto "github.com/garden-raccoon/user-pkg/protocols/user"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

const timeOut = 60

type IUserAPI interface {
	// SignUp is
	SignUp(email string, password []byte, isEmployer bool) ([]byte, error)

	// CheckAuth is
	CheckAuth(token []byte) (*models.User, error)

	// SignIn is
	SignIn(email string, password []byte) ([]byte, error)

	// Close GRPC Api connection
	Close() error
}

// UsersAPI is profile-service GRPC UsersAPI
// structure with client Connection
type UsersAPI struct {
	addr    string
	timeout time.Duration
	*grpc.ClientConn
	proto.UserServiceClient
}

// New create new Users IEmployerAPI instance
func New(addr string) (IUserAPI, error) {
	api := &UsersAPI{timeout: timeOut * time.Second}

	if err := api.initConn(addr); err != nil {
		return nil, fmt.Errorf("create Users UsersAPI:  %w", err)
	}

	api.UserServiceClient = proto.NewUserServiceClient(api.ClientConn)
	return api, nil
}

// initConn initialize connection to Grpc servers
func (api *UsersAPI) initConn(addr string) (err error) {
	var kacp = keepalive.ClientParameters{
		Time:                10 * time.Second, // send pings every 10 seconds if there is no activity
		Timeout:             time.Second,      // wait 1 second for ping ack before considering the connection dead
		PermitWithoutStream: true,             // send pings even without active streams
	}

	api.ClientConn, err = grpc.Dial(addr, grpc.WithInsecure(), grpc.WithKeepaliveParams(kacp))
	return
}
func (api *UsersAPI) CheckAuth(token []byte) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), api.timeout)
	defer cancel()
	fmt.Printf("\n %s token \n", token)
	resp, err := api.UserServiceClient.CheckAuth(ctx, &proto.TokenRequest{Token: token})
	if err != nil {
		return nil, fmt.Errorf("checkAuth api request: %w", err)
	}

	return models.UserFromProto(resp), nil
}

// SignUp is
func (api *UsersAPI) SignUp(email string, password []byte, isEmployer bool) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), api.timeout)
	defer cancel()

	opts := &proto.SignUpRequest{
		Email:      email,
		Password:   password,
		IsEmployer: isEmployer,
	}
	fmt.Printf("epts email is %s\n", opts.Email)
	resp, err := api.UserServiceClient.SignUp(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("signUp api request has been failed: %w", err)
	}
	return resp.Token, nil
}

// SignIn is
func (api *UsersAPI) SignIn(email string, password []byte) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), api.timeout)
	defer cancel()

	opts := &proto.SignInRequest{
		Email:    email,
		Password: password,
	}

	resp, err := api.UserServiceClient.SignIn(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("signIn api request: %w", err)
	}

	return resp.Token, nil
}
