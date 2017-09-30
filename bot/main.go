package bot

import (
	"fmt"
	"net/http"
	"os"

	"golang.org/x/net/context"

	"github.com/line/line-bot-sdk-go/linebot"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/urlfetch"
)

/**
 * LINE Botクライアントインスタンスを生成
 */
func createBotClient(c context.Context, client *http.Client) (bot *linebot.Client, err error) {
	var (
		channelSecret = os.Getenv("LINEBOT_CHANNEL_SECRET")
		channelToken  = os.Getenv("LINEBOT_CHANNEL_ACCESS_TOKEN")
	)

	bot, err = linebot.New(channelSecret, channelToken, linebot.WithHTTPClient(client)) //Appengineのurlfetchを使用する
	if err != nil {
		log.Errorf(c, "Error occurred at create linebot client: %v", err)
		return bot, err
	}
	return bot, nil
}

/**
 * Get event sender's id
 */
func getSenderID(c context.Context, event *linebot.Event) string {
	switch event.Source.Type {
	case linebot.EventSourceTypeGroup:
		return event.Source.GroupID
	case linebot.EventSourceTypeRoom:
		return event.Source.RoomID
	case linebot.EventSourceTypeUser:
		return event.Source.UserID
	}
	log.Warningf(c, "Can not get sender id. type: %v", event.Source.Type)
	return ""
}

/**
 * 送信者の表示名を取得する
 *
 * ユーザしか取得できないので、ルームおよびグループではidをそのまま返す
 * グループメンバーのUserIDの場合、そのユーザが直接Botと友だち登録していなければ取得できない
 */
func getSenderName(c context.Context, bot *linebot.Client, from string) string {
	if len(from) == 0 {
		log.Warningf(c, "Parameter `mid` was not specified.")
		return from
	}
	if from[0:1] == "U" {
		senderProfile, err := bot.GetProfile(from).Do()
		if err != nil {
			log.Warningf(c, "Error occurred at get sender profile. from: %v, err: %v", from, err)
			return from
		}
		return senderProfile.DisplayName
	}
	return from
}

/**
 * LINE Messaging APIからのコールバックをハンドリング
 */
func lineBotCallback(w http.ResponseWriter, r *http.Request) {

	c := appengine.NewContext(r)
	bot, err := createBotClient(c, urlfetch.Client(c))
	if err != nil {
		return
	}

	events, err := bot.ParseRequest(r)
	if err != nil {
		if err == linebot.ErrInvalidSignature {
			log.Warningf(c, "Linebot request status: 400")
			w.WriteHeader(400)
		} else {
			log.Warningf(c, "linebot request status: 500\n\terror: %v", err)
			w.WriteHeader(500)
		}
		return
	}

	for _, event := range events {
		switch event.Type {
		case linebot.EventTypeFollow, linebot.EventTypeJoin:
			sender := getSenderName(c, bot, getSenderID(c, event))
			message := sender + " さん、友だち登録ありがとうございます！"
			if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(message)).Do(); err != nil {
				log.Errorf(c, "Error occurred at reply-message for follow/join. err: %v", err)
			}

		case linebot.EventTypeUnfollow, linebot.EventTypeLeave:
			sender := getSenderName(c, bot, getSenderID(c, event))
			message := sender + " さん、さようなら"
			if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(message)).Do(); err != nil {
				log.Errorf(c, "Error occurred at reply-message for unfollow/leave. err: %v", err)
			}

		case linebot.EventTypeMessage:
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				var replyMessage string
				if message.Text == "/version" {
					//バージョン
					replyMessage = "version: " + version

				} else if message.Text == "/mention" || message.Text == "/mention1" {
					//IDでメンション
					replyMessage = "@" + event.Source.UserID + " メンションになっていますか？"

				} else if message.Text == "/mention2" {
					//名前でメンション
					sender := getSenderName(c, bot, event.Source.UserID)
					replyMessage = "@" + sender + " メンションになっていますか？"

				} else {
					//オウム返し
					sender := getSenderName(c, bot, event.Source.UserID)
					replyMessage = sender + "さんの発言:\n" + message.Text
				}
				if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(replyMessage)).Do(); err != nil {
					log.Errorf(c, "Error occurred at reply-message. err: %v", err)
				}
			}

		default:
			log.Debugf(c, "Unsupported event type. type: %v", event.Type)
		}
	}
}

//
func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "version: %v", version)
}

//
func init() {
	//LINE Bot
	http.HandleFunc("/linebot/callback", lineBotCallback)

	//Web App
	http.HandleFunc("/", index)
}
