package spy

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/sboon-gg/sboon-bot/pkg/discord"
)

const (
	rolesCmdName      = "roles"
	rolesListCmdName  = "list"
	rolesSetCmdName   = "set"
	rolesUnsetCmdName = "unset"
)

var mapRoleSubCommand = &discordgo.ApplicationCommandOption{
	Name:        rolesCmdName,
	Description: "Map discord roles to set when player is active on PRSPY",
	Type:        discordgo.ApplicationCommandOptionSubCommandGroup,
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        rolesListCmdName,
			Description: "List existing mappings",
		},
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        rolesSetCmdName,
			Description: "Set new mapping",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionRole,
					Name:        "role",
					Description: "Role user has that will map to active role (@everyone for no restriction)",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionRole,
					Name:        "active-role",
					Description: "Role user will receive when active in PRSPY",
					Required:    true,
				},
			},
		},
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        rolesUnsetCmdName,
			Description: "Unset existing mapping",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionRole,
					Name:        "role",
					Description: "Role user has",
					Required:    true,
				},
			},
		},
	},
}

func (b *Spy) rolesHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options[0].Options

	switch options[0].Name {
	case rolesListCmdName:
		b.rolesListHandler(s, i)
	case rolesSetCmdName:
		b.rolesSetHandler(s, i)
	case rolesUnsetCmdName:
		b.rolesUnsetHandler(s, i)
	}
}

func (b *Spy) rolesListHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	fields := []*discordgo.MessageEmbedField{}

	for _, role := range b.roleRepo.FindAll() {
		fields = append(fields, &discordgo.MessageEmbedField{
			Value: fmt.Sprintf("<@&%s> -> <@&%s>", role.DiscordID, role.ActiveRoleID),
		})
	}

	if len(fields) == 0 {
		fields = append(fields, &discordgo.MessageEmbedField{
			Value: "No mappings",
		})
	}

	resp := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "List of mapped roles",
			Embeds: []*discordgo.MessageEmbed{
				{
					Fields: fields,
				},
			},
			Flags: discordgo.MessageFlagsEphemeral,
		},
	}

	err := s.InteractionRespond(i.Interaction, resp)
	if err != nil {
		log.Print(err)
	}
}

func (b *Spy) rolesSetHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options[0].Options[0].Options

	role := options[0].Value.(string)
	activeRole := options[1].Value.(string)

	_, err := b.roleRepo.SetMapping(role, activeRole)
	if err != nil {
		discord.ErrorResponse(s, i.Interaction, err)
		return
	}

	b.rolesListHandler(s, i)
}

func (b *Spy) rolesUnsetHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options[0].Options[0].Options
	roleID := options[0].Value.(string)

	err := b.roleRepo.UnsetMapping(roleID)
	if err != nil {
		discord.ErrorResponse(s, i.Interaction, err)
		return
	}

	b.rolesListHandler(s, i)
}
