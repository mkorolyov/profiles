package main

import (
	"context"
	"github.com/mkorolyov/core/server"
	profile "github.com/mkorolyov/profiles"
)

func main() {
	profileService := profile.NewService()
	s := server.New(profileService)
	s.Serve(context.Background())
}
