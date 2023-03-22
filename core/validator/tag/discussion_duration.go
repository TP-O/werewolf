package tag

// import (
// 	"strconv"
// 	"uwwolf/util"

// 	ut "github.com/go-playground/universal-translator"
// 	"github.com/go-playground/validator/v10"
// )

// const DiscussionDurationTag = "discussion_duration"

// func AddDiscussionDurationTag(ut ut.Translator) error {
// 	message := "{0} must be from " +
// 		strconv.Itoa(int(util.Config().Game.MinDiscussionDuration)) +
// 		" to " +
// 		strconv.Itoa(int(util.Config().Game.MaxDiscussionDuration)) +
// 		" seconds"

// 	return ut.Add(DiscussionDurationTag, message, true)
// }

// func RegisterDiscussionDurationMessage(ut ut.Translator, fe validator.FieldError) string {
// 	t, _ := ut.T(DiscussionDurationTag, fe.Field())

// 	return t
// }
