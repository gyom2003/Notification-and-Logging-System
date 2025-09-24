package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"
)


type Notification struct {
	ID        int
	Channel   string
	Target    string
	Message   string
	Timestamp time.Time
}
type Notifier interface {
	Send(target string, message string) error
	GetName() string
}

// EmailNotifier
type EmailNotifier struct{}

func (e EmailNotifier) Send(target string, message string) error {
	if !strings.Contains(target, "@") {
		return errors.New("❌ Adresse email invalide")
	}
	fmt.Printf("[EMAIL] Envoi de '%s' à %s réussi ✅\n", message, target)
	return nil
}
func (e EmailNotifier) GetName() string { return "Email" }

// SMSNotifier
type SMSNotifier struct{}

func (s SMSNotifier) Send(target string, message string) error {
	if len(target) < 2 || target[:2] != "06" {
		return errors.New("❌ Numéro invalide : doit commencer par '06'")
	}
	fmt.Printf("[SMS] Envoi de '%s' au %s réussi ✅\n", message, target)
	return nil
}
func (s SMSNotifier) GetName() string { return "SMS" }

type Storer interface {
	Add(n *Notification) error
	GetAll() ([]*Notification, error)
}

type MemoryStore struct {
	notifications []*Notification
	nextID        int
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		notifications: []*Notification{},
		nextID:        1,
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
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("=== Simulateur de Notifications ===")

	for {
		fmt.Println("\n--- Menu ---")
		fmt.Println("1. Créer une notification")
		fmt.Println("2. Voir l'historique")
		fmt.Println("3. Quitter")
		fmt.Print("Votre choix: ")

		choice := readLine(reader)

		switch choice {
		case "1":
			handleCreateNotification(reader, store)
		case "2":
			handleListNotifications(store)
		case "3":
			fmt.Println("👋 Au revoir!")
			return
		default:
			fmt.Println("❌ Choix invalide, réessayez.")
		}
	}
}

func handleCreateNotification(reader *bufio.Reader, store Storer) {
	fmt.Print("Choisissez le type (sms/email): ")
	ntype := strings.ToLower(readLine(reader))

	var notifier Notifier
	switch ntype {
	case "sms":
		notifier = SMSNotifier{}
	case "email":
		notifier = EmailNotifier{}
	default:
		fmt.Println("type non supporté")
		return
	}

	fmt.Print("Entrez le destinataire: ")
	target := readLine(reader)

	fmt.Print("Entrez le message: ")
	message := readLine(reader)

	err := notifier.Send(target, message)
	if err != nil {
		fmt.Printf("[ERREUR - %s] %v\n", notifier.GetName(), err)
		return
	}

	// Archivage
	notification := &Notification{
		Channel:   notifier.GetName(),
		Target:    target,
		Message:   message,
		Timestamp: time.Now(),
	}
	store.Add(notification)

	fmt.Println("✅ Notification envoyée et archivée.")
}

func handleListNotifications(store Storer) {
	archived, _ := store.GetAll()
	if len(archived) == 0 {
		fmt.Println("Aucune notification envoyée avec succès.")
		return
	}

	fmt.Println("\n=== Historique des envois réussis ===")
	for _, n := range archived {
		fmt.Printf("ID: %d | Canal: %s | Destinataire: %s | Message: '%s' | Date: %s\n",
			n.ID, n.Channel, n.Target, n.Message, n.Timestamp.Format("2006-01-02 15:04:05"))
	}
}

func readLine(reader *bufio.Reader) string {
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}
