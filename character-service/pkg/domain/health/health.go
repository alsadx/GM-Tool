package health

import (
	"github.com/alsadx/GM-Tool/character-service/pkg/domain/dice"
)

type amount struct {
	maxAvailable int
	available    int
}

func (a *amount) useDice(amountUsed int) error {
	if a.available >= amountUsed {
		a.available -= amountUsed
		return nil
	}
	return ErrNoHitDiceAvailable
}

type HealthPoint struct {
	currentHP int
	maxHP     int
	tempHP    int
	hitDice   map[dice.Dice]amount
}

func (hp *HealthPoint) SetMaxHP(maxHP int) {
	hp.maxHP = maxHP
}

func (hp *HealthPoint) AddHitDice(hitDiceType dice.Dice) {
	if dice, ok := hp.hitDice[hitDiceType]; ok {
		dice.maxAvailable++
		dice.available++
	} else {
		hp.hitDice[hitDiceType] = amount{maxAvailable: 1, available: 1}
	}
}

func (hp *HealthPoint) RemoveHitDice(hitDiceType dice.Dice) error {
	if dice, ok := hp.hitDice[hitDiceType]; ok {
		dice.maxAvailable--
		if dice.available > dice.maxAvailable {
			dice.available = dice.maxAvailable
		}
		if dice.maxAvailable == 0 {
			delete(hp.hitDice, hitDiceType)
		}
		return nil
	} else {
		return ErrWrongTypeHitDice
	}
}

func (hp *HealthPoint) RollHitDiceRest(rollingDice map[dice.Dice]int) (result []int, err error) {
	for diceType, needAmount := range rollingDice {
		if amount, ok := hp.hitDice[diceType]; ok {
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

func (hp *HealthPoint) AddTempHP(tempHP int) {
	hp.tempHP = max(tempHP, hp.tempHP)
}

func (hp *HealthPoint) TakeDamage(damage int) {
	if hp.tempHP >= damage {
		hp.tempHP -= damage
	} else {
		damage -= hp.tempHP
		hp.tempHP = 0
		if hp.currentHP > damage {
			hp.currentHP -= damage
		} else {
			hp.currentHP = 0
		}
	}
}
