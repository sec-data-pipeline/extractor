package storage

import (
	"fmt"
	"log"
)

type Logger interface {
	Log(msg string)
}

type cloudWatch struct{}

func NewCloudWatch() *cloudWatch {
	return &cloudWatch{}
}

func (l *cloudWatch) Log(msg string) {
	fmt.Println(msg)
}

type console struct{}

func NewConsole() *console {
	return &console{}
}

func (l *console) Log(msg string) {
	log.Println(msg)
}
