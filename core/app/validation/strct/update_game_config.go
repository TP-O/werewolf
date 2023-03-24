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

func AddTurnDurationTag(ut ut.Translator) error {
	message := "{0} must be from " +
		strconv.Itoa(int(config.Game().MinTurnDuration.Seconds())) +
		" to " +
		strconv.Itoa(int(config.Game().MaxTurnDuration.Seconds())) +
		" seconds"
	return ut.Add(TurnDurationTag, message, true)
}

func RegisterTurnDurationMessage(ut ut.Translator, fe validator.FieldError) string {
	t, _ := ut.T(TurnDurationTag, fe.Field())
	return t
}

func validateTurnDuration(sl validator.StructLevel, cfg *dto.ReplaceGameConfigDto) {
	if cfg.TurnDuration < uint16(config.Game().MinTurnDuration.Seconds()) ||
		cfg.TurnDuration > uint16(config.Game().MaxTurnDuration.Seconds()) {
		sl.ReportError(
			cfg.TurnDuration,
			"turn_duration",
			"TurnDuration",
			TurnDurationTag,
			"",
		)
	}
}

const DiscussionDurationTag = "discussion_duration"

func AddDiscussionDurationTag(ut ut.Translator) error {
	message := "{0} must be from " +
		strconv.Itoa(int(config.Game().MinDiscussionDuration.Seconds())) +
		" to " +
		strconv.Itoa(int(config.Game().MaxDiscussionDuration.Seconds())) +
		" seconds"
	return ut.Add(DiscussionDurationTag, message, true)
}

func RegisterDiscussionDurationMessage(ut ut.Translator, fe validator.FieldError) string {
	t, _ := ut.T(DiscussionDurationTag, fe.Field())
	return t
}

func validateDiscussionDuration(sl validator.StructLevel, cfg *dto.ReplaceGameConfigDto) {
	if cfg.DiscussionDuration < uint16(config.Game().MinDiscussionDuration.Seconds()) ||
		cfg.DiscussionDuration > uint16(config.Game().MaxDiscussionDuration.Seconds()) {
		sl.ReportError(
			cfg.DiscussionDuration,
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

func validateRoleIDs(sl validator.StructLevel, cfg *dto.ReplaceGameConfigDto) {
	for _, roleID := range cfg.RoleIDs {
		if !slices.Contains(maps.Keys(vars.RoleSets), roleID) {
			sl.ReportError(
				cfg.RoleIDs,
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

func validateRequiredRoleIDs(sl validator.StructLevel, cfg *dto.ReplaceGameConfigDto) {
	for _, roleID := range cfg.RequiredRoleIDs {
		if !slices.Contains(cfg.RoleIDs, roleID) {
			sl.ReportError(
				cfg.RequiredRoleIDs,
				RequiredRoleIDsTag,
				"RequiredRoleIDs",
				RequiredRoleIDsTag,
				"",
			)
			return
		}
	}
}

func ValidateGameConfig(sl validator.StructLevel) {
	cfg := sl.Current().Interface().(dto.ReplaceGameConfigDto)
	validateTurnDuration(sl, &cfg)
	validateDiscussionDuration(sl, &cfg)
	validateRoleIDs(sl, &cfg)
	validateRequiredRoleIDs(sl, &cfg)
}
