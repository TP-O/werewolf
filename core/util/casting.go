package util

import "uwwolf/game/types"

func CastPlayerIDs2Strings(pIDs []types.PlayerID) []string {
	var castedArr []string
	for _, pID := range pIDs {
		castedArr = append(castedArr, string(pID))
	}

	return castedArr
}
