package main

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type Bank struct {
	store *redis.Client
}

func NewBank(client *redis.Client) *Bank {
	client.HSet(context.Background(), "bank01", "balance", 0).Result()
	return &Bank{
		store: client,
	}
}

func (b *Bank) Deposit(ctx context.Context, boss *Boss, amount int) {
	lock := NewRedisLock(b.store, "bank01", boss.Name, 5*time.Second)
	done := false
	for i := range maxRetries {
		if done {
			break
		}

		locked, err := lock.Aquire(ctx)
		defer func() {
			err := lock.Release(ctx)
			if err != nil {
				fmt.Println("cannot release lock")
				return
			}
			fmt.Printf("%v release the lock! \n", boss.Name)
		}()

		if err != nil {
			return
		}
		if !locked {
			fmt.Printf("%v cannot aquire lock when deposit, retry %d times... \n", boss.Name, i)
			time.Sleep(retryDelay)
			continue
		}

		fmt.Printf("%v get the lock! \n", boss.Name)
		balance := b.GetBalance(ctx)
		balance += amount
		b.store.HSet(ctx, "bank01", "balance", balance).Result()
		boss.Balance += amount
		done = true
	}
}

func (b *Bank) Withdraw(ctx context.Context, e *Employee, amount int) error {
	lock := NewRedisLock(b.store, "bank01", e.Name, 5*time.Second)
	done := false

	for i := range maxRetries {
		if done {
			break
		}

		locked, err := lock.Aquire(ctx)

		defer func() {
			err := lock.Release(ctx)
			if err != nil {
				fmt.Println("cannot release lock")
				return
			}
			fmt.Printf("%v release the lock! \n", e.Name)
		}()

		if err != nil {
			return fmt.Errorf("error acquiring lock: %v", err)
		}
		if b.GetBalance(ctx) < amount {
			return fmt.Errorf("insufficient balance")
		}

		if !locked {
			fmt.Printf("%v could not acquire lock, retry %d times...\n", e.Name, i+1)
			time.Sleep(retryDelay)
			continue
		}

		balance := b.GetBalance(ctx)
		balance -= amount
		b.store.HSet(ctx, "bank01", "balance", balance).Result()
		b.store.HSet(ctx, e.Name, "balance", e.GetBalance(ctx)+amount).Result()
		done = true
	}

	return nil
}

func (b *Bank) GetBalance(ctx context.Context) int {
	res, _ := strconv.Atoi(b.store.HGet(ctx, "bank01", "balance").Val())
	return res
}
