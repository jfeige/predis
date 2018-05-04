package predis

import (
	"errors"
	"fmt"
)

var(
	ErrNil = errors.New("predis: nil returned!")
	ErrUnknown = errors.New("predis: unknown response type!")
)


func checkErrResponse(tp interface{})error{
	return fmt.Errorf("predis: unknown response type:%T",tp)
}