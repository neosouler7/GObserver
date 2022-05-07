package tg

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/neosouler7/GObserver/collector"
	"github.com/neosouler7/GObserver/db"
	"github.com/neosouler7/GObserver/utils"
)

const (
	bucketClearMsg = "Bucket successfully cleared."
)

// Returns personal informations.
func whoami(update tgbotapi.Update) string {
	messageFromID := update.Message.From.ID
	firstName := update.Message.From.FirstName
	lastName := update.Message.From.LastName
	return fmt.Sprintf("Hi %s %s !\nYou're ID is %d", firstName, lastName, messageFromID)
}

// Returns chat informations.
func whereami(update tgbotapi.Update) string {
	messageFromID := update.Message.From.ID
	chatId := update.Message.Chat.ID
	return fmt.Sprintf("Hi %d !\nYou're now in chat %d", messageFromID, chatId)
}

// Clears database manually.
func clear() string {
	err := db.InitBucket()
	if err != nil {
		return errors.New("something wrong").Error()
	}
	return bucketClearMsg
}

func orderbook(args string) string {
	fmt.Printf(args)
	s := strings.Split(args, " ")
	price, ok := collector.ObMap.Load(fmt.Sprintf("%s:%s:%s", s[0], s[1], s[2]))
	if !ok {
		HandleErr(errors.New("no such key on ObMap."))
	}
	return price.(string)
}

// Check current collected data.
func status() string {
	cp := fmt.Sprintf("%s", db.GetCheckPoint(db.LastUpdatedAt))
	cpAsInt, err := strconv.Atoi(cp)
	HandleErr(err)

	tm := time.Unix(int64(cpAsInt), 0)
	cpAsTime := fmt.Sprintf("%s", tm.Format(TimeFormat))
	statusMsg := fmt.Sprintf("LastUpdatedAt: %s", cpAsTime)
	return statusMsg
}

// to be deprecated. sample code.
func save(args string) string {
	s := "test"
	m := db.MoldStruct{
		Payload: []byte(s),
	}
	db.UpdateMold(utils.ToBytes(m))
	HandleErr(db.SaveCheckPoint(args))
	// db.SaveCheckPoint(db.LastUpdatedAt)
	return "successfully saved!"
}
