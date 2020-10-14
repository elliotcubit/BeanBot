package help

import (
	"beanbot/handlers"

	"github.com/bwmarrin/discordgo"
)

func init() {
	handlers.RegisterActiveModule(
		Help{},
	)
}

type Help struct{}

func (h Help) Do(s *discordgo.Session, m *discordgo.MessageCreate) {
	embed := &discordgo.MessageEmbed{Color: 0x3498DB}
	embed.Title = "BeanBot Help"

	helpMessage := ""
	helpMessage += "Learn to count, and get beans for it.\n"
	helpMessage += "Usage: !bean [command] [options]\n\n"
	helpMessage += "Commands:\n"
	helpMessage += "bet: Gamble your beans on a coin flip. !bean bet [amount] @Someone\n"
	helpMessage += "configure: Set the channel I'll listen to for counting. !bean configure\n"
	helpMessage += "give: Give away your beans. !bean give [amount] @Someone\n"
	helpMessage += "help: Show this message. !bean help\n"
	helpMessage += "leaderboard: Who has the most/least beans? !bean leaderboard [top/bottom] [limit]\n"
	helpMessage += "max: What is the highest number this server can count to? !bean max\n"
	helpMessage += "query: How many beans does someone have? !bean query @Someone. [No @ means yourself]\n"
	helpMessage += "risk: How many beans will I lose if I make a mistake? !bean risk\n"

	embed.Description = helpMessage
	s.ChannelMessageSendEmbed(m.ChannelID, embed)
}

func (h Help) Prefixes() []string {
	return []string{"help"}
}
