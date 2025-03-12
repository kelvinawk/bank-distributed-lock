package main

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

type Employee struct {
	Name    string `redis:"name"`
	Balance int    `redis:"balance"`
	bank    *Bank
	store   *redis.Client
}

func NewEmployee(name string, bank *Bank, client *redis.Client) *Employee {
	emp := &Employee{
		Name:    name,
		Balance: 0,
		bank:    bank,
		store:   client,
	}
	client.HSet(context.Background(), name, emp).Result()

	return emp
}

func (e *Employee) GetBalance(ctx context.Context) int {
	emp, _ := e.store.HGet(ctx, e.Name, "balance").Result()
	balance, _ := strconv.Atoi(emp)
	return balance
}

func (e *Employee) GetMoney(wg *sync.WaitGroup) {
	defer wg.Done()
	for range 5 {
		time.Sleep(500 * time.Millisecond)
		amount := randRange(1, 10)
		fmt.Printf("%s withdraw: %d \n", e.Name, amount)
		err := e.bank.Withdraw(context.Background(), e, amount)
		if err != nil {
			fmt.Printf("Error withdrawing: %v \n", err)
			continue
		}
		fmt.Printf("%v withdraw success", e.Name)
		fmt.Printf("Current balance after withdraw by %s: %d \n", e.Name, e.bank.GetBalance(context.Background()))
	}
}
