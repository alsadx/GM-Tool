package level_test

import (
	"testing"

	"github.com/alsadx/GM-Tool/character-service/pkg/domain/level"
	"github.com/stretchr/testify/assert"
)

func TestNewLevelSystem(t *testing.T) {
	tests := []struct {
		name           string
		wantLevel      int
		wantEarned     int
		wantExp        int
		wantNextThresh int
	}{
		{
			name:           "Initial state",
			wantLevel:      1,
			wantEarned:     1,
			wantExp:        0,
			wantNextThresh: 300,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ls := level.NewLevelSystem()
			assert.Equal(t, tt.wantLevel, ls.CurrentLevel())
			assert.Equal(t, tt.wantEarned, ls.EarnedLevel())
			assert.Equal(t, tt.wantExp, ls.CurrentExp())
			assert.Equal(t, tt.wantNextThresh, ls.ExpToNextLevel())
		})
	}
}

func TestAddExp(t *testing.T) {
	tests := []struct {
		name           string
		setup          func() *level.LevelSystem
		expToAdd       int
		wantLevel      int
		wantEarned     int
		wantExp        int
		wantNextThresh int
	}{
		{
			name:           "Add exp below threshold",
			setup:          level.NewLevelSystem,
			expToAdd:       200,
			wantLevel:      1,
			wantEarned:     1,
			wantExp:        200,
			wantNextThresh: 100,
		},
		{
			name: "Add exp crossing multiple thresholds",
			setup: func() *level.LevelSystem {
				ls := level.NewLevelSystem()
				ls.SetLevel(3)
				return ls
			},
			expToAdd:       2000,
			wantLevel:      3,
			wantEarned:     4,
			wantExp:        900 + 2000, // 2900
			wantNextThresh: 6500 - 2900,
		},
		{
			name: "Add exp to max level",
			setup: func() *level.LevelSystem {
				ls := level.NewLevelSystem()
				ls.SetLevel(len(level.Thresholds))
				return ls
			},
			expToAdd:       100000,
			wantLevel:      len(level.Thresholds),
			wantEarned:     len(level.Thresholds),
			wantExp:        level.Thresholds[len(level.Thresholds)-1],
			wantNextThresh: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ls := tt.setup()
			ls.AddExp(tt.expToAdd)
			assert.Equal(t, tt.wantLevel, ls.CurrentLevel())
			assert.Equal(t, tt.wantEarned, ls.EarnedLevel())
			assert.Equal(t, tt.wantExp, ls.CurrentExp())
			assert.Equal(t, tt.wantNextThresh, ls.ExpToNextLevel())
		})
	}
}

func TestSetLevel(t *testing.T) {
	tests := []struct {
		name           string
		setup          func() *level.LevelSystem
		levelToSet     int
		wantLevel      int
		wantEarned     int
		wantExp        int
		wantNextThresh int
	}{
		{
			name:           "Set level higher than current",
			setup:          level.NewLevelSystem,
			levelToSet:     5,
			wantLevel:      5,
			wantEarned:     5,
			wantExp:        level.Thresholds[4],                       // 6500
			wantNextThresh: level.Thresholds[5] - level.Thresholds[4], // 14000 - 6500 = 7500
		},
		{
			name: "Set level lower than earned",
			setup: func() *level.LevelSystem {
				ls := level.NewLevelSystem()
				ls.SetLevel(5)
				return ls
			},
			levelToSet:     3,
			wantLevel:      3,
			wantEarned:     5,
			wantExp:        level.Thresholds[2], // 900
			wantNextThresh: level.Thresholds[5] - level.Thresholds[2],
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ls := tt.setup()
			ls.SetLevel(tt.levelToSet)
			assert.Equal(t, tt.wantLevel, ls.CurrentLevel())
			assert.Equal(t, tt.wantEarned, ls.EarnedLevel())
			assert.Equal(t, tt.wantExp, ls.CurrentExp())
			assert.Equal(t, tt.wantNextThresh, ls.ExpToNextLevel())
		})
	}
}

func TestRemoveExp(t *testing.T) {
	tests := []struct {
		name           string
		setup          func() *level.LevelSystem
		expToRemove    int
		wantLevel      int
		wantEarned     int
		wantExp        int
		wantNextThresh int
	}{
		{
			name: "Remove exp without level change",
			setup: func() *level.LevelSystem {
				ls := level.NewLevelSystem()
				ls.AddExp(500)
				return ls
			},
			expToRemove:    400,
			wantLevel:      1,
			wantEarned:     1,
			wantExp:        100,
			wantNextThresh: 200,
		},
		{
			name: "Remove exp with level down",
			setup: func() *level.LevelSystem {
				ls := level.NewLevelSystem()
				ls.SetLevel(5)
				return ls
			},
			expToRemove:    6000,
			wantLevel:      5,
			wantEarned:     2,
			wantExp:        6500 - 6000, // 500
			wantNextThresh: 900 - 500,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ls := tt.setup()
			ls.RemoveExp(tt.expToRemove)
			assert.Equal(t, tt.wantLevel, ls.CurrentLevel())
			assert.Equal(t, tt.wantEarned, ls.EarnedLevel())
			assert.Equal(t, tt.wantExp, ls.CurrentExp())
			assert.Equal(t, tt.wantNextThresh, ls.ExpToNextLevel())
		})
	}
}

func TestLevelUpDown(t *testing.T) {
	tests := []struct {
		name        string
		setup       func() *level.LevelSystem
		actions     func(*level.LevelSystem)
		wantLevel   int
		wantCanUp   bool
		wantCanDown bool
	}{
		{
			name: "Level up when possible",
			setup: func() *level.LevelSystem {
				ls := level.NewLevelSystem()
				ls.AddExp(300)
				return ls
			},
			actions: func(ls *level.LevelSystem) {
				ls.LevelUp()
			},
			wantLevel:   2,
			wantCanUp:   false,
			wantCanDown: false,
		},
		{
			name: "Сan raise the level",
			setup: func() *level.LevelSystem {
				ls := level.NewLevelSystem()
				ls.AddExp(300)
				return ls
			},
			wantLevel:   1,
			wantCanUp:   true,
			wantCanDown: false,
		},
		{
			name: "Level down when possible",
			setup: func() *level.LevelSystem {
				ls := level.NewLevelSystem()
				ls.SetLevel(3)
				return ls
			},
			actions: func(ls *level.LevelSystem) {
				ls.RemoveExp(900)
				ls.LevelDown()
				ls.LevelDown()
			},
			wantLevel:   1,
			wantCanUp:   false,
			wantCanDown: false,
		},
		{
			name: "Level dowwn possible",
			setup: func() *level.LevelSystem {
				ls := level.NewLevelSystem()
				ls.SetLevel(3)
				ls.RemoveExp(900)
				return ls
			},
			wantLevel:   3,
			wantCanUp:   false,
			wantCanDown: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ls := tt.setup()
			if tt.actions != nil {
				tt.actions(ls)
			}
			assert.Equal(t, tt.wantLevel, ls.CurrentLevel())
			assert.Equal(t, tt.wantCanUp, ls.CanLevelUp())
			assert.Equal(t, tt.wantCanDown, ls.CanLevelDown())
		})
	}
}
