package common

import "errors"

var ErrInternal = errors.New("InternalError")
var ErrNotFound = errors.New("NotFound")
var ErrChannelDeadlock = errors.New("ChannelDeadlock")
