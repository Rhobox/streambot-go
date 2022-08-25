package commands

import (
	"errors"
	"github.com/bwmarrin/discordgo"
	"strings"

	"github.com/aricodes-oss/std"
	"streambot/config"
	"streambot/twitch"
)

var Config = config.Config

var helix = twitch.Helix
var logger = std.Logger

var ErrNotACommand = errors.New("message does not contain a command")
var ErrNotAdmin = errors.New("user is not an administrator of this channel")

type Command struct {
	Name      string
	Arguments []string

	Raw          string
	RawArguments string

	Member *discordgo.Member

	Session *discordgo.Session
	Event   *discordgo.MessageCreate

	Reply func(string)
}

func messageIsCommand(message string) bool {
	return len(message) >= 2 && string(message[0]) == "!"
}

func userIsAdmin(s *discordgo.Session, m *discordgo.MessageCreate) bool {
	perms, err := s.UserChannelPermissions(m.Author.ID, m.ChannelID)
	if err != nil {
		return false
	}

	return perms&discordgo.PermissionAdministrator != 0
}

func Parse(s *discordgo.Session, m *discordgo.MessageCreate) (result *Command, err error) {
	result = &Command{
		Member:  m.Member,
		Session: s,
		Event:   m,
		Reply:   func(msg string) { s.ChannelMessageSend(m.ChannelID, msg) },
		Raw:     m.Content,
	}
	if !messageIsCommand(m.Content) {
		return nil, ErrNotACommand
	}

	if !userIsAdmin(s, m) {
		return nil, ErrNotAdmin
	}

	segments := strings.Split(m.Content[1:], " ")

	result.Name = segments[0]
	result.Arguments = segments[1:]
	result.RawArguments = strings.Join(segments[1:], " ")

	return
}
