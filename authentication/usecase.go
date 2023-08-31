package authentication

import (
	"errors"
	"fmt"
	"sync"
)

var (
	database map[int]User
	counter  int
)

var verificationQueue chan User

type UseCase interface {
	CreateUser(*User) error
	GetAllUsers() ([]User, error)
	VerifyUser(int) error
	QueueForVerification(User)
	GetUser(int) (*User, error)
	UpdateUser(User) error
}

type useCaseImpl struct {
	*sync.RWMutex
}

func NewUseCase() UseCase {
	counter = 0
	database = map[int]User{}
	verificationQueue = make(chan User, 30000000)
	return &useCaseImpl{}
}

func (u *useCaseImpl) CreateUser(user *User) error {
	counter += 1
	user.UserID = counter

	if _, ok := database[counter]; ok {
		return errors.New("duplicated user id")
	}

	database[counter] = *user

	putToQueue(*user)
	return nil
}

func putToQueue(user User) {
	verificationQueue <- user
}

func (u *useCaseImpl) VerifyUser(workerNumber int) error {
	select {
	case currentUser := <-verificationQueue:
		if currentUser.Verified {
			fmt.Println("user already verified...")
			return nil
		}
		currentUser.Verified = true

		database[currentUser.UserID] = currentUser

		return nil
	default:
		fmt.Printf("no user in queue, exiting worker %d\n", workerNumber)
		return nil
	}
}

func (u *useCaseImpl) GetAllUsers() (users []User, err error) {
	for _, user := range database {
		users = append(users, user)
	}

	return
}

func (u *useCaseImpl) GetUser(userId int) (*User, error) {
	if user, ok := database[userId]; !ok {
		return nil, errors.New("user doesn't exist")
	} else {
		return &user, nil
	}
}

func (u *useCaseImpl) UpdateUser(user User) error {
	if _, ok := database[user.UserID]; !ok {
		return errors.New("user doesn't exist")
	} else {
		database[user.UserID] = user
	}

	return nil
}

func (u *useCaseImpl) QueueForVerification(user User) {
	verificationQueue <- user
}
