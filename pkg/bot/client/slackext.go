package client

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"go.uber.org/zap"

	"github.com/nlopes/slack"
	"github.com/pkg/errors"
)

type SlackExt interface {
	Respond(responseURL string, payload RespondPayload) (*RespondResponse, error)
}

type RespondPayload struct {
	Text         string             `json:"text"`
	ResponseType string             `json:"response_type"`
	Attachments  []slack.Attachment `json:"attachments"`
}

type RespondResponse struct {
	OK bool `json:"ok"`
}

type clientImpl struct {
	appLog  *zap.Logger
	httpCli *http.Client
}

func NewSlackExt(appLog *zap.Logger, httpCli *http.Client) SlackExt {
	return &clientImpl{
		appLog:  appLog,
		httpCli: httpCli,
	}
}

func (c *clientImpl) Respond(responseURL string, payload RespondPayload) (*RespondResponse, error) {

	payloadBytes, err := json.Marshal(&payload)
	if err != nil {
		return nil, errors.Wrapf(err, "Marshal json error. [%+v]", payload)
	}

	req, err := http.NewRequest("POST", responseURL, bytes.NewReader(payloadBytes))
	if err != nil {
		return nil, errors.Wrapf(err, "create http.Request is failed.")
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpCli.Do(req)
	if err != nil {
		return nil, errors.Wrapf(err, "send http request is failed. request=%+v", req)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			c.appLog.With(zap.Error(err)).Error("Failed to callback respond")
		}
	}()

	res, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrapf(err, "read response is failed. response=%+v", resp.Body)
	}

	var response RespondResponse
	if err := json.Unmarshal(res, &response); err != nil {
		return nil, errors.Wrapf(err, "json unmarshal is failed. data=%+v", string(res))
	}

	return &response, nil
}
