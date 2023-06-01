package cogs

import (
	"math/rand"
	"time"

	"github.com/disgoorg/disgo/bot"
)

type Cards struct {
	bot *bot.Client
}

func NewCards(bot *bot.Client) *Cards {
	return &Cards{
		bot: bot,
	}
}

func (c *Cards) RandomizeRarity() int {
	randomRarity := rand.Intn(1000)
	rarity := 0
	if randomRarity < 1 {
		rarity = 3
	} else if randomRarity < 90 {
		rarity = 2
	} else if randomRarity < 1000 {
		rarity = 1
	}
	return rarity
}

// Helper function to generate a random hexadecimal string of given length
func RandomHex(length int) string {
	const hexChars = "0123456789ABCDEF"
	source := rand.NewSource(time.Now().UnixNano())
	random := rand.New(source)
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		result[i] = hexChars[random.Intn(len(hexChars))]
	}
	return string(result)
}
