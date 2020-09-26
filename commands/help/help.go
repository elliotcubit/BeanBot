
	// Help will not be its own module
  /*
	if strings.HasPrefix(m.Content, "!help") {
		embed := &discordgo.MessageEmbed{Color: 0x3498DB}
		embed.Title = "Commands"

		helpMessage := ""
		for _, handler := range activeModules {
			helpMessage += strings.Join(handler.Prefixes(), ", ")
			helpMessage += ": " + handler.Help() + "\n\n"
		}

		embed.Description = helpMessage

		s.ChannelMessageSendEmbed(m.ChannelID, embed)
	}
  */
