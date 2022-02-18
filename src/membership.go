package main

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/spidernest-go/logger"
)

var REQVOTES = 4
var BALLOTCOUNT = make(map[string]int8)
var BALLOTMUTEX sync.RWMutex
var VOTEMSG *discordgo.Message
var VOTETYPE = "undefined"
var VOTEUSER *discordgo.Member
var VOTEUSERSTR string
var VOTEEMBED *discordgo.MessageEmbed
var VOTEDESCRIPT string

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
			VOTEUSER = target

			// vote stuff
			VOTETYPE = "undefined"
			votecolor := 0xffffff
			targetname := "noone"
			VOTEDESCRIPT = "No Descrption Given."
			switch args[1] {
			case "kick":
				VOTETYPE = "kick"
				votecolor = 0xff0000
				targetname = j.Mentions[0].Username + "#" + j.Mentions[0].Discriminator
				VOTEDESCRIPT = "*KICK* User " + targetname + " from the server?"

			case "ban":
				VOTETYPE = "ban"
				votecolor = 0xff0000
				targetname = j.Mentions[0].Username + "#" + j.Mentions[0].Discriminator
				VOTEDESCRIPT = "*BAN* User " + targetname + " from the server?"

			case "promote":
				VOTETYPE = "promote"
				votecolor = 0x00ff10
				targetname = j.Mentions[0].Username + "#" + j.Mentions[0].Discriminator
				VOTEDESCRIPT = "*PROMOTE* User " + targetname + " to the next rank?"

			default:
				s.ChannelMessageSend(j.ChannelID, "Error: A vote type was unrecognized or unspecified.")

				return
			}
			VOTEUSERSTR = targetname

			//TODO: caculate how many votes are needed to pass, currently hardcoded
			emb := NewEmbed().
				SetTitle("🗳️ " + strings.Title(VOTETYPE) + " " + targetname).
				SetDescription(VOTEDESCRIPT + " The vote ends if the count is over (" +
					strconv.Itoa(REQVOTES) + ") or 24 hours from now.\n\nYES: 0\nNAH: 0").
				SetColor(votecolor).MessageEmbed

			VOTEEMBED = emb

			yesbutt := discordgo.Button{
				Label:    "YES",
				Style:    discordgo.SuccessButton,
				Disabled: false,
				Emoji: discordgo.ButtonEmoji{
					Name:     "",
					ID:       "738842746525057065",
					Animated: false,
				},
				URL:      "",
				CustomID: "yes",
			}
			nobutt := discordgo.Button{
				Label:    "NAH",
				Style:    discordgo.DangerButton,
				Disabled: false,
				Emoji: discordgo.ButtonEmoji{
					Name:     "❌",
					ID:       "",
					Animated: false,
				},
				URL:      "",
				CustomID: "no",
			}

			msg := discordgo.MessageSend{
				Content:    "@here",
				TTS:        false,
				Components: []discordgo.MessageComponent{discordgo.ActionsRow{Components: []discordgo.MessageComponent{yesbutt, nobutt}}},
				Files:      nil,
				Embed:      emb,
			}

			VOTEMSG, err = s.ChannelMessageSendComplex(j.ChannelID, &msg)
			fmt.Println(err)
			if err != nil {
				s.ChannelMessageSend(j.ChannelID, "Error: Discord refused the request, please try again.")
				return
			}

			deactivateVoting()
		}
	}
}

func evtCastVote(s *discordgo.Session, in *discordgo.InteractionCreate) {
	/*fmt.Println(j.MessageReaction.ChannelID)
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
	}*/

	// convert vote into int
	vote := int8(-1)
	rxvote := in.MessageComponentData().CustomID
	if rxvote == "yes" {
		vote = 1
	} else {
		vote = 0
	}

	// store the vote
	BALLOTMUTEX.Lock()
	BALLOTCOUNT[in.Member.User.ID] = vote
	BALLOTMUTEX.Unlock()

	// recount
	yes := 0
	no := 0
	BALLOTMUTEX.RLock()
	for _, v := range BALLOTCOUNT {
		if v == 1 {
			yes += 1
		} else {
			no += 1
		}
	}
	BALLOTMUTEX.RUnlock()

	// if we are over end voting
	if yes >= REQVOTES {
		go performVoteTypeAction(s)
		reenableVoting()
		VOTEACTIVE = false

		BALLOTMUTEX.Lock()
		defer BALLOTMUTEX.Unlock()

		VOTEEMBED.Description = VOTEDESCRIPT + " The vote ends if the count is over (" +
			strconv.Itoa(REQVOTES) + ") or 24 hours from now.\n\nYES: " +
			strconv.Itoa(yes) + "\nNAH: " + strconv.Itoa(no) +
			"\n\n_Voting has concluded._"
		resp := discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseUpdateMessage,
			Data: &discordgo.InteractionResponseData{
				TTS:             false,
				Content:         "@here",
				Components:      []discordgo.MessageComponent{},
				Embeds:          []*discordgo.MessageEmbed{VOTEEMBED},
				AllowedMentions: nil,
			},
		}

		err := s.InteractionRespond(in.Interaction, &resp)
		fmt.Println(err)
	} else if no >= REQVOTES {
		BALLOTMUTEX.Lock()
		defer BALLOTMUTEX.Unlock()

		VOTEEMBED.Description = VOTEDESCRIPT + " The vote ends if the count is over (" +
			strconv.Itoa(REQVOTES) + ") or 24 hours from now.\n\nYES: " +
			strconv.Itoa(yes) + "\nNAH: " + strconv.Itoa(no) +
			"\n\n_Voting has concluded._"
		resp := discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseUpdateMessage,
			Data: &discordgo.InteractionResponseData{
				TTS:             false,
				Content:         "@here",
				Components:      []discordgo.MessageComponent{},
				Embeds:          []*discordgo.MessageEmbed{VOTEEMBED},
				AllowedMentions: nil,
			},
		}

		err := s.InteractionRespond(in.Interaction, &resp)
		fmt.Println(err)
	}

	BALLOTMUTEX.Lock()
	defer BALLOTMUTEX.Unlock()

	VOTEEMBED.Description = VOTEDESCRIPT + " The vote ends if the count is over (" +
		strconv.Itoa(REQVOTES) + ") or 24 hours from now.\n\nYES: " +
		strconv.Itoa(yes) + "\nNAH: " + strconv.Itoa(no)
	resp := discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			TTS:             false,
			Content:         "@here",
			Components:      nil,
			Embeds:          []*discordgo.MessageEmbed{VOTEEMBED},
			AllowedMentions: nil,
		},
	}

	err := s.InteractionRespond(in.Interaction, &resp)
	fmt.Println(err)
}

func deactivateVoting() {
	VOTEACTIVE = false
	VOTEENABLED = false
	go Dispatch(time.Hour*24*7, time.Minute, reenableVoting, "reactivateVoting")
}

func reenableVoting() {
	VOTEENABLED = true
}

func testButtons(s *discordgo.Session, j *discordgo.MessageCreate) {
	/*if command(j, "test", ROLEELDER, false) {
	            butt := discordgo.Button{
	                    Label:    "YES",
	                    Style:    discordgo.SuccessButton,
	                    Disabled: false,
	                    Emoji: discordgo.ButtonEmoji{
	                            Name:     "",
	                            ID:       "738842746525057065",
	                            Animated: false,
	                    },
	                    URL:      "",
	                    CustomID: "test",
	            }
	            msg := discordgo.MessageSend{
	                    Content:    "test",
	                    TTS:        false,
	                    Components: []discordgo.MessageComponent{discordgo.ActionsRow{Components: []discordgo.MessageComponent{butt}}},
	                    Files:      nil,
	            }

	            _, err := s.ChannelMessageSendComplex(j.ChannelID, &msg)
	            fmt.Println(err)
	}*/
}

func getInteract(s *discordgo.Session, n *discordgo.InteractionCreate) {
	/*fmt.Println(n.Interaction.MessageComponentData().CustomID)
	          resp := discordgo.InteractionResponse{
	                  Type: discordgo.InteractionResponseUpdateMessage,
	                  Data: &discordgo.InteractionResponseData{
	                          TTS:             false,
	                          Content:         "voting has stopped.",
	                          Components:      []discordgo.MessageComponent{},
	                          Embeds:          nil,
	                          AllowedMentions: nil,
	                  },
	          }

	          err := s.InteractionRespond(n.Interaction, &resp)

	fmt.Println(err)*/
}

func performVoteTypeAction(s *discordgo.Session) {
	switch VOTETYPE {
	case "ban":
		err := s.GuildBanCreate(GUILDID, VOTEUSER.User.ID, 0)
		if err != nil {
			fmt.Println(err)
		}

	case "kick":
		err := s.GuildMemberDelete(GUILDID, VOTEUSER.User.ID)
		if err != nil {
			fmt.Println(err)
		}

	case "promote":
		err := s.GuildMemberRoleAdd(GUILDID, VOTEUSER.User.ID, ROLEHONORARY)
		if err != nil {
			fmt.Println(err)
		}
	}

	reenableVoting()
	VOTEACTIVE = false
}
