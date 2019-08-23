package main

import (
	"context"
	"github.com/mkorolyov/core/config"
	"github.com/mkorolyov/core/server"
	profile "github.com/mkorolyov/profiles"
)

func main() {
	profileService := profile.NewService()
	loader := config.Configure()
	s := server.New(loader, profileService)
	s.Serve(context.Background())
}
