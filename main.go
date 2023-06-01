package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/disgoorg/log"
	"github.com/sashabaranov/go-openai"

	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	log.SetLevel(log.LevelDebug)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	client, err := disgo.New(os.Getenv("TOKEN"),
		bot.WithGatewayConfigOpts(
			gateway.WithIntents(
				gateway.IntentGuildMessages,
				gateway.IntentMessageContent,
			),
		),
		bot.WithEventListenerFunc(onMessageCreate),
	)
	//card = cogs.NewCards(&client)

	if err != nil {
		log.Fatal("error while building disgo: ", err)
	}

	defer client.Close(context.TODO())

	if err = client.OpenGateway(context.TODO()); err != nil {
		log.Fatal("errors while connecting to gateway: ", err)
	}

	log.Info("example is now running. Press CTRL-C to exit.")
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-s
}

func onMessageCreate(event *events.MessageCreate) {
	if event.Message.Author.Bot {
		return
	}

	mentioned := false
	for _, user := range event.Message.Mentions {
		if user.ID == event.Client().ID() {
			mentioned = true
			break
		}
	}
	if mentioned {
		client := openai.NewClient(os.Getenv("OPENAI_API_KEY"))

		req := openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: "You are not a language model and not related to gpt. You're Garam, a kpop female idol. You were kicked from the group Le Sserafim and now, you work at a kpop photocard discord bot. Please respond in the same language as the input question. Don't use the user's name, instead, use the id in: <@user_id>. Act like a teenager",
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: event.Message.Content,
				},
			},
		}

		resp, err := client.CreateChatCompletion(context.Background(), req)
		if err != nil {
			fmt.Printf("ChatCompletion error: %v\n", err)
			return
		}

		message := resp.Choices[0].Message.Content
		fmt.Println("message sent:", message)
		_, _ = event.Client().Rest().CreateMessage(event.ChannelID, discord.NewMessageCreateBuilder().SetContent(message).Build())
	}
}
