package middleware

import (
	"errors"
	tele "gopkg.in/telebot.v3"
	"hamsterbot/internal/app/constants"
	"hamsterbot/internal/app/models"
	"hamsterbot/pkg/metrics"
)

type User interface {
	GetById(id int64) (models.User, error)
	Add(user models.User) error
	Update(user models.User) error
	Delete(id int64) error
}

type Endpoint struct {
	Bot  *tele.Bot
	User User
}

func (e *Endpoint) IsUser(next tele.HandlerFunc) tele.HandlerFunc {
	return func(c tele.Context) error {
		metrics.IncomingMessages.Inc()

		data, err := e.User.GetById(c.Sender().ID)
		if err != nil {
			if errors.Is(err, constants.ErrUserNotFound) {
				user := models.User{
					ID:        c.Sender().ID,
					Username:  c.Sender().Username,
					Firstname: c.Sender().FirstName,
					Lastname:  c.Sender().LastName,
					IsPremium: c.Sender().IsPremium,
				}

				err = e.User.Add(user)
				if err != nil {
					return err
				}
			}

			return err
		}

		if data.Username != c.Sender().Username || data.Firstname != c.Sender().FirstName || data.Lastname != c.Sender().LastName || data.IsPremium != c.Sender().IsPremium {
			err = e.User.Update(data)
			if err != nil {
				return err
			}
		}

		return next(c)
	}
}
