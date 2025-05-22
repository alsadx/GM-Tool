package level

type LevelSystem struct {
	currentLevel  int
	earnedLevel   int
	currentExp    int
	nextThreshold int
	thresholds    []int
}

func NewLevelSystem() *LevelSystem {
	thresholds := []int{
		0, 300, 900, 2700, 6500, 14000, 23000, 34000,
		48000, 64000, 85000, 100000, 120000, 140000,
		165000, 195000, 225000, 265000, 305000, 355000,
	}

	return &LevelSystem{
		currentLevel:  1,
		earnedLevel:   1,
		currentExp:    0,
		nextThreshold: thresholds[1],
		thresholds:    thresholds,
	}
}

func (ls *LevelSystem) AddExp(amount int) {
	if amount <= 0 || ls.earnedLevel >= len(ls.thresholds) {
		return
	}

	ls.currentExp += amount

	if ls.currentExp >= ls.nextThreshold {
		ls.updateEarnedLevel()
	}
}

func (ls *LevelSystem) SetLevel(lvl int) {
	ls.currentExp = ls.thresholds[lvl-1]
	ls.currentLevel = lvl
	ls.updateEarnedLevel()
}

func (ls *LevelSystem) RemoveExp(amount int) {
	if amount <= 0 || ls.currentExp == 0 {
		return
	}

	ls.currentExp = max(ls.currentExp-amount, 0)

	if ls.currentExp < ls.thresholds[ls.earnedLevel-1] {
		ls.updateEarnedLevelDown()
	}
}

func (ls *LevelSystem) updateEarnedLevelDown() {
	newEarnedLevel := 1

	for i := len(ls.thresholds) - 1; i >= 0; i-- {
		if ls.currentExp >= ls.thresholds[i] {
			newEarnedLevel = i + 1
			break
		}
	}

	if newEarnedLevel != ls.earnedLevel {
		ls.earnedLevel = newEarnedLevel
		if ls.earnedLevel < len(ls.thresholds) {
			ls.nextThreshold = ls.thresholds[ls.earnedLevel]
		} else {
			ls.nextThreshold = 0
		}
	}
}

func (ls *LevelSystem) updateEarnedLevel() {
	newEarnedLevel := ls.earnedLevel

	for i := ls.earnedLevel; i < len(ls.thresholds); i++ {
		if ls.currentExp >= ls.thresholds[i] {
			newEarnedLevel = i + 1
		} else {
			break
		}
	}

	if newEarnedLevel != ls.earnedLevel {
		ls.earnedLevel = newEarnedLevel

		if ls.earnedLevel < len(ls.thresholds) {
			ls.nextThreshold = ls.thresholds[ls.earnedLevel]
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

func (ls *LevelSystem) CanLevelUp() bool    { return ls.currentLevel < ls.earnedLevel }
func (ls *LevelSystem) CanLevelDown() bool  { return ls.currentLevel > ls.earnedLevel }
func (ls *LevelSystem) CurrentLevel() int   { return ls.currentLevel }
func (ls *LevelSystem) EarnedLevel() int    { return ls.earnedLevel }
func (ls *LevelSystem) CurrentExp() int     { return ls.currentExp }
func (ls *LevelSystem) ExpToNextLevel() int { return ls.nextThreshold - ls.currentExp }
