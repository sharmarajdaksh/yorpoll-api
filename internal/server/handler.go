package server

import "github.com/sharmarajdaksh/yorpoll-api/internal/db"

type handler struct {
	dbc db.Connection
}
