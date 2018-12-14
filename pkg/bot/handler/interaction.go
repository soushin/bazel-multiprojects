package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"google.golang.org/grpc/metadata"

	"go.uber.org/zap"

	"github.com/nlopes/slack"
	"github.com/pkg/errors"
	"github.com/soushin/bazel-multiprojects/proto/ops"
)

type interactionHandler struct {
	appLog       *zap.Logger
	slackCli     *slack.Client
	opsDeployCli ops.DeployClient
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

		selected := strings.Split(message.Actions[0].SelectedOptions[0].Value, ",")
		fullName := selected[0]
		packageName := selected[1]

		state := DialogState{
			ResponseURL:    message.ResponseURL,
			SubmissionType: SubmissionInputTag,
			Values: map[string]string{
				"fullName": fullName,
				"package":  packageName,
			},
		}

		stateStr, err := json.Marshal(state)
		if err != nil {
			h.appLog.With(zap.Any("state", state)).Error("invalid process")
			return nil, errors.Wrapf(err, "failed to marshal state")
		}

		md := metadata.New(map[string]string{})
		ctx := metadata.NewOutgoingContext(context.Background(), md)
		splitName := strings.Split(fullName, "/")
		res, err := h.opsDeployCli.GetBranches(ctx, &ops.BranchInbound{
			Owner:      splitName[0],
			Repository: splitName[1],
		})
		if err != nil {
			return nil, errors.Wrapf(err, "failed to get branches")
		}

		branches := make([]slack.DialogSelectOption, len(res.Branches))
		for i, branch := range res.Branches {
			branches[i] = slack.DialogSelectOption{
				Label: branch.Name,
				Value: branch.Name,
			}
		}

		dialog := slack.Dialog{
			TriggerID:  message.TriggerID,
			CallbackID: message.CallbackID,
			Title:      "ブランチを入力してね",
			State:      string(stateStr),
			Elements: []slack.DialogElement{
				slack.NewStaticSelectDialogInput(
					"branch",
					"ブランチ",
					branches,
				),
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
