module github.com/Starz0r/Decision

go 1.13

replace github.com/bwmarrin/discordgo => ./lib/discordgo

require (
	github.com/bwmarrin/discordgo v0.22.0
	github.com/natefinch/lumberjack v2.0.0+incompatible // indirect
	github.com/rs/zerolog v1.19.0 // indirect
	github.com/spidernest-go/logger v0.0.0-20191128190838-520d89ea00af
)
