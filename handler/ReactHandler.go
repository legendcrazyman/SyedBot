package messageHandler

import (
	"github.com/bwmarrin/discordgo"
)

func ReactHandler(s *discordgo.Session, r *discordgo.MessageReactionAdd) {


	if r.UserID == s.State.User.ID {
		return
	}
	if r.Emoji.Name == "📌" {
		reactionMessage, _ := s.ChannelMessage(r.ChannelID, r.MessageID)
		for _, x := range reactionMessage.Reactions {
			if x.Emoji.Name == "📌" && x.Count > 3 {
				s.ChannelMessagePin(r.ChannelID, r.MessageID)

			}
		}

	}
}