package account_store

import "context"

type AccountStore interface {
	Init(ctx context.Context) error
	CreateAppUser(ctx context.Context, login string, password string) (*AppUser, error)
	UpdateAppUser(ctx context.Context, user AppUser) error
	DeleteAppUserByID(ctx context.Context, id int64) error

	GetAppUserByID(ctx context.Context, id int64) (*AppUser, error)
	FindAppUser(ctx context.Context, filter AppUserFilter, skip int64, count int) ([]*AppUser, int64, error)
}

var accountStoreImpl AccountStore

func SetupAccountStoreImplementation(s AccountStore) {
	accountStoreImpl = s
}

func GetAccountStore() AccountStore {
	if accountStoreImpl == nil {
		panic("account store not initialized")
	}
	return accountStoreImpl
}
