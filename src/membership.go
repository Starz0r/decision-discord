package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/spidernest-go/logger"
)

var REQVOTES = 7
var BALLOTCOUNT = make(map[string]int8)
var VOTEMSG *discordgo.Message
var VOTETYPE = "undefined"
var VOTEUSER *discordgo.Member

func cmdVote(s *discordgo.Session, j *discordgo.MessageCreate) {
	if command(j, "vote", ROLEELDER, false) {
		args := strings.Split(j.Content, " ")

		if VOTEACTIVE == true {
			s.ChannelMessageSend(j.ChannelID, "Another vote cannot be started while another vote is currently running.")

			return
		} else if VOTEENABLED == false {
			s.ChannelMessageSend(j.ChannelID, "A new vote cannot be started at this time.")

			return
		}

		// make sure the channel is the secret one
		channel, err := discord.State.Channel(j.ChannelID)
		if err != nil {
			logger.Error().
				Err(err).
				Msg("Could not retrieve the channel.")

			s.ChannelMessageSend(j.ChannelID, "An internal error occurred, please try again.")

			return
		}

		if channel.ID == CHANSECRET {
			// check if anyone was properly mentioned
			if len(j.Mentions) <= 0 {
				s.ChannelMessageSend(j.ChannelID, "Error: No target was mentioned.")

				return
			}

			// check if the user can be targetted
			target, err := discord.State.Member(GUILDID, j.Mentions[0].ID)
			if err != nil {
				logger.Error().
					Err(err).
					Msg("Could not retrieve the member.")

				s.ChannelMessageSend(j.ChannelID, "Error: An internal issue has occurred, please try again.")

				return
			}

			for _, role := range target.Roles {
				if role == ROLEHONORARY {
					s.ChannelMessageSend(j.ChannelID, "Error: Targetted user's role is too high and currently cannot be chosen.")
					return
				}

				if role == ROLEELDER {
					s.ChannelMessageSend(j.ChannelID, "Error: Targetted user's role is too high and currently cannot be chosen.")
					return
				}
			}

			// vote stuff
			VOTETYPE = "undefined"
			votecolor := 0xffffff
			targetname := "noone"
			switch args[1] {
			case "kick":
				VOTETYPE = "kick"
				votecolor = 0xff0000
				targetname = j.Mentions[0].Username + "#" + j.Mentions[0].Discriminator

			case "ban":
				VOTETYPE = "ban"
				votecolor = 0xff0000
				targetname = j.Mentions[0].Username + "#" + j.Mentions[0].Discriminator

			default:
				s.ChannelMessageSend(j.ChannelID, "Error: A vote type was unrecognized or unspecified.")

				return
			}

			//TODO: caculate how many votes are needed to pass, currently hardcoded
			msg := NewEmbed().
				SetTitle("🗳️ " + VOTETYPE + " " + targetname).
				SetDescription("@here\n\n👍YES  \\ 👎NO\n\nThe vote end if the count is over " + strconv.Itoa(REQVOTES) + " or 24 hours from now.").
				SetColor(votecolor).MessageEmbed

			VOTEMSG, err = s.ChannelMessageSendEmbed(j.ChannelID, msg)
			if err != nil {
				s.ChannelMessageSend(j.ChannelID, "Error: Discord refused the request, please try again.")
				return
			}

			s.MessageReactionAdd(j.ChannelID, VOTEMSG.ID, "👍")
			s.MessageReactionAdd(j.ChannelID, VOTEMSG.ID, "👎")
			VOTEUSER = target
			deactivateVoting()
		}
	}
}

func evtCastVote(s *discordgo.Session, j *discordgo.MessageReactionAdd) {
	fmt.Println(j.MessageReaction.ChannelID)
	fmt.Println(j.MessageReaction.MessageID)
	fmt.Println(j.MessageReaction.Emoji.Name)
	if j.MessageReaction.MessageID != VOTEMSG.ID {
		return
	}

	voters, err := s.MessageReactions(j.MessageReaction.ChannelID, j.MessageReaction.MessageID, j.MessageReaction.Emoji.Name, 100, "", "")
	// passing vote
	if (len(voters) >= REQVOTES-1) && (j.MessageReaction.Emoji.Name == "👍") {
		s.ChannelMessageSend(VOTEMSG.ChannelID, "✔️ Vote Passed...")

		switch VOTETYPE {
		case "ban":
			s.ChannelMessageSend(VOTEMSG.ChannelID, "🔨Banning User...")
			err = s.GuildBanCreate(GUILDID, VOTEUSER.User.ID, 0)

			if err != nil {
				s.ChannelMessageSend(VOTEMSG.ChannelID, "Error: User could not be banned, will try again next period.")
				return
			}

			s.ChannelMessageSend(VOTEMSG.ChannelID, "🔨User was banned!")
			s.ChannelMessageSend(VOTEMSG.ChannelID, "Shutting down for the next 24 hours.")
			panic("")
		}
	}
}

func deactivateVoting() {
	VOTEACTIVE = false
	VOTEENABLED = false
	go Dispatch(time.Hour*24*7, time.Minute, reactivateVoting, "reactivateVoting")
}

func reactivateVoting() {
	VOTEENABLED = true
}
