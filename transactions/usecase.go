package transactions

import (
	"LemFi/authentication"
	"errors"
	"fmt"
)

var transactionQueue chan Transaction

type UseCase interface {
	Transact(Transaction) error
	ProcessTransaction(int) error
}

type useCaseImpl struct {
	authentication.UseCase
}

func NewUseCase(authUseCase authentication.UseCase) UseCase {
	transactionQueue = make(chan Transaction, 300000)
	return &useCaseImpl{
		UseCase: authUseCase,
	}
}

func (u *useCaseImpl) Transact(transaction Transaction) error {
	fmt.Println("here.....")
	sender, err := u.UseCase.GetUser(transaction.SenderId)
	if err != nil {
		return errors.New("sender doesn't exist")
	}

	fmt.Println("here......")
	_, err = u.UseCase.GetUser(transaction.ReceiverId)
	if err != nil {
		return errors.New("receiver doesn't exist")
	}

	fmt.Println("here.......")
	if sender.Balance < transaction.Amount {
		return errors.New("insufficient balance")
	}

	fmt.Println("here........")
	transactionQueue <- transaction
	return nil
}

func (u *useCaseImpl) ProcessTransaction(workerNumber int) error {
	select {
	case currentTransaction := <-transactionQueue:
		fmt.Println("processing transaction...")
		sender, _ := u.UseCase.GetUser(currentTransaction.SenderId)

		receiver, _ := u.UseCase.GetUser(currentTransaction.ReceiverId)

		if !sender.Verified {
			fmt.Println("sender not verified, queueing for verification...")
			u.UseCase.QueueForVerification(*sender)
			transactionQueue <- currentTransaction
			return nil
		}
		if !receiver.Verified {
			fmt.Println("receiver not verified, queueing for verification...")
			u.UseCase.QueueForVerification(*receiver)
			transactionQueue <- currentTransaction
			return nil
		}

		sender.Balance -= currentTransaction.Amount
		receiver.Balance += currentTransaction.Amount

		go func() {
			_ = u.UseCase.UpdateUser(*sender)
		}()
		go func() {
			_ = u.UseCase.UpdateUser(*receiver)
		}()
	default:
		fmt.Printf("no transaction in queue, exiting worker %d\n", workerNumber)
		return nil
	}
	return nil
}
