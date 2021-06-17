package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

var (
	ErrUserPresent    = errors.New("user already present")
	ErrUserNotPresent = errors.New("user not yet present")
)

type User struct {
	DiscordID      string `json:"id"`
	FirstName      string `json:"first-name"`
	LastName       string `json:"last-name"`
	GraduatingYear int    `json:"class-of"`
	LastFmAccount  string `json:"last-fm"`
}

type UserDB interface {
	Get(key string) (*User, bool)
	Add(key string, user *User) error
	Put(key string, user *User) error
	Save() error
}

type UserDBImpl struct {
	path  string
	users map[string]User // mapped by id
}

func NewUserDB(path string) (UserDB, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	db := UserDBImpl{
		path:  path,
		users: make(map[string]User),
	}
	err = json.Unmarshal(data, &db.users)
	if err != nil {
		return nil, err
	}
	return &db, nil
}

func (u *UserDBImpl) Get(key string) (*User, bool) {
	if user, present := u.users[key]; !present {
		return nil, false
	} else {
		return &user, true
	}
}

func (u *UserDBImpl) Add(key string, user *User) error {
	if _, present := u.users[key]; present {
		return ErrUserPresent
	}
	return u.Put(key, user)
}

func (u *UserDBImpl) Put(key string, user *User) error {
	u.users[key] = *user
	if err := u.Save(); err != nil {
		delete(u.users, key)
		return err
	}
	return nil
}

func (u *UserDBImpl) Save() error {
	marshalled, err := json.Marshal(u.users)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(u.path, marshalled, 0644)
	if err != nil {
		return err
	}
	return nil
}

type UserInt interface {
	ById(id string) (*User, error)
	Register(id, firstName, lastName string, gradYear int) (*User, error)
	Update(id, firstName, lastName string, gradYear int) (*User, error)
}

type UserIntImpl struct {
	db UserDB
}

func NewUserInt(path string) (UserInt, error) {
	db, err := NewUserDB(path)
	if err != nil {
		return nil, err
	}
	return &UserIntImpl{
		db: db,
	}, nil
}

func (u *UserIntImpl) ById(id string) (*User, error) {
	if user, present := u.db.Get(id); !present {
		return nil, ErrUserNotPresent
	} else {
		return user, nil
	}
}

func (u *UserIntImpl) Register(id, firstName, lastName string, gradYear int) (*User, error) {
	if _, present := u.db.Get(id); present {
		return nil, ErrUserPresent
	}

	if gradYear >= time.Now().Year()+7 {
		return nil, fmt.Errorf("Graduation year is too far in the future.")
	} else if gradYear <= 2019 { // sus lower bound
		return nil, fmt.Errorf("Graduation year is too far in the past.")
	}

	user := &User{
		DiscordID:      id,
		FirstName:      firstName,
		LastName:       lastName,
		GraduatingYear: gradYear,
	}

	if err := u.db.Add(id, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (u *UserIntImpl) Update(id, firstName, lastName string, gradYear int) (*User, error) {
	user, present := u.db.Get(id)
	if !present {
		return nil, ErrUserNotPresent
	}

	if firstName != "" {
		user.FirstName = firstName
	}
	if lastName != "" {
		user.LastName = lastName
	}
	if gradYear != 0 {
		user.GraduatingYear = gradYear
	}

	if err := u.db.Put(id, user); err != nil {
		return nil, err
	}
	return user, nil
}
