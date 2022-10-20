package usecase

import (
	"context"

	"alukart32.com/bank/entity"
	"github.com/google/uuid"
)

type (
	AccountService interface {
		Create(ctx context.Context, a *entity.Account) error
		Get(ctx context.Context, id uuid.UUID) (*entity.Account, error)
		UpdateOwner(ctx context.Context, id uuid.UUID, owner string) (*entity.Account, error)
		AddBalance(ctx context.Context, id uuid.UUID, amount int64) (*entity.Account, error)
		Delete(ctx context.Context, id uuid.UUID) error
		ListEntries(ctx context.Context, id uuid.UUID) ([]entity.Entry, error)
		ListTransfers(ctx context.Context, id uuid.UUID) ([]entity.Transfer, error)
	}

	EntryService interface {
		Create(ctx context.Context, e *entity.Entry) (*entity.Entry, error)
		Get(ctx context.Context, id int64) (*entity.Entry, error)
		Update(ctx context.Context, e *entity.Entry) (*entity.Entry, error)
		List(ctx context.Context, accountId uuid.UUID) ([]entity.Entry, error)
		Delete(ctx context.Context, id int64) error
	}

	TransferService interface {
		Transfer(ctx context.Context, t entity.Transfer) (*entity.TransferRes, error)
		Get(ctx context.Context, id int64) (*entity.Transfer, error)
		List(ctx context.Context, params ListTransferParams) ([]entity.Transfer, error)
		Rollback(ctx context.Context, id int64) error
	}

	AccountRepo interface {
		Create(ctx context.Context, a *entity.Account) (*entity.Account, error)
		Get(ctx context.Context, id uuid.UUID) (*entity.Account, error)
		UpdateOwner(ctx context.Context, id uuid.UUID, owner string) (*entity.Account, error)
		UpdateBalance(ctx context.Context, id uuid.UUID, amount int64) (*entity.Account, error)
		Delete(ctx context.Context, id uuid.UUID) error
	}
	EntryRepo interface {
		Create(ctx context.Context, e *entity.Entry) (*entity.Entry, error)
		Get(ctx context.Context, id int64) (*entity.Entry, error)
		Update(ctx context.Context, e *entity.Entry) error
		List(ctx context.Context, accountId uuid.UUID) ([]entity.Entry, error)
		Delete(ctx context.Context, id int64) error
	}

	TransferRepo interface {
		Create(ctx context.Context, transfer *entity.Transfer) (*entity.TransferRes, error)
		Get(ctx context.Context, id int64) (*entity.Transfer, error)
		List(ctx context.Context, params ListTransferParams) ([]entity.Transfer, error)
		Rollback(ctx context.Context, id int64) error
	}

	PaggingParams struct {
		Limit  int32
		Offset int32
	}

	ListTransferMode byte

	ListTransferParams struct {
		FromAccountId uuid.UUID
		ToAccountId   uuid.UUID
		Mode          ListTransferMode
		PaggingParams
	}
)

const (
	ListFromAccount = iota
	ListToAccount
	ListByAccounts
)
