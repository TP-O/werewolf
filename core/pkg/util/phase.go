package util

import (
	"uwwolf/internal/app/game/logic/constants"
	"uwwolf/internal/app/game/logic/types"
)

// NextPhasePhaseID returns the next phase ID based on the given lastPhaseID.
func NextPhasePhaseID(currentPhaseId types.PhaseId) types.PhaseId {
	if currentPhaseId == constants.DuskPhaseId {
		return 1
	}
	return currentPhaseId + 1
}

// PreviousPhasePhaseID returns the previous phase ID based on the given lastPhaseID.
func PreviousPhaseID(currentPhaseId types.PhaseId) types.PhaseId {
	if currentPhaseId == 1 {
		return constants.DuskPhaseId
	}
	return currentPhaseId - 1
}
