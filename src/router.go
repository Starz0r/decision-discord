package main

import (
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/spidernest-go/logger"
)

func command(context *discordgo.MessageCreate, cmd string, roleID string, private bool) bool {
	// check if the caller is it's self, or another bot
	if !context.Author.Bot {
		if strings.HasPrefix(context.Content, strings.Join([]string{"!", cmd}, "")) {

			//TODO: check if private here

			//BUG: This crashes Go for some inexplicable reason
			logger.Debug().Msg("Getting Channel from State.")
			channel, err := discord.State.Channel(context.ChannelID)
			if err != nil {
				logger.Error().
					Err(err).
					Msg("Could not retrieve the channel.")

				return false
			}
			logger.Debug().Msg("Finished, Got Channel.")

			member, err := discord.State.Member(channel.GuildID, context.Author.ID)
			//member, err := discord.State.Member("67092563995136000", context.Author.ID)
			if err != nil {
				logger.Error().
					Err(err).
					Msg("Could not retrieve the member.")

				return false
			}

			for i := range member.Roles {
				switch roleID {
				case ROLEEVERYONE:
					if member.Roles[i] == ROLEEVERYONE {
						return true
					}
					fallthrough
				case ROLEHONORARY:
					if member.Roles[i] == ROLEHONORARY {
						return true
					}
					fallthrough
				case ROLEELDER:
					if member.Roles[i] == ROLEELDER {
						return true
					}
					break
				}
			}

			logger.Info().
				Msg(context.Author.Username + " did not meet role requirement for command [" + cmd + "]")

		}
	}

	return false
}
