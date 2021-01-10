package server

import (
	"fmt"

	"github.com/sharmarajdaksh/yorpoll-api/internal/poll"
)

var idValidationString = fmt.Sprintf("required,len=%d", poll.IDFieldsSize)
