package main

import (
	"github.com/EDDYCJY/go-gin-example/models"
	"github.com/robfig/cron"
	"log"
	"time"
)

func CronSetup() {

	go func() {
		log.Println("Starting...")

		c := cron.New()
		c.AddFunc("5 * * * * *", func() {
			log.Println("Run load book...")
			books, err := models.HotGtZeroBooks()
			if err != nil {
				log.Fatal("get hot gt 0 from db failed")
				log.Fatal(err)
			}
			for _, book := range books {
				print(book.BookName)
			}
		})

		c.Start()

		t1 := time.NewTimer(time.Second * 10)
		for {
			select {
			case <-t1.C:
				t1.Reset(time.Second * 10)
			}
		}
	}()

}
