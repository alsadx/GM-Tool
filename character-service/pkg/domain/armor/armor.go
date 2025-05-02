package armor

type Type int

const (
	NoArmor Type = iota
	LightArmor
	MediumArmor
	HeavyArmor
)

type Armor struct {
	Name           string
	Type           Type
	BaseAC         int
	MaxDex         int
	StrengthReq    int
	StealthPenalty bool
}

func (a *Armor) CalculateAC(dexMod int) int {
	switch a.Type {
	case LightArmor:
		return a.BaseAC + dexMod
	case MediumArmor:
		return a.BaseAC + min(dexMod, a.MaxDex)
	case HeavyArmor:
		return a.BaseAC
	}
	return 10 + dexMod
}

func (a *Armor) CheckRequirements(strenght int) bool {
	if strenght < a.StrengthReq {
		return true
	}
	return false
}
