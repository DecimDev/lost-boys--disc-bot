package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

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

	botToken := os.Getenv("DISCORD_BOT_TOKEN")
	if botToken == "" {
		log.Fatal("DISCORD_BOT_TOKEN is not set in .env file")
	}

	client, err := disgo.New(botToken,
		// set gateway options
		bot.WithGatewayConfigOpts(
			gateway.WithIntents(
				gateway.IntentGuildPresences, // Listen for guild message events
				gateway.IntentMessageContent, // Listen for message content (for reading message text)
				gateway.IntentGuildMembers,
				gateway.IntentGuildMessages, // Listen for reactions on messages
			),
		),
		bot.WithEventListenerFunc(onMessageCreate),
	)
	if err != nil {
		log.Println("Error while building disgo client", err)
		return
	}

	// Ensure the client closes gracefully
	defer client.Close(context.TODO())

	// Connect to the Discord gateway
	if err = client.OpenGateway(context.TODO()); err != nil {
		log.Println("Error connecting to Discord gateway", err)
		return
	}

	// Bot is now running
	log.Println("lost-boys-bot-go is now running. Press CTRL-C to exit.")

	// Listen for termination signals to shut down gracefully
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-s
}

func onMessageCreate(event *events.MessageCreate) {
	if event.Message.Author.Bot {
		return
	}

	var message string
	if event.Message.Content == "hows that jonze?" {
		message = "yeah fuck you jonze"
	} else if event.Message.Content == "pong" {
		message = "ping"
	}

	if message != "" {
		_, err := event.Client().Rest().CreateMessage(event.ChannelID, discord.NewMessageCreateBuilder().SetContent(message).Build())
		if err != nil {
			log.Printf("Error sending message: %v", err)
		}
	}
}
