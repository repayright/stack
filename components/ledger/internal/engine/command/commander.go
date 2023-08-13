package command

import (
	"context"
	"fmt"
	"math/big"
	"sync"

	ledger "github.com/formancehq/ledger/internal"
	"github.com/formancehq/ledger/internal/engine/utils/batching"
	"github.com/formancehq/ledger/internal/machine"
	"github.com/formancehq/ledger/internal/machine/vm"
	storageerrors "github.com/formancehq/ledger/internal/storage"
	"github.com/formancehq/stack/libs/go-libs/collectionutils"
	"github.com/formancehq/stack/libs/go-libs/errorsutil"
	"github.com/formancehq/stack/libs/go-libs/metadata"
	"github.com/pkg/errors"
)

type Parameters struct {
	DryRun         bool
	IdempotencyKey string
}

type Commander struct {
	*batching.Batcher[*ledger.ChainedLog]
	store      Store
	locker     Locker
	compiler   *Compiler
	running    sync.WaitGroup
	lastTXID   *big.Int
	referencer *Referencer
	mu         sync.Mutex

	lastLog *ledger.ChainedLog
}

func New(
	store Store,
	locker Locker,
	compiler *Compiler,
	referencer *Referencer,
) *Commander {
	log, err := store.ReadLastLogWithType(context.Background(), ledger.NewTransactionLogType, ledger.RevertedTransactionLogType)
	if err != nil && !storageerrors.IsNotFoundError(err) {
		panic(err)
	}

	lastTxID := big.NewInt(-1)
	if err == nil {
		switch payload := log.Data.(type) {
		case ledger.NewTransactionLogPayload:
			lastTxID = payload.Transaction.ID
		case ledger.RevertedTransactionLogPayload:
			lastTxID = payload.RevertTransaction.ID
		default:
			panic(fmt.Sprintf("unhandled payload type: %T", payload))
		}
	} else {

	}

	lastLog, err := store.GetLastLog(context.Background())
	if err != nil && !storageerrors.IsNotFoundError(err) {
		panic(err)
	}

	return &Commander{
		store:      store,
		locker:     locker,
		compiler:   compiler,
		referencer: referencer,
		lastTXID:   lastTxID,
		lastLog:    lastLog,
		Batcher:    batching.NewBatcher(store.InsertLogs, 1, 4096),
	}
}

func (commander *Commander) GetLedgerStore() Store {
	return commander.store
}

func (commander *Commander) exec(ctx context.Context, parameters Parameters, script ledger.RunScript,
	logComputer func(tx *ledger.Transaction, accountMetadata map[string]metadata.Metadata) *ledger.Log) (*ledger.ChainedLog, error) {

	if script.Script.Plain == "" {
		return nil, ErrNoScript
	}

	if script.Timestamp.IsZero() {
		script.Timestamp = ledger.Now()
	}

	execContext := newExecutionContext(commander, parameters)
	return execContext.run(ctx, func(executionContext *executionContext) (*ledger.ChainedLog, chan struct{}, error) {
		if script.Reference != "" {
			if err := commander.referencer.take(referenceTxReference, script.Reference); err != nil {
				return nil, nil, ErrConflictError
			}
			defer commander.referencer.release(referenceTxReference, script.Reference)

			_, err := commander.store.GetTransactionByReference(ctx, script.Reference)
			if err == nil {
				return nil, nil, ErrConflictError
			}
			if err != storageerrors.ErrNotFound && err != nil {
				return nil, nil, err
			}
		}

		program, err := commander.compiler.Compile(ctx, script.Plain)
		if err != nil {
			return nil, nil, errorsutil.NewError(ErrCompilationFailed, errors.Wrap(err, "compiling numscript"))
		}

		m := vm.NewMachine(*program)

		if err := m.SetVarsFromJSON(script.Vars); err != nil {
			return nil, nil, errorsutil.NewError(ErrCompilationFailed,
				errors.Wrap(err, "could not set variables"))
		}

		involvedAccounts, involvedSources, err := m.ResolveResources(ctx, commander.store)
		if err != nil {
			return nil, nil, errorsutil.NewError(ErrCompilationFailed,
				errors.Wrap(err, "could not resolve program resources"))
		}

		worldFilter := collectionutils.FilterNot(collectionutils.FilterEq("world"))
		lockAccounts := Accounts{
			Read:  collectionutils.Filter(involvedAccounts, worldFilter),
			Write: collectionutils.Filter(involvedSources, worldFilter),
		}

		unlock, err := commander.locker.Lock(ctx, lockAccounts)
		if err != nil {
			return nil, nil, errors.Wrap(err, "locking accounts for tx processing")
		}
		unlock(ctx)

		err = m.ResolveBalances(ctx, commander.store)
		if err != nil {
			return nil, nil, errorsutil.NewError(ErrCompilationFailed,
				errors.Wrap(err, "could not resolve balances"))
		}

		result, err := machine.Run(m, script)
		if err != nil {
			return nil, nil, errors.Wrap(err, "running numscript")
		}

		if len(result.Postings) == 0 {
			return nil, nil, ErrNoPostings
		}

		tx := ledger.NewTransaction().
			WithPostings(result.Postings...).
			WithMetadata(result.Metadata).
			WithDate(script.Timestamp).
			WithID(commander.nextTXID()).
			WithReference(script.Reference)

		log := logComputer(tx, result.AccountMetadata)
		if parameters.IdempotencyKey != "" {
			log = log.WithIdempotencyKey(parameters.IdempotencyKey)
		}

		return executionContext.AppendLog(ctx, log)
	})
}

func (commander *Commander) CreateTransaction(ctx context.Context, parameters Parameters, script ledger.RunScript) (*ledger.Transaction, error) {
	log, err := commander.exec(ctx, parameters, script, ledger.NewTransactionLog)
	if err != nil {
		return nil, err
	}
	return log.Data.(ledger.NewTransactionLogPayload).Transaction, nil
}

func (commander *Commander) SaveMeta(ctx context.Context, parameters Parameters, targetType string, targetID interface{}, m metadata.Metadata) error {
	if m == nil {
		return nil
	}

	if targetType == "" {
		return errorsutil.NewError(ErrValidation, errors.New("empty target type"))
	}
	if targetID == "" {
		return errorsutil.NewError(ErrValidation, errors.New("empty target id"))
	}

	execContext := newExecutionContext(commander, parameters)
	_, err := execContext.run(ctx, func(executionContext *executionContext) (*ledger.ChainedLog, chan struct{}, error) {
		var (
			log *ledger.Log
			at  = ledger.Now()
		)
		switch targetType {
		case ledger.MetaTargetTypeTransaction:
			_, err := commander.store.GetTransaction(ctx, targetID.(uint64))
			if err != nil {
				return nil, nil, err
			}
			log = ledger.NewSetMetadataLog(at, ledger.SetMetadataLogPayload{
				TargetType: ledger.MetaTargetTypeTransaction,
				TargetID:   targetID.(uint64),
				Metadata:   m,
			})
		case ledger.MetaTargetTypeAccount:
			log = ledger.NewSetMetadataLog(at, ledger.SetMetadataLogPayload{
				TargetType: ledger.MetaTargetTypeAccount,
				TargetID:   targetID.(string),
				Metadata:   m,
			})
		default:
			return nil, nil, errorsutil.NewError(ErrValidation, errors.Errorf("unknown target type '%s'", targetType))
		}

		return executionContext.AppendLog(ctx, log)
	})
	return err
}

func (commander *Commander) RevertTransaction(ctx context.Context, parameters Parameters, id uint64) (*ledger.Transaction, error) {

	if err := commander.referencer.take(referenceReverts, id); err != nil {
		return nil, ErrRevertOccurring
	}
	defer commander.referencer.release(referenceReverts, id)

	tx, err := commander.store.GetTransaction(ctx, id)
	if err != nil {
		if errors.Is(err, storageerrors.ErrNotFound) {
			return nil, errors.New("tx not found")
		}
		return nil, err
	}
	if tx.Reverted {
		return nil, ErrAlreadyReverted
	}

	transactionToRevert, err := commander.store.GetTransaction(ctx, id)
	if storageerrors.IsNotFoundError(err) {
		return nil, errorsutil.NewError(err, errors.Errorf("transaction %d not found", id))
	}
	if err != nil {
		return nil, err
	}

	rt := transactionToRevert.Reverse()
	rt.Metadata = ledger.MarkReverts(metadata.Metadata{}, transactionToRevert.ID)

	log, err := commander.exec(ctx, parameters,
		ledger.TxToScriptData(ledger.TransactionData{
			Postings: rt.Postings,
			Metadata: rt.Metadata,
		}),
		func(tx *ledger.Transaction, accountMetadata map[string]metadata.Metadata) *ledger.Log {
			return ledger.NewRevertedTransactionLog(tx.Date, transactionToRevert.ID, tx)
		})
	if err != nil {
		return nil, err
	}

	return log.Data.(ledger.RevertedTransactionLogPayload).RevertTransaction, nil
}

func (commander *Commander) Close() {
	commander.Batcher.Close()
	commander.running.Wait()
}

func (commander *Commander) chainLog(log *ledger.Log) *ledger.ChainedLog {
	commander.mu.Lock()
	defer commander.mu.Unlock()

	commander.lastLog = log.ChainLog(commander.lastLog)
	return commander.lastLog
}

func (commander *Commander) nextTXID() *big.Int {
	commander.mu.Lock()
	defer commander.mu.Unlock()

	ret := big.NewInt(0).Add(commander.lastTXID, big.NewInt(1))
	commander.lastTXID = ret

	return ret
}
