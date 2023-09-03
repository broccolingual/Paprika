package common

type LineController _LineController

type _LineController struct {
	stepIndex int
	stepSize  int
	startPos  []int
}

func NewLineController() *LineController {
	lc := new(LineController)
	lc.stepIndex = 0
	lc.stepSize = 0
	lc.startPos = []int{0}
	return lc
}

func (lc *LineController) lineStartPosition(line int) int {
	index := lc.startPos[line]
	if line > lc.stepIndex {
		index += lc.stepSize
	}
	return index
}

func (lc *LineController) textInserted(line int, delta int) {
	if lc.stepSize == 0 {
		lc.stepIndex = line
		lc.stepSize = delta
	} else {
		if line != lc.stepIndex {
			lc.setStepIndex(line);
		}
		lc.stepSize += delta
	}
}

func (lc *LineController) setStepIndex(line int) {
	if line == lc.stepIndex {
		return
	}
	if lc.stepSize == 0 {
		return
	}
	if line > lc.stepIndex {
		for line > lc.stepIndex {
			lc.stepIndex++
			lc.startPos[lc.stepIndex] += lc.stepSize
		}
	} else {
		for line < lc.stepIndex {
			lc.startPos[lc.stepIndex] -= lc.stepSize
			lc.stepIndex--
		}
	}
}