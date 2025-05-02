package health

import "github.com/alsadx/GM-Tool/character-service/pkg/domain/dice"

type amount struct {
	maxAvailable int
	available int
}

func (a *amount) useDice() error {
	if a.available > 0 {
		a.available--
		return nil
	}
	return ErrNoHitDiceAvailable
}

type HealthPoint struct {
	currentHP int
	maxHP int
	tempHP int
	hitDice map[dice.Dice]amount
}

func (hp *HealthPoint) SetMaxHP(maxHP int) {
	hp.maxHP = maxHP
}

func (hp *HealthPoint) rollHitDice(hitDiceType dice.Dice) (int, error) {
	diceAmount, ok := hp.hitDice[hitDiceType]
	if !ok {
		return 0, WrongTypeHitDice
	}
	if err := diceAmount.useDice(); err != nil {
		return 0, err
	}
	return dice.RollDice(hitDiceType), nil
}

func (hp *HealthPoint) AddHitDice(hitDiceType dice.Dice) {
	if dice, ok := hp.hitDice[hitDiceType]; ok {
		dice.maxAvailable++
		dice.available++
	} else {
		hp.hitDice[hitDiceType] = amount{maxAvailable: 1, available: 1}
	}
}