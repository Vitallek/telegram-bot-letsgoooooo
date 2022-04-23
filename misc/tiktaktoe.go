package misc

import "strconv"

func getTikTakToePatterns() [8][3][3]int {
	var patterns [8][3][3]int
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			patterns[i][i][j] = 1
			patterns[i + 3][j][i] = 1
		}
		patterns[6][i][i] = 1
		patterns[7][2 - i][i] = 1
	}
	return patterns
}

func TikTakToeCheckForWin(field [][]string) int {
	patterns := getTikTakToePatterns()
	for _, player := range [2]string{"1", "2"} {
		for _, pattern := range patterns {
			win := true
			for i := 0; i < 3; i++ {
				for j := 0; j < 3; j++ {
					if pattern[i][j] == 0 {
						continue
					}
					if field[i][j] != player {
						win = false
						break
					}
				}
				if !win {
					break
				}
			}
			if win {
				playerInt, _ := strconv.Atoi(player)
				return playerInt
			}
		}
	}
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if field[i][j] == "0" {
				return 0
			}
		}
	}
	return 3
}
