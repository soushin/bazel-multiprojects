package handler

import (
	"encoding/json"
	"fmt"

	"github.com/nlopes/slack"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type interactionHandler struct {
	appLog   *zap.Logger
	slackCli *slack.Client
}

const (
	InteractionSelect = "select"
	InteractionCancel = "cancel"
)

func (h interactionHandler) Handle(message slack.InteractionCallback) (*slack.Message, error) {
	action := message.Actions[0]
	switch action.Name {
	case InteractionSelect:
		originalMessage := message.OriginalMessage
		originalMessage.ReplaceOriginal = true

		state := DialogState{
			ResponseURL:    message.ResponseURL,
			SubmissionType: SubmissionInputTag,
		}

		stateStr, err := json.Marshal(state)
		if err != nil {
			h.appLog.With(zap.Any("state", state)).Error("invalid process")
			return nil, errors.Wrapf(err, "failed to marshal state")
		}

		dialog := slack.Dialog{
			TriggerID:  message.TriggerID,
			CallbackID: message.CallbackID,
			Title:      "タグを入力してね",
			State:      string(stateStr),
			Elements: []slack.DialogElement{
				slack.NewTextInput("tag", "タグ", "latest"),
			},
		}

		if err := h.slackCli.OpenDialog(message.TriggerID, dialog); err != nil {
			h.appLog.With(
				zap.String("triggerID", message.TriggerID),
				zap.Any("dialog", dialog),
			).Error("invalid process")
			return nil, errors.Wrapf(err, "failed to open dialog")
		}

		return &originalMessage, nil
	case InteractionCancel:
		title := fmt.Sprintf(":x: @%s canceled the request", message.User.Name)
		res := h.responseMessage(message.OriginalMessage, title, "")
		return &res, nil
	default:
		h.appLog.With(zap.String("action", action.Name)).Error("invalid process")
		return nil, errors.New("invalid action was submitted")
	}
}

func (h interactionHandler) responseMessage(original slack.Message, title, value string) slack.Message {
	original.Attachments[0].Actions = []slack.AttachmentAction{} // empty buttons
	original.Attachments[0].Fields = []slack.AttachmentField{
		{
			Title: title,
			Value: value,
			Short: false,
		},
	}
	original.ReplaceOriginal = true

	return original
}
