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

type HealthPoint struct {
	CurrentHP int
	MaxHP     int
	TempHP    int
	HitDice   map[dice.Dice]*Amount
}

func New(maxHp int, hitDice dice.Dice) *HealthPoint {
	hitDiceMap := make(map[dice.Dice]*Amount, 1)
	hitDiceMap[hitDice] = &Amount{MaxAvailable: 1, Available: 1}
	return &HealthPoint{
		CurrentHP: maxHp,
		MaxHP:     maxHp,
		TempHP:    0,
		HitDice:   hitDiceMap,
	}
}

func (hp *HealthPoint) SetMaxHP(maxHP int) {
	hp.MaxHP = maxHP
	if hp.CurrentHP > maxHP {
		hp.CurrentHP = maxHP
	}
}

func (hp *HealthPoint) AddHitDice(hitDiceType dice.Dice) {
	if dice, ok := hp.HitDice[hitDiceType]; ok {
		dice.MaxAvailable++
		dice.Available++
	} else {
		hp.HitDice[hitDiceType] = &Amount{MaxAvailable: 1, Available: 1}
	}
}

func (hp *HealthPoint) RemoveHitDice(hitDiceType dice.Dice) error {
	if dice, ok := hp.HitDice[hitDiceType]; ok {
		dice.MaxAvailable--
		if dice.Available > dice.MaxAvailable {
			dice.Available = dice.MaxAvailable
		}
		if dice.MaxAvailable == 0 {
			delete(hp.HitDice, hitDiceType)
		}
		return nil
	} else {
		return ErrWrongTypeHitDice
	}
}

func (hp *HealthPoint) RollHitDiceRest(rollingDice map[dice.Dice]int) (result []int, err error) {
	for diceType, needAmount := range rollingDice {
		if amount, ok := hp.HitDice[diceType]; ok {
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

func (hp *HealthPoint) ResetHitDice(dicesReset map[dice.Dice]int) error {
	for diceType, resetAmount := range dicesReset {
		if amount, ok := hp.HitDice[diceType]; ok {
			if err := amount.resetDice(resetAmount); err != nil {
				return err
			}
		} else {
			return ErrWrongTypeHitDice
		}
	}
	return nil
}

func (hp *HealthPoint) AddTempHP(tempHP int) {
	hp.TempHP = max(tempHP, hp.TempHP)
}

func (hp *HealthPoint) TakeDamage(damage int) {
	if hp.TempHP >= damage {
		hp.TempHP -= damage
	} else {
		damage -= hp.TempHP
		hp.TempHP = 0
		if hp.CurrentHP > damage {
			hp.CurrentHP -= damage
		} else {
			hp.CurrentHP = 0
		}
	}
}

func (hp *HealthPoint) Heal(heal int) {
	hp.CurrentHP = min(hp.CurrentHP+heal, hp.MaxHP)
}
