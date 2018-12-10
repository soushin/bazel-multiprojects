package handler

import (
	"encoding/json"
	"fmt"

	"go.uber.org/zap"

	"github.com/nlopes/slack"
	"github.com/pkg/errors"
	"github.com/soushin/bazel-multiprojects/pkg/bot/client"
)

type submissionHandler struct {
	appLog      *zap.Logger
	slackExtCli client.SlackExt
}

const (
	SubmissionInputTag = "input_tag"
)

func (h submissionHandler) Handle(message slack.InteractionCallback) error {

	var state DialogState
	err := json.Unmarshal([]byte(message.State), &state)
	if err != nil {
		h.appLog.With(zap.Any("state", message.State)).Error("invalid process")
		return errors.Wrapf(err, "failed to unmarshal state")
	}

	switch state.SubmissionType {
	case SubmissionInputTag:
		attachments := []slack.Attachment{
			{
				Title: fmt.Sprintf(":ok: %sをデプロイするね", message.Submission["tag"]),
			},
		}
		response := client.RespondPayload{
			ResponseType: "in_channel",
			Attachments:  attachments,
		}
		if _, err := h.slackExtCli.Respond(state.ResponseURL, response); err != nil {
			h.appLog.With(zap.Error(err)).Error("invalid process")
			return errors.Wrapf(err, "respond to slack is failed")
		}
		return nil
	default:
		h.appLog.With(zap.String("submissionType", state.SubmissionType)).Error("invalid process")
		return errors.New("invalid submissionType was submitted")
	}

	return nil
}
