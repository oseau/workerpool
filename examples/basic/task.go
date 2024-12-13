package main

import (
	"log"
	"time"
)

type TaskSendMail struct {
}

func (t *TaskSendMail) Execute() error {
	time.Sleep(time.Second * 1)
	log.Println("mail sent")
	return nil
}

type TaskSendSMS struct {
}

func (t *TaskSendSMS) Execute() error {
	time.Sleep(time.Second * 1)
	log.Println("sms sent")
	return nil
}
