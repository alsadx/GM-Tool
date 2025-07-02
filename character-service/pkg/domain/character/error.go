package character

import "errors"

var ErrNotEnoughLvlForAddHidDice = errors.New("not enough lvl for add hit dice")

var ErrInvalidCharacterName = errors.New("invalid or empty character name")

var ErrInvalidOwnerID = errors.New("invalid or empty owner id")
