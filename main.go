package main

import (
	"errors"
	"fmt"
	"time"
)
//obj part
type Notification struct {
	ID    int
	Channel  string
	Message string
	TimeRef time.Time
}
type Notifier interface {
	Send(message string) error
	GetName() string
}

//verid caneaux de notifs:
type EmailNotifier struct {}
type PushNotifier struct {}
type SMSNotifer struct {
	PhoneNumber string
}

//matcher les méthodes avec l'interface pr email, push notif et SMS
func (e EmailNotifier) Send(message string) error {
	fmt.Printf("email envoie à partir de '%s' ", message)
	return nil
}
func (e EmailNotifier) GetName() string {return "Email"}

func (p PushNotifier) Send(message string) error {
	fmt.Printf("[PUSH] Envoi de '%s' réussi", message)
	return nil
}
func (p PushNotifier) GetName() string { return "Push" }

func (s SMSNotifer) Send(message string) error {
	if len(s.PhoneNumber) < 2 || s.PhoneNumber[:2] != "06" {
		return errors.New("Numéro invalide : doit commencer par '06'")
	}

	fmt.Printf("[SMS] Envoi de '%s' au %s réussi", message, s.PhoneNumber)
	return nil
}
func (s SMSNotifer) GetName() string {return "SMS"}

//storer la data
type Storer interface {
	Add(n *Notification) error
	GetAll() ([]*Notification, error)	
}
type MemoryStore struct {
	notifications []*Notification
	nextID int
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		notifications: []*Notification{},
		nextID: 1,
	}
}

func (ms *MemoryStore) Add(n *Notification) error {
	n.ID = ms.nextID
	ms.nextID++
	ms.notifications = append(ms.notifications, n)
	return nil
}

func (ms *MemoryStore) GetAll() ([]*Notification, error) {
	return ms.notifications, nil
}



func main() {
	store := NewMemoryStore()
	notifiers := []Notifier{
		EmailNotifier{},
		SMSNotifer{PhoneNumber: "0612345678"},
		SMSNotifer{PhoneNumber: "0711111111"}, // invalide
		PushNotifier{},
	}

	message := "Hello"

	for _, notifier := range notifiers {
		err := notifier.Send(message)
		if err != nil {
			fmt.Printf("[ERREUR - %s] %v\n", notifier.GetName(), err)
			continue
		}

		//si pas err
		notification := &Notification{
			Channel:   notifier.GetName(),
			Message:   message,
			TimeRef: time.Now(),
		}
		store.Add(notification)
	}

	fmt.Printf("résumé des notifications réussites")
	archived, _ := store.GetAll()
	if len(archived) == 0 {
		fmt.Println("aucune notification pour l'instant")
	} else {
		for _, n := range archived {
			fmt.Printf("ID: %d | Canal: %s | Message: '%s' | Date: %s\n",
			n.ID, n.Channel, n.Message, n.TimeRef.Format("2006-01-02 15:04:05"))

		}
	}
}



