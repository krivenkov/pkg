package auth

import (
	"fmt"
)

var ErrNoUserFound = fmt.Errorf("no user found")
var ErrNoTokenFound = fmt.Errorf("no token found")
