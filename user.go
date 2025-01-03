package user

import (
	"context"
	"errors"
	"fmt"
	"github.com/garden-raccoon/user-pkg/models"
	"github.com/gofrs/uuid"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health/grpc_health_v1"
	"sync"
	"time"

	proto "github.com/garden-raccoon/user-pkg/protocols/user"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

const timeOut = 60

type IUserAPI interface {
	// SignUp is
	SignUp(email string, password []byte, userType int) ([]byte, error)

	// CheckAuth is
	CheckAuth(token []byte) (*models.User, error)

	// UserByUUID is
	UserByUUID(userUUID uuid.UUID) (*models.User, error)

	UpdateUser(user *models.UpdateUserRequest) (*models.User, error)

	// SignIn is
	SignIn(email string, password []byte) ([]byte, error)

	HealthCheck() error

	// Close GRPC Api connection
	Close() error
}

// UsersAPI is profile-service GRPC UsersAPI
// structure with client Connection
type UsersAPI struct {
	addr    string
	timeout time.Duration
	mu      sync.Mutex
	*grpc.ClientConn
	proto.UserServiceClient
	grpc_health_v1.HealthClient
}

// New create new Users IEmployerAPI instance
func New(addr string) (IUserAPI, error) {
	api := &UsersAPI{timeout: timeOut * time.Second}

	if err := api.initConn(addr); err != nil {
		return nil, fmt.Errorf("create Users UsersAPI:  %w", err)
	}
	api.HealthClient = grpc_health_v1.NewHealthClient(api.ClientConn)

	api.UserServiceClient = proto.NewUserServiceClient(api.ClientConn)
	return api, nil
}

func (api *UsersAPI) UpdateUser(user *models.UpdateUserRequest) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), api.timeout)
	defer cancel()
	protoUser := models.Proto(*user)
	resp, err := api.UserServiceClient.UpdateUser(ctx, protoUser)
	if err != nil {
		return nil, fmt.Errorf("updateUser api request: %w", err)
	}
	fmt.Printf("resp is %s \n", resp)
	return models.UserFromProto(resp), nil
}

// initConn initialize connection to Grpc servers
func (api *UsersAPI) initConn(addr string) (err error) {
	var kacp = keepalive.ClientParameters{
		Time:                10 * time.Second, // send pings every 10 seconds if there is no activity
		Timeout:             time.Second,      // wait 1 second for ping ack before considering the connection dead
		PermitWithoutStream: true,             // send pings even without active streams
	}

	api.ClientConn, err = grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithKeepaliveParams(kacp))
	if err != nil {
		return fmt.Errorf("failed to dial: %w", err)
	}
	return
}
func (api *UsersAPI) CheckAuth(token []byte) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), api.timeout)
	defer cancel()

	protoToken := &proto.TokenRequest{Token: token}
	resp, err := api.UserServiceClient.CheckAuth(ctx, protoToken)
	if err != nil {
		return nil, fmt.Errorf("checkAuth api request: %w", err)
	}

	return models.UserFromProto(resp), nil
}

// SignUp is
func (api *UsersAPI) SignUp(email string, password []byte, userType int) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), api.timeout)
	defer cancel()

	opts := &proto.SignUpRequest{
		Email:    email,
		Password: password,
		UserType: int64(userType),
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

func (api *UsersAPI) HealthCheck() error {
	ctx, cancel := context.WithTimeout(context.Background(), api.timeout)
	defer cancel()

	api.mu.Lock()
	defer api.mu.Unlock()

	resp, err := api.HealthClient.Check(ctx, &grpc_health_v1.HealthCheckRequest{Service: "userapi"})
	if err != nil {
		return fmt.Errorf("healthcheck error: %w", err)
	}

	if resp.Status != grpc_health_v1.HealthCheckResponse_SERVING {
		return fmt.Errorf("node is %s", errors.New("service is unhealthy"))
	}
	return nil
}

func (api *UsersAPI) UserByUUID(userUUID uuid.UUID) (*models.User, error) {
	opts := &proto.UserGetter{
		Getter: &proto.UserGetter_UserUuid{
			UserUuid: userUUID.Bytes(),
		},
	}
	return api.getUser(opts)
}

func (api *UsersAPI) getUser(opts *proto.UserGetter) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), api.timeout)
	defer cancel()

	resp, err := api.UserServiceClient.UserBy(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("userapi request failed: %w", err)
	}

	return models.UserFromProto(resp), nil
}
