package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/nlopes/slack"
	"github.com/soushin/bazel-multiprojects/pkg/bot/client"

	"go.uber.org/zap"
)

type slackHandler struct {
	appLog            *zap.Logger
	slackCli          *slack.Client
	slackExtCli       client.SlackExt
	verificationToken string
}

func NewSlackHandler(appLog *zap.Logger, slackCli *slack.Client, slackExtCli client.SlackExt,
	verificationToken string) *slackHandler {
	return &slackHandler{
		appLog:            appLog,
		slackCli:          slackCli,
		slackExtCli:       slackExtCli,
		verificationToken: verificationToken,
	}
}

func (h slackHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.appLog.With(zap.String("method", r.Method)).Error("Invalid method")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		h.appLog.With(zap.Error(err)).Error("Failed to read request body")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonStr, err := url.QueryUnescape(string(buf)[8:])
	if err != nil {
		h.appLog.With(zap.Error(err)).Error("Failed to unespace request body")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	h.appLog.With(zap.String("jsonStr", jsonStr)).Info("Debug")

	var message slack.InteractionCallback
	if err := json.Unmarshal([]byte(jsonStr), &message); err != nil {
		h.appLog.With(zap.String("json", jsonStr)).Error("Failed to decode json message from slack")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if message.Token != h.verificationToken {
		h.appLog.With(zap.String("token", message.Token)).Error("Invalid token")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if strings.Contains(jsonStr, `"type":"dialog_submission"`) {

		submission := submissionHandler{
			appLog:      h.appLog,
			slackExtCli: h.slackExtCli,
		}

		if err := submission.Handle(message); err != nil {
			h.appLog.With(zap.Error(err)).Error("Failed to handle submission")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		return
	} else {

		interaction := interactionHandler{
			appLog:   h.appLog,
			slackCli: h.slackCli,
		}

		originalMessage, err := interaction.Handle(message)
		if err != nil {
			h.appLog.With(zap.Error(err)).Error("Failed to handle interactive")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Add("Content-type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(&originalMessage)
		return
	}
}
