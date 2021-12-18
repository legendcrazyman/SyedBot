package utilities

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
)

// todo: add something for url parsing here
func CountVotes(s *discordgo.Session, m *discordgo.MessageCreate, amount int) bool {
	s.MessageReactionAdd(m.ChannelID, m.ID, "✅")
	s.MessageReactionAdd(m.ChannelID, m.ID, "🖕")
	time.Sleep(10 * time.Second)
	reactionMessage, _ := s.ChannelMessage(m.ChannelID, m.ID)

	upvote := 0
	downvote := 0
	for _, x := range reactionMessage.Reactions {
		if x.Emoji.Name == "✅" {
			upvote = x.Count
		} else if x.Emoji.Name == "🖕" {
			downvote = x.Count
		}
	} 
	if upvote > amount && upvote - downvote > 1 {
		return true
	} else {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Not enough upvotes! (need at least %d)", amount))
		return false
	}
}