// Package tg notifies stored informations to subscribers via telegram.
package tg

import (
	"errors"
	"fmt"
	"log"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/neosouler7/GObserver/config"
	"github.com/neosouler7/GObserver/utils"
)

type asString struct {
	s string
}

// Returns error as a string.
func (e *asString) Error() string {
	return fmt.Sprintf("#ERROR\n\n%s", e.s)
}

type botConfig struct {
	token        string
	chatId       int64
	commanderIds []int64
}

var (
	bc  *botConfig
	bot *tgbotapi.BotAPI
)

const (
	TimeFormat            = "2006-01-02 15:04:05"
	wrongCommandArguments = "Check your command's arguments."
)

// Initialize and returns tg config & bot object.
func initBot() (*botConfig, *tgbotapi.BotAPI) {
	c := config.TgConfig()
	bc = &botConfig{
		token:        c.Token,
		chatId:       c.ChatId,
		commanderIds: c.CommanderIds,
	}

	botPointer, err := tgbotapi.NewBotAPI(bc.token)
	if err != nil {
		log.Panic(errors.New("Failed initializing bot"))
	}
	// botPointer.Debug = true

	bot = botPointer
	log.Printf("Authorized on account %s.", bot.Self.UserName)
	return bc, bot
}

// Returns message for non-commanders.
func getNonCommanderMsg(update tgbotapi.Update) string {
	messageFromID := update.Message.From.ID
	userName := update.Message.From.UserName
	firstName := update.Message.From.FirstName
	lastName := update.Message.From.LastName
	return fmt.Sprintf("Trying to approach!\n\nID: %d\nUser: %s\nFirst: %s\nLast: %s\n", messageFromID, userName, firstName, lastName)
}

// Get message updates and handles commands.
func listenMsg(bc *botConfig, bot *tgbotapi.BotAPI) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)
	for update := range updates {
		if update.Message == nil {
			continue
		}

		if !update.Message.IsCommand() { // server ignores non-command messages
			continue
		}

		args := update.Message.CommandArguments()
		commander := utils.Contains(bc.commanderIds, update.Message.From.ID)

		var msgText string
		switch update.Message.Command() {
		case "whoami":
			msgText = whoami(update)
		case "whereami":
			msgText = whereami(update)
		case "cl":
			if commander {
				msgText = clear()
			} else { // someone is trying to approach!
				msgText = getNonCommanderMsg(update)
			}
		case "ob":
			if commander {
				msgText = orderbook(args)
			} else { // someone is trying to approach!
				msgText = getNonCommanderMsg(update)
			}
		case "st":
			if commander {
				msgText = status()
			} else { // someone is trying to approach!
				msgText = getNonCommanderMsg(update)
			}
		case "save":
			if commander {
				if args == "lastUpdatedAt" || args == "createdAt" {
					msgText = save(args)
				} else {
					msgText = wrongCommandArguments
				}
			} else { // someone is trying to approach!
				msgText = getNonCommanderMsg(update)
			}
		default:
			msgText = "Wrong command :("
		}

		log.Printf("%s → %s", update.Message.Text, msgText)
		SendMsg(msgText)
	}
}

// Send tg message and handle errors.
func SendMsg(msgText string) {
	msg := tgbotapi.NewMessage(bc.chatId, fmt.Sprintf("%s\n\n%s", msgText, time.Now().Format(TimeFormat)))
	if _, err := bot.Send(msg); err != nil {
		log.Panic(err)
	}
}

// Sends error message that system admin should know & shutdown.
func HandleErr(err error) {
	if err != nil {
		SendMsg(err.Error())
	}
}

// Starts tg package.
func Start() {
	bc, bot := initBot()
	listenMsg(bc, bot)
}
