package users

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"hamsterbot/internal/app/constants"
	"hamsterbot/internal/app/models"
	"hamsterbot/pkg/cache"
	"hamsterbot/pkg/db"
)

type Service struct{}

func New() *Service {
	return &Service{}
}

func (s Service) GetUserById(id int64) (models.User, error) {
	var data models.User

	cacheKey := fmt.Sprintf("user:%d", id)
	cacheValue, err := cache.Rdb.Get(cache.Ctx, cacheKey).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return models.User{}, err
	} else if cacheValue != "" {
		err = json.Unmarshal([]byte(cacheValue), &data)
		if err != nil {
			return models.User{}, err
		}
	}

	rows, err := db.Conn.Query(`SELECT * FROM users WHERE id = $1`)
	if err != nil {
		return models.User{}, err
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&data)
		if err != nil {
			return models.User{}, err
		}
	}
	if data.ID == 0 {
		return models.User{}, constants.ErrUserNotFound
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return models.User{}, err
	}
	err = cache.Rdb.Set(cache.Ctx, cacheKey, jsonData, 0).Err()
	if err != nil {
		return models.User{}, err
	}

	return data, nil
}

func (s Service) DeleteUser(id int64) error {
	rows, err := db.Conn.Queryx(`DELETE FROM users WHERE id = $1`, id)
	if err != nil {
		return err
	}
	defer rows.Close()

	return nil
}

func (s Service) UpdateUser(user models.User) error {
	userOld, err := s.GetUserById(user.ID)
	if err != nil {
		return err
	}

	tx, err := db.Conn.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if user.Username != userOld.Username {
		_, err = tx.Exec(`UPDATE users SET username = $1 WHERE id = $2`, user.Username, user.ID)
		if err != nil {
			return err
		}
	}

	if user.Firstname != userOld.Firstname {
		_, err = tx.Exec(`UPDATE users SET firstname = $1 WHERE id = $2`, user.Firstname, user.ID)
		if err != nil {
			return err
		}
	}

	if user.Lastname != userOld.Lastname {
		_, err = tx.Exec(`UPDATE users SET lastname = $1 WHERE id = $2`, user.Lastname, user.ID)
		if err != nil {
			return err
		}
	}

	if user.IsPremium != userOld.IsPremium {
		_, err = tx.Exec(`UPDATE users SET ispremium = $1 WHERE id = $2`, user.IsPremium, user.ID)
		if err != nil {
			return err
		}
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (s Service) AddUser(user models.User) error {
	rows, err := db.Conn.Queryx(`INSERT INTO users (id, username, firstname, lastname, ispremium) VALUES ($1, $2, $3, $4, $5)`, user.ID, user.Username, user.Firstname, user.Lastname, user.IsPremium)
	if err != nil {
		return err
	}
	defer rows.Close()

	return nil
}
