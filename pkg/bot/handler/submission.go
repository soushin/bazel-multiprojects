package handler

import (
	"encoding/json"
	"strings"

	"go.uber.org/zap"

	"github.com/nlopes/slack"
	"github.com/pkg/errors"
	"github.com/soushin/bazel-multiprojects/pkg/bot/client"
	"github.com/soushin/bazel-multiprojects/pkg/bot/repository"
	"github.com/soushin/bazel-multiprojects/proto/ops"
)

type submissionHandler struct {
	appLog           *zap.Logger
	slackExtCli      client.SlackExt
	deployRepository repository.DeployRepository
}

const (
	SubmissionInputTag = "input_tag"
)

func (h *submissionHandler) Handle(message slack.InteractionCallback) error {

	var state DialogState
	err := json.Unmarshal([]byte(message.State), &state)
	if err != nil {
		h.appLog.With(zap.Any("state", message.State)).Error("invalid process")
		return errors.Wrapf(err, "failed to unmarshal state")
	}

	switch state.SubmissionType {
	case SubmissionInputTag:

		fullName := state.Values["fullName"]
		targets := strings.Split(fullName, "/")

		owner := targets[0]
		repo := targets[1]
		path := state.Values["package"]
		branch := message.Submission["branch"]

		go func() {
			observer := func(result *ops.DeployOutbound) {
				switch result.Progress {
				case ops.DeployProgress_STARTED:
					h.respond(state.ResponseURL, result.Message)
				case ops.DeployProgress_RUNNING:
					h.respond(state.ResponseURL, result.Message)
				case ops.DeployProgress_SUCCESS:
					h.respond(state.ResponseURL, result.Message)
				case ops.DeployProgress_ERROR:
					h.respond(state.ResponseURL, result.Message)
				}
			}

			if err := h.deployRepository.Execute(owner, repo, path, branch, observer); err != nil {
				h.appLog.With(zap.Error(err)).Error("invalid process")
			}
		}()

		return nil
	default:
		h.appLog.With(zap.String("submissionType", state.SubmissionType)).Error("invalid process")
		return errors.New("invalid submissionType was submitted")
	}

	return nil
}

func (h *submissionHandler) respond(responseURL, message string) {
	attachments := []slack.Attachment{
		{
			Title: message,
		},
	}
	response := client.RespondPayload{
		ResponseType: "in_channel",
		Attachments:  attachments,
	}
	if _, err := h.slackExtCli.Respond(responseURL, response); err != nil {
		h.appLog.With(zap.Error(err)).Error("invalid process")
	}
}
