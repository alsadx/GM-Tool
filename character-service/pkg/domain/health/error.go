package health

import "errors"

var ErrNoHitDiceAvailable = errors.New("not available hit dice")
var ErrWrongTypeHitDice = errors.New("wrong type of hit dice")
var ErrCantResetHitDice = errors.New("you can't recover more hit dice than you have")
