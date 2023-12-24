package common

import "errors"

var ErrNotFound = errors.New("NotFound")
var ErrChannelDeadlock = errors.New("ChannelDeadlock")
