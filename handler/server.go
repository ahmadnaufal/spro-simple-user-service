package handler

import "github.com/SawitProRecruitment/UserService/repository"

type Server struct {
	Repository repository.RepositoryInterface
	JWT        JWT
}

type NewServerOptions struct {
	Repository repository.RepositoryInterface
	JWT        JWT
}

func NewServer(opts NewServerOptions) *Server {
	return &Server{
		Repository: opts.Repository,
		JWT:        opts.JWT,
	}
}
