package level

import "github.com/alsadx/GM-Tool/character-service/gen"

var Thresholds = []int{
	0, 300, 900, 2700, 6500, 14000, 23000, 34000,
	48000, 64000, 85000, 100000, 120000, 140000,
	165000, 195000, 225000, 265000, 305000, 355000,
}

type LevelSystem struct {
	currentLevel  int
	earnedLevel   int
	currentExp    int
	nextThreshold int
}

func NewLevelSystem() *LevelSystem {

	return &LevelSystem{
		currentLevel:  1,
		earnedLevel:   1,
		currentExp:    0,
		nextThreshold: Thresholds[1],
	}
}

func (ls *LevelSystem) ToProto() *gen.LevelSystem {
	return &gen.LevelSystem{
		CurrentLevel:  int32(ls.currentLevel),
		EarnedLevel:   int32(ls.earnedLevel),
		CurrentExp:    int32(ls.currentExp),
		NextThreshold: int32(ls.nextThreshold),
	}
}

func FromProto(protoLevelSystem *gen.LevelSystem) *LevelSystem {
	return &LevelSystem{
		currentLevel:  int(protoLevelSystem.CurrentLevel),
		earnedLevel:   int(protoLevelSystem.EarnedLevel),
		currentExp:    int(protoLevelSystem.CurrentExp),
		nextThreshold: int(protoLevelSystem.NextThreshold),
	}
}

func (ls *LevelSystem) AddExp(amount int) {
	if amount <= 0 || ls.earnedLevel >= len(Thresholds) {
		return
	}

	ls.currentExp += amount

	if ls.currentExp >= ls.nextThreshold {
		ls.updateEarnedLevel()
	}
}

func (ls *LevelSystem) SetLevel(lvl int) {
	ls.currentExp = Thresholds[lvl-1]
	ls.currentLevel = lvl
	ls.updateEarnedLevel()
}

func (ls *LevelSystem) RemoveExp(amount int) {
	if amount <= 0 || ls.currentExp == 0 {
		return
	}

	ls.currentExp = max(ls.currentExp-amount, 0)

	if ls.currentExp < Thresholds[ls.earnedLevel-1] {
		ls.updateEarnedLevelDown()
	}
}

func (ls *LevelSystem) updateEarnedLevelDown() {
	newEarnedLevel := 1

	for i := len(Thresholds) - 1; i >= 0; i-- {
		if ls.currentExp >= Thresholds[i] {
			newEarnedLevel = i + 1
			break
		}
	}

	if newEarnedLevel != ls.earnedLevel {
		ls.earnedLevel = newEarnedLevel
		if ls.earnedLevel < len(Thresholds) {
			ls.nextThreshold = Thresholds[ls.earnedLevel]
		} else {
			ls.nextThreshold = 0
		}
	}
}

func (ls *LevelSystem) updateEarnedLevel() {
	newEarnedLevel := ls.earnedLevel

	for i := ls.earnedLevel; i < len(Thresholds); i++ {
		if ls.currentExp >= Thresholds[i] {
			newEarnedLevel = i + 1
		} else {
			break
		}
	}

	if newEarnedLevel != ls.earnedLevel {
		ls.earnedLevel = newEarnedLevel

		if ls.earnedLevel < len(Thresholds) {
			ls.nextThreshold = Thresholds[ls.earnedLevel]
		} else {
			ls.nextThreshold = 0
		}
	}
}

func (ls *LevelSystem) LevelUp() bool {
	if !ls.CanLevelUp() {
		return false
	}

	ls.currentLevel++
	return true
}

func (ls *LevelSystem) LevelDown() bool {
	if !ls.CanLevelDown() {
		return false
	}

	ls.currentLevel--
	return true
}

func (ls *LevelSystem) CanLevelUp() bool   { return ls.currentLevel < ls.earnedLevel }
func (ls *LevelSystem) CanLevelDown() bool { return ls.currentLevel > ls.earnedLevel }
func (ls *LevelSystem) CurrentLevel() int  { return ls.currentLevel }
func (ls *LevelSystem) EarnedLevel() int   { return ls.earnedLevel }
func (ls *LevelSystem) CurrentExp() int    { return ls.currentExp }
func (ls *LevelSystem) ExpToNextLevel() int {
	if ls.currentExp >= Thresholds[19] {
		return 0
	}
	return ls.nextThreshold - ls.currentExp
}
