package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type Boss struct {
	Name    string
	bank    *Bank
	Balance int
}

func NewBoss(bank *Bank) *Boss {
	return &Boss{
		bank: bank,
	}
}

func (b *Boss) Pay(wg *sync.WaitGroup) {
	defer wg.Done()
	for range 3 {
		amount := randRange(1, 10)
		fmt.Printf("Boss pay: %d \n", amount)
		b.bank.Deposit(context.Background(), b, amount)
		fmt.Printf("Current balance after deposit: %d \n", b.bank.GetBalance(context.Background()))
		time.Sleep(2 * time.Second)
	}
}
