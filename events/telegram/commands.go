package telegram

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"GoodDeedDAO/lib/e"
	"GoodDeedDAO/storage"
)

const (
	UserInfoCmd = "/userinfo"
	HelpCmd     = "/help"
	StartCmd    = "/start"
	AddKarmaCmd = "/addkarma"
)

func (p *Processor) doCmd(text string, chatID int, username string) error {
	text = strings.TrimSpace(text)

	log.Printf("got new command '%s' from '%s'", text, username)

	var comm, data string
	comm, data, _ = strings.Cut(text, " ")

	//if isAddCmd(text) {
	//	return p.savePage(chatID, text, username)
	//}

	switch comm {
	case UserInfoCmd:
		return p.sendUserInfo(chatID, username)
	case HelpCmd:
		return p.sendHelp(chatID)
	case StartCmd:
		return p.sendHello(chatID, username)
	case AddKarmaCmd:
		return p.AddKarma(chatID, username, data)
	default:
		return p.tg.SendMessage(chatID, msgUnknownCommand)
	}
}

// TODO remove it
func (p *Processor) savePage(chatID int, pageURL string, username string) (err error) {
	defer func() { err = e.WrapIfErr("can't do command: save page", err) }()

	//page := &storage.User{
	//	URL:      pageURL,
	//	UserName: username,
	//}

	//isExists, err := p.storage.IsExists(context.Background(), page)
	//if err != nil {
	//	return err
	//}
	//if isExists {
	//	return p.tg.SendMessage(chatID, msgAlreadyExists)
	//}
	//
	//if err := p.storage.Save(context.Background(), page); err != nil {
	//	return err
	//}
	//
	//if err := p.tg.SendMessage(chatID, msgSaved); err != nil {
	//	return err
	//}

	return nil
}

func (p *Processor) sendUserInfo(chatID int, username string) (err error) {
	defer func() { err = e.WrapIfErr("can't do command: sendUserInfo", err) }()

	user, err := p.storage.GetUserInfo(context.Background(), username)
	if err != nil && !errors.Is(err, storage.ErrNoSavedPages) {
		return err
	}
	if errors.Is(err, storage.ErrNoSavedPages) {
		return p.tg.SendMessage(chatID, msgUserNotFound)
	}

	text := "User: " + user.UserName +
		"\nKarma: " + strconv.Itoa(user.Karma) +
		"\ndeeds: " + strconv.Itoa(user.Deeds) +
		"\nvalidations: " + strconv.Itoa(user.Validations)

	if err := p.tg.SendMessage(chatID, text); err != nil {
		return err
	}

	return nil
}

func (p *Processor) sendHelp(chatID int) error {
	return p.tg.SendMessage(chatID, msgHelp)
}

func (p *Processor) sendHello(chatID int, username string) error {
	err := p.storage.AddUser(context.Background(), chatID, username)
	if err != nil {
		fmt.Errorf("can't add user: %s", username)
	}
	// fmt.Printf("user %s added", username)
	return p.tg.SendMessage(chatID, msgHello)
}

func (p *Processor) AddKarma(chatID int, username string, data string) error {
	karma, err := strconv.Atoi(data)
	fmt.Printf("Tryin' to add %d karma\n", karma)
	if err == nil {
		return p.storage.AddKarma(context.Background(), username, karma)
	}

	return p.tg.SendMessage(chatID, "Add "+data+" karma")
}

//func isAddCmd(text string) bool {
//	return isURL(text)
//}
//
//func isURL(text string) bool {
//	u, err := url.Parse(text)
//
//	return err == nil && u.Host != ""
//}
