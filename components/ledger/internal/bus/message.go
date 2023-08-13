package bus

import (
	"github.com/formancehq/ledger/internal"
	"github.com/formancehq/stack/libs/go-libs/metadata"
)

const (
	EventVersion = "v1"
	EventApp     = "ledger"

	EventTypeCommittedTransactions = "COMMITTED_TRANSACTIONS"
	EventTypeSavedMetadata         = "SAVED_METADATA"
	EventTypeRevertedTransaction   = "REVERTED_TRANSACTION"
)

type EventMessage struct {
	Date    ledger.Time `json:"date"`
	App     string      `json:"app"`
	Version string      `json:"version"`
	Type    string      `json:"type"`
	Payload any         `json:"payload"`
}

type CommittedTransactions struct {
	Ledger       string               `json:"ledger"`
	Transactions []ledger.Transaction `json:"transactions"`
}

func newEventCommittedTransactions(txs CommittedTransactions) EventMessage {
	return EventMessage{
		Date:    ledger.Now(),
		App:     EventApp,
		Version: EventVersion,
		Type:    EventTypeCommittedTransactions,
		Payload: txs,
	}
}

type SavedMetadata struct {
	Ledger     string            `json:"ledger"`
	TargetType string            `json:"targetType"`
	TargetID   string            `json:"targetId"`
	Metadata   metadata.Metadata `json:"metadata"`
}

func newEventSavedMetadata(metadata SavedMetadata) EventMessage {
	return EventMessage{
		Date:    ledger.Now(),
		App:     EventApp,
		Version: EventVersion,
		Type:    EventTypeSavedMetadata,
		Payload: metadata,
	}
}

type RevertedTransaction struct {
	Ledger              string             `json:"ledger"`
	RevertedTransaction ledger.Transaction `json:"revertedTransaction"`
	RevertTransaction   ledger.Transaction `json:"revertTransaction"`
}

func newEventRevertedTransaction(tx RevertedTransaction) EventMessage {
	return EventMessage{
		Date:    ledger.Now(),
		App:     EventApp,
		Version: EventVersion,
		Type:    EventTypeRevertedTransaction,
		Payload: tx,
	}
}
