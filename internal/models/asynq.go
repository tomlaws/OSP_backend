package models

import "github.com/hibiken/asynq"

type JobSystem struct {
	Client *asynq.Client
	Server *asynq.Server
	Mux    *asynq.ServeMux
}
