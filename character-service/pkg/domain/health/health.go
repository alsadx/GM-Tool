package health

import (
	"github.com/alsadx/GM-Tool/character-service/pkg/domain/dice"
)

type Amount struct {
	MaxAvailable int
	Available    int
}

func (a *Amount) useDice(amountUsed int) error {
	if a.Available >= amountUsed {
		a.Available -= amountUsed
		return nil
	}
	return ErrNoHitDiceAvailable
}

func (a *Amount) resetDice(amountReset int) error {
	if a.Available+amountReset <= a.MaxAvailable {
		a.Available += amountReset
		return nil
	}
	return ErrCantResetHitDice
}

type Health struct {
	CurrentHP int
	MaxHP     int
	TempHP    int
	HitDice   map[dice.Dice]*Amount
}

func New(maxHp int, hitDice dice.Dice) *Health {
	hitDiceMap := make(map[dice.Dice]*Amount, 1)
	hitDiceMap[hitDice] = &Amount{MaxAvailable: 1, Available: 1}
	return &Health{
		CurrentHP: maxHp,
		MaxHP:     maxHp,
		TempHP:    0,
		HitDice:   hitDiceMap,
	}
}

func (h *Health) SetMaxHP(maxHP int) {
	h.MaxHP = maxHP
	if h.CurrentHP > maxHP {
		h.CurrentHP = maxHP
	}
}

func (h *Health) AddHitDice(hitDiceType dice.Dice) {
	if dice, ok := h.HitDice[hitDiceType]; ok {
		dice.MaxAvailable++
		dice.Available++
	} else {
		h.HitDice[hitDiceType] = &Amount{MaxAvailable: 1, Available: 1}
	}
}

func (h *Health) RemoveHitDice(hitDiceType dice.Dice) error {
	if dice, ok := h.HitDice[hitDiceType]; ok {
		dice.MaxAvailable--
		if dice.Available > dice.MaxAvailable {
			dice.Available = dice.MaxAvailable
		}
		if dice.MaxAvailable == 0 {
			delete(h.HitDice, hitDiceType)
		}
		return nil
	} else {
		return ErrWrongTypeHitDice
	}
}

func (h *Health) RollHitDiceRest(rollingDice map[dice.Dice]int) (result []int, err error) {
	for diceType, needAmount := range rollingDice {
		if amount, ok := h.HitDice[diceType]; ok {
			if err = amount.useDice(needAmount); err != nil {
				return nil, ErrNoHitDiceAvailable
			} else {
				result = append(result, dice.MultiRollDice(diceType, needAmount)...)
			}
		} else {
			return nil, ErrWrongTypeHitDice
		}
	}
	return result, nil
}

func (h *Health) ResetHitDice(dicesReset map[dice.Dice]int) error {
	for diceType, resetAmount := range dicesReset {
		if amount, ok := h.HitDice[diceType]; ok {
			if err := amount.resetDice(resetAmount); err != nil {
				return err
			}
		} else {
			return ErrWrongTypeHitDice
		}
	}
	return nil
}

func (h *Health) AddTempHP(tempHP int) {
	h.TempHP = max(tempHP, h.TempHP)
}

func (h *Health) TakeDamage(damage int) {
	if h.TempHP >= damage {
		h.TempHP -= damage
	} else {
		damage -= h.TempHP
		h.TempHP = 0
		if h.CurrentHP > damage {
			h.CurrentHP -= damage
		} else {
			h.CurrentHP = 0
		}
	}
}

func (h *Health) Heal(heal int) {
	h.CurrentHP = min(h.CurrentHP+heal, h.MaxHP)
}

func (h *Health) GetHitDiceCount() int {
	var maxAvailable int
	for _, amount := range h.HitDice {
		maxAvailable += amount.MaxAvailable
	}
	return maxAvailable
}
