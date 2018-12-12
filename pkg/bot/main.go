package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/soushin/bazel-multiprojects/proto/ops"
	"google.golang.org/grpc"

	"go.uber.org/zap"

	"github.com/kelseyhightower/envconfig"
	"github.com/nlopes/slack"
	"github.com/soushin/bazel-multiprojects/pkg/bot/client"
	"github.com/soushin/bazel-multiprojects/pkg/bot/handler"
	"github.com/soushin/bazel-multiprojects/pkg/bot/listener"
)

type envConfig struct {
	Port              string `envconfig:"PORT" default:"3000"`
	BotToken          string `envconfig:"BOT_TOKEN" required:"true" default:"bot_token"`
	VerificationToken string `envconfig:"VERIFICATION_TOKEN" required:"true" default:"verification_token"`
	BotID             string `envconfig:"BOT_ID" required:"true" default:"bit_id"`
	ChannelID         string `envconfig:"CHANNEL_ID" required:"true" default:"channel_id"`

	OpsAddr string `envconfig:"OPS_ADDR" default:"localhost:50051"`
}

func main() {
	os.Exit(_main(os.Args[1:]))
}

func _main(args []string) int {
	appLog, _ := zap.NewProduction()
	var env envConfig
	if err := envconfig.Process("", &env); err != nil {
		appLog.With(zap.Error(err)).Error("Failed to process env var")
		return 1
	}

	defaultHttpCLI := &http.Client{
		Timeout: 10 * time.Second,
	}

	conn, err := grpc.Dial(
		env.OpsAddr,
		grpc.WithInsecure(),
	)
	if err != nil {
		appLog.With(zap.Error(err)).Error("Failed to connect ops server")
		return 1
	}
	defer conn.Close()

	opsDeployCli := ops.NewDeployClient(conn)

	appLog.Info("Start slack event listening")
	slackCli := slack.New(env.BotToken)
	slackExtCli := client.NewSlackExt(appLog, defaultHttpCLI)
	slackListener := listener.NewSlackListener(appLog, slackCli, opsDeployCli, env.BotID, env.ChannelID)
	go slackListener.ListenAndResponse()

	http.Handle("/interaction", handler.NewSlackHandler(appLog, slackCli, slackExtCli, env.VerificationToken))
	http.HandleFunc("/hc", hcHandler)

	appLog.With(zap.String("port", env.Port)).Info("Server listening")
	if err := http.ListenAndServe(":"+env.Port, nil); err != nil {
		appLog.With(zap.Error(err)).Error("Failed to serve bot")
		return 1
	}

	return 0
}

func hcHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "true")
}
