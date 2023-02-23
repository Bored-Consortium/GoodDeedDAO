package storage

import (
	"context"
	"crypto/sha1"
	"errors"
	"fmt"
	"io"

	"GoodDeedDAO/lib/e"
)

type Storage interface {
	AddUser(ctx context.Context, chatID int, username string) error
	AddKarma(ctx context.Context, username string) (*User, error)
	Remove(ctx context.Context, p *User) error
	IsUserInDb(ctx context.Context, username string) (bool, error)
	GetUserInfo(ctx context.Context, userName string) (*User, error)
}

var ErrNoSavedPages = errors.New("no saved pages")

type User struct {
	URL         string // TODO delete it
	Id          int
	UserName    string
	Karma       int
	Deeds       int
	Validations int
}

func (p User) Hash() (string, error) {
	h := sha1.New()

	if _, err := io.WriteString(h, p.URL); err != nil {
		return "", e.Wrap("can't calculate hash", err)
	}

	if _, err := io.WriteString(h, p.UserName); err != nil {
		return "", e.Wrap("can't calculate hash", err)
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
