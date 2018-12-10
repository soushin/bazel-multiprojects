package listener

import (
	"fmt"
	"strings"

	"go.uber.org/zap"

	"github.com/nlopes/slack"
	"github.com/pkg/errors"
	"github.com/soushin/bazel-multiprojects/pkg/bot/handler"
	"github.com/soushin/bazel-multiprojects/proto/ops"
)

type slackListener struct {
	appLog       *zap.Logger
	slackCli     *slack.Client
	opsDeployCli ops.DeployClient
	botID        string
	channelID    string
}

func NewSlackListener(appLog *zap.Logger, slackCli *slack.Client, opsDeployCli ops.DeployClient, botID string, channelID string) *slackListener {
	return &slackListener{
		appLog:       appLog,
		slackCli:     slackCli,
		opsDeployCli: opsDeployCli,
		botID:        botID,
		channelID:    channelID,
	}
}

func (s *slackListener) ListenAndResponse() {
	rtm := s.slackCli.NewRTM()
	go rtm.ManageConnection()

	for msg := range rtm.IncomingEvents {
		switch ev := msg.Data.(type) {
		case *slack.MessageEvent:
			if err := s.handleMessageEvent(ev); err != nil {
				s.appLog.With(
					zap.Error(err),
				).Error("Failed to handle message")
			} else {
				s.appLog.With(
					zap.String("message", fmt.Sprintf("%+v", ev)),
				).Info("Received incoming event")
			}
		}
	}
}

func (s *slackListener) handleMessageEvent(ev *slack.MessageEvent) error {

	if ev.BotID != "" {
		s.appLog.With(
			zap.String("channel", ev.Channel),
			zap.String("botId", ev.BotID),
			zap.String("message", ev.Msg.Text),
		).Info("skip bot message.")
		return nil
	}

	if ev.Channel != s.channelID {
		s.appLog.With(
			zap.String("channel", ev.Channel),
			zap.String("botId", ev.BotID),
			zap.String("message", ev.Msg.Text),
		).Info("Not allowed operation")
		return nil
	}

	if !strings.HasPrefix(ev.Msg.Text, fmt.Sprintf("<@%s> ", s.botID)) {
		return nil
	}

	m := strings.Split(strings.TrimSpace(ev.Msg.Text), " ")[1:]
	if len(m) == 0 {
		return errors.New(fmt.Sprintf("failed to post message: %s", ev.Msg.Text))
	}

	action := handler.NewActionHandler(s.appLog, s.slackCli, s.opsDeployCli)

	return action.Handle(m[0], ev)
}
