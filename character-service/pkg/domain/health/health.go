package health

type HealthDice int

const (
	k6 HealthDice = 6 + (iota*2)
	k8
	k10
	k12
)

type hitDicePool struct {
	Available map[HealthDice]int
}

type HealthPoint struct {
	currentHP int
	maxHP int
	tempHP int
	healthDice map[HealthDice]
}

func (hp *HealthPoint) SetMaxHP(maxHP int) {
	hp.maxHP = maxHP
}

func (hp *HealthPoint) RollHealthDice()