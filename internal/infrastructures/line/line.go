package line

import (
	"net/http"

	"github.com/line/line-bot-sdk-go/v7/linebot"
	"github.com/pkg/errors"
)

type (
	Env struct {
		ChannelSecret       string `envconfig:"channel_secret" default:"mock-secret"`
		ChannelToken        string `envconfig:"channel_token" default:"mock-token"`
		APIEndpointBase     string `envconfig:"api_endpoint_base" default:"https://api.line.me"`
		APIEndpointBaseData string `envconfig:"api_endpoint_base_data" default:"https://api-data.line.me"`
	}
)

func NewClient(env Env, httpCli *http.Client) (*linebot.Client, error) {
	cli, err := linebot.New(
		env.ChannelSecret,
		env.ChannelToken,
		linebot.WithHTTPClient(httpCli),
		linebot.WithEndpointBase(env.APIEndpointBase),
		linebot.WithEndpointBaseData(env.APIEndpointBaseData),
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create linebot client")
	}

	return cli, nil
}
