package discord

import (
	"fmt"
	"github.com/Karitham/WaifuBot/database"
	"log"

	"github.com/Karitham/WaifuBot/query"
	"github.com/andersfylling/disgord"
)

func verify(data *disgord.MessageCreate, args CmdArguments) {
	// Verify if user possesses the Waifu
	desc, isValid := verifyWaifuValid(data, args)

	// Get char
	resp, err := query.CharSearch(args.ParseArgToSearch())
	if err != nil {
		log.Println(err)
	}

	if isValid {
		// Send confirmation Message
		_, err = client.CreateMessage(
			ctx,
			data.Message.ChannelID,
			&disgord.CreateMessageParams{
				Embed: &disgord.Embed{
					Title:       "Waifu Verification",
					Thumbnail:   &disgord.EmbedThumbnail{URL: resp.Character.Image.Large},
					Description: desc,
					Timestamp:   data.Message.Timestamp,
					Color:       0x43e99a,
				},
			},
		)
		if err != nil {
			log.Println("There was an error when sending the verification message: ", err)
		}
	} else {
		// Send message
		_, err := client.CreateMessage(
			ctx,
			data.Message.ChannelID,
			&disgord.CreateMessageParams{
				Embed: &disgord.Embed{
					Title:       "Waifu Verification",
					Thumbnail:   &disgord.EmbedThumbnail{URL: resp.Character.Image.Large},
					Description: desc,
					Timestamp:   data.Message.Timestamp,
					Color:       0xcc0000,
				},
			},
		)
		if err != nil {
			log.Println("Create message returned error :", err)
		}
	}
}

// Verify if user possesses the Waifu
func verifyWaifuValid(data *disgord.MessageCreate, arg CmdArguments) (desc string, isValid bool) {
	if len(arg) > 0 {
		resp, _ := query.CharSearch(arg.ParseArgToSearch())
		switch {
		case resp.Character.ID == 0:
			return fmt.Sprintf("Error, %d is not a valid WaifuID,\nRefer to %shelp to see this command's syntax", resp.Character.ID, conf.Prefix), false
		case data.Message.Mentions == nil:
			return fmt.Sprintf("Error, please tag a discord user,\nRefer to %shelp to see this command's syntax", conf.Prefix), false
		case !database.VerifyWaifuStruct{UserID: data.Message.Mentions[0].ID, CharID: resp.Character.ID}.VerifyWaifu():
			return fmt.Sprintf("%s does not possess %s.", data.Message.Mentions[0].Username, resp.Character.Name.Full), false
		default:
			return fmt.Sprintf("%s possesses %s.", data.Message.Mentions[0].Username, resp.Character.Name.Full), true
		}
	}
	return "Please enter arguments,\nRefer to help to see how to use this command", false
}

// Help function for Verify
func verifyHelp(data *disgord.MessageCreate) {
	_, err := client.CreateMessage(
		ctx,
		data.Message.ChannelID,
		&disgord.CreateMessageParams{
			Embed: &disgord.Embed{
				Title: "Verify Help || alias v",
				Description: fmt.Sprintf(
					"This is the help for the verify functionality\n\n"+
						"This permits you to check if one of your friends has got the waifu you want.\n"+
						"You can give a waifu to another user using the following syntax :\n"+
						"`%sverify ID @User`",
					conf.Prefix,
				),
				Footer: &disgord.EmbedFooter{
					Text: fmt.Sprintf("Help requested by %s", data.Message.Author.Username),
				},
				Timestamp: data.Message.Timestamp,
				Color:     0xeec400,
			},
		},
	)
	if err != nil {
		log.Println("There was an error sending help for verifying waifu: ", err)
	}
}
