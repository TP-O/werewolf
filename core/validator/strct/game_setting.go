package strct

// import (
// 	"math"
// 	"uwwolf/game/types"
// 	"uwwolf/util"
// 	"uwwolf/validator/tag"

// 	"github.com/go-playground/validator/v10"
// 	"golang.org/x/exp/maps"
// 	"golang.org/x/exp/slices"
// )

// func validateGameCapacity(sl validator.StructLevel, setting *types.GameSetting) {
// 	if len(setting.PlayerIDs) < int(util.Config().Game.MinCapacity) ||
// 		len(setting.PlayerIDs) > int(util.Config().Game.MaxCapacity) {
// 		sl.ReportError(
// 			setting.PlayerIDs,
// 			"player_ids",
// 			"PlayerIDs",
// 			tag.GameCapacityTag,
// 			"",
// 		)
// 	}
// }

// func validateTurnDuration(sl validator.StructLevel, setting *types.GameSetting) {
// 	if setting.TurnDuration < util.Config().Game.MinTurnDuration ||
// 		setting.TurnDuration > util.Config().Game.MaxTurnDuration {
// 		sl.ReportError(
// 			setting.TurnDuration,
// 			"turn_duration",
// 			"TurnDuration",
// 			tag.TurnDurationTag,
// 			"",
// 		)
// 	}
// }

// func validateDiscussionDuration(sl validator.StructLevel, setting *types.GameSetting) {
// 	if setting.DiscussionDuration < util.Config().Game.MinDiscussionDuration ||
// 		setting.DiscussionDuration > util.Config().Game.MaxDiscussionDuration {
// 		sl.ReportError(
// 			setting.DiscussionDuration,
// 			"discussion_duration",
// 			"DiscussionDuration",
// 			tag.DiscussionDurationTag,
// 			"",
// 		)
// 	}
// }

// func validateRoleID(sl validator.StructLevel, setting *types.GameSetting) {
// 	for _, roleID := range setting.RoleIDs {
// 		if !slices.Contains(maps.Keys(types.RoleIDSets), roleID) {
// 			sl.ReportError(
// 				setting.RoleIDs,
// 				"role_ids",
// 				"RoleIDs",
// 				tag.RoleIDTag,
// 				"",
// 			)

// 			break
// 		}
// 	}

// 	for _, roleID := range setting.RequiredRoleIDs {
// 		if !slices.Contains(setting.RoleIDs, roleID) {
// 			sl.ReportError(
// 				setting.RequiredRoleIDs,
// 				"required_role_ids",
// 				"RequiredRoleIDs",
// 				tag.RoleIDTag,
// 				"",
// 			)

// 			break
// 		}
// 	}
// }

// func validateNumberWerewolves(sl validator.StructLevel, setting *types.GameSetting) {
// 	capacity := len(setting.PlayerIDs)
// 	maxNumberWerewolves := uint8(math.Round(float64(capacity)/2)) - 2

// 	if capacity < 5 && setting.NumberWerewolves != 1 ||
// 		(capacity > 4 && setting.NumberWerewolves > maxNumberWerewolves) {
// 		sl.ReportError(
// 			setting.NumberWerewolves,
// 			"number_werewolves",
// 			"NumberWerewolves",
// 			tag.NumberWerewolvesTag,
// 			"",
// 		)
// 	}
// }

// func ValidateGameSetting(sl validator.StructLevel) {
// 	setting := sl.Current().Interface().(types.GameSetting)
// 	validateGameCapacity(sl, &setting)
// 	validateTurnDuration(sl, &setting)
// 	validateDiscussionDuration(sl, &setting)
// 	validateRoleID(sl, &setting)
// 	validateNumberWerewolves(sl, &setting)
// }
