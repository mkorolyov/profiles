package profile

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
	"math/rand"
	"strconv"
	"time"
)

type Profile struct {
	FirstName string
	LastName  string
}

type ProfileService struct {
	storage map[string]Profile
}

func (profileService ProfileService) GRPCRegisterer() func(s *grpc.Server) {
	return func(s *grpc.Server) {
		RegisterProfileServer(s, profileService)
	}
}

func (profileService ProfileService) HTTPRegisterer() func(ctx context.Context, mux *runtime.ServeMux, endpoint string,
	opts []grpc.DialOption) (err error) {
	return RegisterProfileHandlerFromEndpoint
}

func NewService() *ProfileService {
	rand.Seed(time.Now().UnixNano())
	return &ProfileService{storage: map[string]Profile{}}
}

func (ps ProfileService) Get(ctx context.Context, r *GetRequest) (*GetResponse, error) {
	p, ok := ps.storage[r.Id]
	if !ok {
		return nil, fmt.Errorf("no profile with id %s", r.Id)
	}

	return &GetResponse{FirstName: p.FirstName, LastName: p.LastName}, nil
}

func (ps ProfileService) Create(ctx context.Context, r *CreateRequest) (*CreateResponse, error) {
	id := strconv.Itoa(rand.Int())
	ps.storage[id] = Profile{
		FirstName: r.FirstName,
		LastName:  r.LastName,
	}
	return &CreateResponse{Id: id}, nil
}
