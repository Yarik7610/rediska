package memory

import "fmt"

func handleRangeIndexes(startIdx, stopIdx, l int) (int, int, error) {
	if startIdx < 0 {
		startIdx += l
		if startIdx < 0 {
			startIdx = 0
		}
	}
	if stopIdx < 0 {
		stopIdx += l
		if stopIdx >= l {
			stopIdx = l - 1
		}
	}

	if stopIdx < startIdx {
		return 0, 0, fmt.Errorf("start index (%d) is bigger than stop index (%d)", startIdx, stopIdx)
	}

	if startIdx >= l {
		return 0, 0, fmt.Errorf("start index (%d) is bigger or equal than len (%d)", startIdx, l)
	}
	if stopIdx >= l {
		stopIdx = l - 1
	}

	return startIdx, stopIdx, nil
}
