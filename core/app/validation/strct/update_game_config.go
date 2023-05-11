package strct

import (
	"strconv"
	"uwwolf/app/dto"
	"uwwolf/config"
	"uwwolf/game/vars"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

const TurnDurationTag = "turn_duration"

func AddTurnDurationTag(config config.Game) func(ut ut.Translator) error {
	return func(ut ut.Translator) error {
		message := "{0} must be from " +
			strconv.Itoa(int(config.MinTurnDuration.Seconds())) +
			" to " +
			strconv.Itoa(int(config.MaxTurnDuration.Seconds())) +
			" seconds"
		return ut.Add(TurnDurationTag, message, true)
	}
}

func RegisterTurnDurationMessage(ut ut.Translator, fe validator.FieldError) string {
	t, _ := ut.T(TurnDurationTag, fe.Field())
	return t
}

func validateTurnDuration(sl validator.StructLevel, config config.Game, gameCfg *dto.ReplaceGameConfigDto) {
	if gameCfg.TurnDuration < uint16(config.MinTurnDuration.Seconds()) ||
		gameCfg.TurnDuration > uint16(config.MaxTurnDuration.Seconds()) {
		sl.ReportError(
			gameCfg.TurnDuration,
			"turn_duration",
			"TurnDuration",
			TurnDurationTag,
			"",
		)
	}
}

const DiscussionDurationTag = "discussion_duration"

func AddDiscussionDurationTag(config config.Game) func(ut ut.Translator) error {
	return func(ut ut.Translator) error {
		message := "{0} must be from " +
			strconv.Itoa(int(config.MinDiscussionDuration.Seconds())) +
			" to " +
			strconv.Itoa(int(config.MaxDiscussionDuration.Seconds())) +
			" seconds"
		return ut.Add(DiscussionDurationTag, message, true)
	}
}

func RegisterDiscussionDurationMessage(ut ut.Translator, fe validator.FieldError) string {
	t, _ := ut.T(DiscussionDurationTag, fe.Field())
	return t
}

func validateDiscussionDuration(sl validator.StructLevel, config config.Game, gameCfg *dto.ReplaceGameConfigDto) {
	if gameCfg.DiscussionDuration < uint16(config.MinDiscussionDuration.Seconds()) ||
		gameCfg.DiscussionDuration > uint16(config.MaxDiscussionDuration.Seconds()) {
		sl.ReportError(
			gameCfg.DiscussionDuration,
			DiscussionDurationTag,
			"DiscussionDuration",
			DiscussionDurationTag,
			"",
		)
	}
}

const RoleIDTag = "role_ids"

func AddRoleIDsTag(ut ut.Translator) error {
	message := "{0} contains invalid role id"
	return ut.Add(RoleIDTag, message, true)
}

func RegisterRoleIDsMessage(ut ut.Translator, fe validator.FieldError) string {
	t, _ := ut.T(RoleIDTag, fe.Field())
	return t
}

func validateRoleIDs(sl validator.StructLevel, gameCfg *dto.ReplaceGameConfigDto) {
	for _, roleID := range gameCfg.RoleIDs {
		if !slices.Contains(maps.Keys(vars.RoleSets), roleID) {
			sl.ReportError(
				gameCfg.RoleIDs,
				RoleIDTag,
				"RoleIDs",
				RoleIDTag,
				"",
			)
			return
		}
	}
}

const RequiredRoleIDsTag = "required_role_ids"

func AddRequiredRoleIDsTag(ut ut.Translator) error {
	message := "{0} must be in role_ids"
	return ut.Add(RequiredRoleIDsTag, message, true)
}

func RegisterRequiredRoleIDsMessage(ut ut.Translator, fe validator.FieldError) string {
	t, _ := ut.T(RequiredRoleIDsTag, fe.Field())
	return t
}

func validateRequiredRoleIDs(sl validator.StructLevel, gameCfg *dto.ReplaceGameConfigDto) {
	for _, roleID := range gameCfg.RequiredRoleIDs {
		if !slices.Contains(gameCfg.RoleIDs, roleID) {
			sl.ReportError(
				gameCfg.RequiredRoleIDs,
				RequiredRoleIDsTag,
				"RequiredRoleIDs",
				RequiredRoleIDsTag,
				"",
			)
			return
		}
	}
}

func ValidateGameConfig(config config.Game) func(sl validator.StructLevel) {
	return func(sl validator.StructLevel) {
		gameCfg := sl.Current().Interface().(dto.ReplaceGameConfigDto)
		validateTurnDuration(sl, config, &gameCfg)
		validateDiscussionDuration(sl, config, &gameCfg)
		validateRoleIDs(sl, &gameCfg)
		validateRequiredRoleIDs(sl, &gameCfg)
	}
}
