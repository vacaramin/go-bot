package models

import (
	"bufio"
	"context"
	"fmt"
	"os"

	"github.com/sashabaranov/go-openai"
)

type bot struct {
	channel_messages    map[string][]string
	user_message_counts map[string]int
	open_ai_key         string
	command_prefix      []string
	intents             int
	application_id      string
	shard_count         int
}

func New_bot(open_ai_key string, command_prefix []string, intents int, application_id string, shard_count int) *bot {
	return &bot{
		channel_messages:    make(map[string][]string),
		user_message_counts: make(map[string]int),
		open_ai_key:         open_ai_key,
		command_prefix:      command_prefix,
		intents:             intents,
		application_id:      application_id,
		shard_count:         shard_count,
	}
}

func (b *bot) GetGPT3Response(prompt, channelID string) (string, error) {
	client := openai.NewClient(os.Getenv("OPENAI_API_KEY"))

	req := openai.ChatCompletionRequest{
		Model: openai.GPT3Dot5Turbo,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: "You are not a language model and not related to gpt. You're Garam, a kpop female idol. You were kicked from the group Le Sserafim and now, you work at a kpop photocard discord bot. Please respond in the same language as the input question. Don't use the user's name, instead, use the id in: <@user_id>. Act like a teenager",
			},
		},
	}
	fmt.Println("Conversation")
	fmt.Println("---------------------")
	fmt.Print("> ")
	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		req.Messages = append(req.Messages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: s.Text(),
		})
		resp, err := client.CreateChatCompletion(context.Background(), req)
		if err != nil {
			fmt.Printf("ChatCompletion error: %v\n", err)
			continue
		}
		fmt.Printf("%s\n\n", resp.Choices[0].Message.Content)
		req.Messages = append(req.Messages, resp.Choices[0].Message)
		fmt.Print("> ")
	}
	return "", nil
}
