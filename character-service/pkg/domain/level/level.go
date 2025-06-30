package level

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

func (ls *LevelSystem) WithCurrentExp(current_exp int) *LevelSystem {
	ls.currentExp = current_exp
	ls.updateEarnedLevel()
	return ls
}

func (ls *LevelSystem) WithCurrentLvl(current_lvl int) *LevelSystem {
	current_lvl = min(max(current_lvl, 1), len(Thresholds))
	ls.currentLevel = current_lvl
	return ls
}

func (ls *LevelSystem) WithNextThreshold(next_threshols int) *LevelSystem {
	ls.nextThreshold = next_threshols
	return ls
}

func (ls *LevelSystem) AddExp(amount int) {
	if ls.earnedLevel >= len(Thresholds) {
		return
	}

	ls.currentExp += amount

	if ls.currentExp >= ls.nextThreshold {
		ls.updateEarnedLevel()
	}
}

func (ls *LevelSystem) SetLevel(lvl int) {
    lvl = min(max(lvl, 1), len(Thresholds))
    ls.currentExp = Thresholds[lvl-1]
    ls.currentLevel = lvl
    ls.updateEarnedLevel()
}

func (ls *LevelSystem) RemoveExp(amount int) {
    if amount <= 0 || ls.currentExp == 0 {
        return
    }
    ls.currentExp = max(ls.currentExp-amount, 0)
    ls.updateEarnedLevel()
}

func (ls *LevelSystem) updateEarnedLevel() {
    newEarnedLevel := 1
    for i := 1; i < len(Thresholds); i++ {
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
	if ls.earnedLevel >= len(Thresholds) {
		return 0
	}
	return ls.nextThreshold - ls.currentExp
}
