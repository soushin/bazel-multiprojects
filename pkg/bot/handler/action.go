package handler

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"

	"google.golang.org/grpc/metadata"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/nlopes/slack"
	"github.com/pkg/errors"
	"github.com/soushin/bazel-multiprojects/proto/ops"
)

type actionHandler struct {
	appLog       *zap.Logger
	slackCli     *slack.Client
	opsDeployCli ops.DeployClient
}

const (
	ActionDeploy = "deploy"
)

func NewActionHandler(appLog *zap.Logger, slackCli *slack.Client, opsDeployCli ops.DeployClient) *actionHandler {
	return &actionHandler{
		appLog:       appLog,
		slackCli:     slackCli,
		opsDeployCli: opsDeployCli,
	}
}

func (h actionHandler) Handle(action string, ev *slack.MessageEvent) error {

	switch action {
	case ActionDeploy:

		md := metadata.New(map[string]string{})
		ctx := metadata.NewOutgoingContext(context.Background(), md)

		res, err := h.opsDeployCli.GetTargets(ctx, &empty.Empty{})
		if err != nil {
			return errors.Wrapf(err, "failed to get targets")
		}

		h.appLog.With(zap.Any("res", res)).Info("Debug")

		options := make([]slack.AttachmentActionOption, len(res.Targets))
		for i, target := range res.Targets {
			options[i] = slack.AttachmentActionOption{
				Text:  target.Pkg,
				Value: fmt.Sprintf("%s,%s", target.Repo, target.Pkg),
			}
		}

		attachment := slack.Attachment{
			Text:       "どのパッケージをデプロイするの？",
			Color:      "#f9a41b",
			CallbackID: fmt.Sprintf("%d", time.Now().Unix()),
			Actions: []slack.AttachmentAction{
				{
					Name:    InteractionSelect,
					Type:    "select",
					Options: options,
				},
				{
					Name:  InteractionCancel,
					Text:  "Cancel",
					Type:  "button",
					Style: "danger",
				},
			},
		}

		if _, _, err := h.slackCli.PostMessage(ev.Channel, slack.MsgOptionAttachments(attachment)); err != nil {
			return fmt.Errorf("failed to post message: %s", err)
		}
		return nil
	default:
		if _, _, err := h.slackCli.PostMessage(ev.Channel, slack.MsgOptionText("そのオペレーションは知らいないなぁ:thinking_face:", false)); err != nil {
			return fmt.Errorf("failed to post message: %s", err)
		}
		return nil
	}
}
