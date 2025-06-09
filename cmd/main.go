package main

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"

	"workerPool/handlers"
)

func main() {
	log := slog.New(slog.NewTextHandler(os.Stdout, nil))

	pool := handlers.NewPool(log)

	pool.AddWorker()
	pool.AddWorker()

	scanner := bufio.NewScanner(os.Stdin)
	handlers.PrintHelp()

	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}

		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			continue
		}

		parts := strings.SplitN(input, " ", 2)
		command := parts[0]

		switch command {
		case "add":
			pool.AddWorker()

		case "addid":
			if len(parts) < 2 {
				fmt.Println("Укажите ID воркера (например: addid 5)")
				continue
			}
			id, err := strconv.Atoi(parts[1])
			if err != nil || id < 1 {
				fmt.Println("Некорректный ID воркера. ID должен быть числом и больше нуля")
				continue
			}
			if err := pool.AddWorkerByID(id); err != nil {
				fmt.Println("Ошибка:", err)
			}

		case "remove":
			pool.RemoveLastWorker()

		case "removeid":
			if len(parts) < 2 {
				fmt.Println("Укажите ID воркера (например: removeid 1)")
				continue
			}
			id, err := strconv.Atoi(parts[1])
			if err != nil {
				fmt.Println("Некорректный ID воркера")
				continue
			}
			pool.RemoveWorkerByID(id)

		case "process":
			if len(parts) < 2 {
				fmt.Println("Укажите задачу (например: process task1)")
				continue
			}
			pool.JobQueue <- parts[1]

		case "ls":
			pool.PrintWorkers()

		case "help":
			handlers.PrintHelp()
		case "exit":
			pool.Stop()
			fmt.Println("Завершение программы")
			return

		default:
			fmt.Println("Неизвестная команда")
		}
	}
}
