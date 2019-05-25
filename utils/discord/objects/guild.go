package objects

type Guild struct {
	Id string
	Name string
	Icon string
	Splash string
	Owner bool
	OwnerId string
	Permissions int
	Region string
	AfkChannelid string
	AfkTimeout int
	EmbedEnabled bool
	EmbedChannelId string
	VerificationLevel int
	DefaultMessageNotifications int
	ExplicitContentFilter int
	Roles []Role
	Emojis []Emoji
	Features []string
	MfaLevel int
	ApplicationId string
	WidgetEnabled bool
	WidgetChannelId string
	SystemChannelId string
	JoinedAt string
	Large bool
	Unavailable bool
	MemberCount int
	VoiceStates []VoiceState
	Members []Member
	Channels []Channel
	Presences []Presence
	MaxPresences int
	Maxmembers int
	VanityUrlCode string
	Description string
	Banner string
}

func (g *Guild) GetCategories() []Channel {
	var categories []Channel

	for _, channel := range g.Channels {
		if channel.Type == 4 {
			categories = append(categories, channel)
		}
	}

	return categories
}
