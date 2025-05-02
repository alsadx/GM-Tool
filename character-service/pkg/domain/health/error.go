package health

import "errors"

var ErrNoHitDiceAvailable = errors.New("Not available hit dice")
var WrongTypeHitDice = errors.New("Wrong type of hit dice")