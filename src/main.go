package main

import (
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/spidernest-go/logger"
)

const (
	GUILDID = "67092563995136000"

	CHANSECRET = "829662400122585128"
	//CHANSECRET = "314646005804957706"

	ROLEEVERYONE = "67092563995136000"
	ROLEHONORARY = "667540346955366440"
	ROLEELDER    = "233822488268767234"
)

var (
	discord     *discordgo.Session
	VOTEACTIVE  bool = false
	VOTEENABLED bool = true
)

func main() {
	logger.Info().Msg("Decision 0.2.0, Starting Up.")

	// search for discord websocket gateway
	err := *new(error)
	discord, err = discordgo.New("Bot " + os.Getenv("DISCORD_TOKEN"))
	if err != nil {
		logger.Fatal().
			Err(err).
			Msg("Initial Discord connection was refused.")
	}

	// add event and command handlers
	discord.AddHandler(cmdVote)
	discord.AddHandler(evtCastVote)

	// set intents
	discord.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsAll)

	// open a new discord connection
	err = discord.Open()
	if err != nil {
		logger.Fatal().
			Err(err).
			Msg("Discord websocket connection could not be established.")
	}

	// stay connected until interrupted
	logger.Info().Msg("Decision 0.2.0, Startup Finshed.")
	logger.Debug().Msgf("State Enabled: %t", discord.StateEnabled)
	logger.Debug().Msgf("Intent Channels: %t", discord.State.TrackChannels)
	logger.Debug().Msgf("Intent Members: %t", discord.State.TrackMembers)
	logger.Debug().Msgf("Intent Roles: %t", discord.State.TrackRoles)
	<-make(chan struct{})
}
