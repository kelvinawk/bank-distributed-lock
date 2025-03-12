package main

import (
	"context"
	"fmt"
	"math/rand/v2"
	"sync"
)

func main() {
	rds := NewRedisClient()
	bank := NewBank(rds)
	boss := NewBoss(bank)
	var wg sync.WaitGroup

	wg.Add(1)

	go boss.Pay(&wg)

	for i := 1; i < 5; i++ {
		emp := NewEmployee(fmt.Sprintf("Employee %d", i), bank, rds)
		wg.Add(1)
		go emp.GetMoney(&wg)
	}

	wg.Wait()
	fmt.Printf("Boss paid: %d \n", boss.Balance)
	fmt.Printf("Bank balance: %d \n", bank.GetBalance(context.Background()))
}

func randRange(min, max int) int {
	return rand.IntN(max-min) + min
}
