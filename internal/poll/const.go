package poll

import "fmt"

// ErrPollExpired implies that the poll being accessed is expired and cannot be modified
var ErrPollExpired = fmt.Errorf("the poll being accessed is expired")

// IDFieldsSize represents the size of ID fields for various objects
const IDFieldsSize = 36
