package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/labstack/gommon/log"
	"github.com/nlopes/slack"
	"net/http"
	"net/url"
	"os"
)

const (
	UrlVerificationEvent  = "url_verification"
	ChannelCreatedEvent   = "channel_created"
	ChannelDeletedEvent   = "channel_deleted"
	ChannelRenameEvent    = "channel_rename"
	ChannelArchiveEvent   = "channel_archive"
	ChannelUnarchiveEvent = "channel_unarchive"
	EmojiChangedEvent     = "emoji_changed"
	SubTeamCreatedEvent   = "subteam_created"
	SubTeamUpdatedEvent   = "subteam_updated"
)

const (
	SlackIcon = ":bb8-flame:"
	SlackName = "BB-8"
)

type ApiEvent struct {
	Type      string `json:"type"`
	Challenge string `json:"challenge"`
	Token     string `json:"token"`
	Event     struct {
		Type string `json:"type"`
	}
}

type ChannelEvent struct {
	Type       string     `json:"type"`
	Challenge  string     `json:"challenge"`
	Token      string     `json:"token"`
	SlackEvent SlackEvent `json:"event"`
}

type SlackEvent struct {
	Type    string       `json:"type"`
	Channel SlackChannel `json:"channel"`
}

type SlackChannel struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type EmojiEvent struct {
	Type      string `json:"type"`
	Challenge string `json:"challenge"`
	Token     string `json:"token"`
	Event     struct {
		Type string `json:"type"`
		Name string `json:"name"`
	}
}

func main() {
	lambda.Start(eventApiHandler)
}

func eventApiHandler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	params, _ := url.ParseQuery(req.Body)
	// Print out server request logs of slack
	log.Infof("params:%v type:%T", params, params)

	// Get env variables
	botToken := os.Getenv("BOT_TOKEN")
	channelID := os.Getenv("CHANNEL_ID")
	signingSecrets := os.Getenv("SIGNING_SECRETS")

	res := events.APIGatewayProxyResponse{}
	apiEvent := &ApiEvent{}
	for key, _ := range params {
		err := json.Unmarshal([]byte(key), apiEvent)
		if err != nil {
			log.Error(err)
			return res, err
		}
	}

	slackClient := slack.New(botToken)

	if err := verify(signingSecrets, req); err != nil {
		log.Error(err)
		return res, err
	}

	res.Headers = make(map[string]string)
	res.Headers["Content-Type"] = "text/plain"
	res.StatusCode = http.StatusOK

	if apiEvent.Type == UrlVerificationEvent {
		res.Body = apiEvent.Challenge
		return res, nil
	}

	var msg string
	msgParams := slack.PostMessageParameters{
		Username:  SlackName,
		IconEmoji: SlackIcon,
		LinkNames: 1,
	}

	log.Infof("Type is %s", apiEvent.Event.Type)
	switch apiEvent.Event.Type {
	case ChannelCreatedEvent:
		log.Info("ChannelCreatedEvent")
		channelEvent := &ChannelEvent{}
		for key, _ := range params {
			err := json.Unmarshal([]byte(key), channelEvent)
			if err != nil {
				log.Error(err)
				return res, err
			}
		}

		msg = fmt.Sprintf("New channel called #%s was created :tada:", channelEvent.SlackEvent.Channel.Name)
	case EmojiChangedEvent:
		log.Info("EmojiChangedEvent")
		emojiEvent := &EmojiEvent{}
		for key, _ := range params {
			err := json.Unmarshal([]byte(key), emojiEvent)
			if err != nil {
				log.Error(err)
				return res, err
			}
		}

		msg = fmt.Sprintf("Found a new Emoji, `:%s:` :%s:", emojiEvent.Event.Name, emojiEvent.Event.Name)

	default:
		log.Info("default")

		res.StatusCode = http.StatusOK
		return res, nil
	}

	msgOptText := slack.MsgOptionText(msg, true)
	msgOptParams := slack.MsgOptionPostMessageParameters(msgParams)

	if _, _, err := slackClient.PostMessage(channelID, msgOptText, msgOptParams); err != nil {
		return res, fmt.Errorf("メッセージ送信に失敗: %s", err)
	}

	res.StatusCode = http.StatusOK
	return res, nil
}

func verify(signingSecrets string, req events.APIGatewayProxyRequest) error {
	httpHeader := http.Header{}
	for key, value := range req.Headers {
		httpHeader.Set(key, value)
	}
	sv, err := slack.NewSecretsVerifier(httpHeader, signingSecrets)
	if err != nil {
		log.Error(err)
		return err
	}

	if _, err := sv.Write([]byte(req.Body)); err != nil {
		log.Error(err)
		return err
	}

	if err := sv.Ensure(); err != nil {
		log.Error("Invalid SIGNING_SECRETS")
		return err
	}
	return nil
}
