package packetslogic

import (
	"fmt"
	"little/packets"
	"little/types"
)

func SplitFromToOnSteps(cleanFromX, cleanToX, cleanFromY, cleanToY uint32) (steps []packets.ActionPosition) {
	diffX := int(cleanToX) - int(cleanFromX)
	diffY := int(cleanToY) - int(cleanFromY)

	for {
		nextStep := packets.ActionPosition{}
		if diffX != 0 {
			if diffX < 0 {
				cleanFromX = uint32(int(cleanFromX) - 1)
				nextStep.X.Value = cleanFromX
				diffX += 1
			}

			if diffX > 0 {
				cleanFromX = uint32(int(cleanFromX) + 1)
				nextStep.X.Value = cleanFromX
				diffX -= 1
			}
		} else {
			nextStep.X.Value = cleanToX
		}

		if diffY != 0 {
			if diffY < 0 {
				cleanFromY = uint32(int(cleanFromY) - 1)
				nextStep.Y.Value = cleanFromY
				diffY += 1
			}

			if diffY > 0 {
				cleanFromY = uint32(int(cleanFromY) + 1)
				nextStep.Y.Value = cleanFromY
				diffY -= 1
			}
		} else {
			nextStep.Y.Value = cleanToY
		}

		if diffX == 0 && diffY == 0 {
			break
		}

		// Normalize
		nextStep.X.Value = nextStep.X.Value * 100
		nextStep.Y.Value = nextStep.Y.Value * 100

		steps = append(steps, nextStep)

		fmt.Println("Next step:", nextStep)
	}

	// Last step
	steps = append(steps, packets.ActionPosition{
		X: types.HardUInt32{Value: cleanToX * 100},
		Y: types.HardUInt32{Value: cleanToY * 100},
	})

	return
}
