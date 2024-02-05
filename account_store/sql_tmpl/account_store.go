package sql_tmpl

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	"github.com/senomas/gosvc_account/account_store"
	"github.com/senomas/gosvc_store/store"
)

var (
	errCtxNoDB             = errors.New("no db defined in context")
	errNoaccount_storeTmpl = errors.New("user account_store template not initialized")
)

type AccountStoreImpl struct{}

type AccountStoreTemplate interface {
	InsertAppUser(t *account_store.AppUser) (string, []any)
	UpdateAppUser(t *account_store.AppUser) (string, []any)
	DeleteAppUserByID(id any) (string, []any)

	GetAppUserByID(id any) (string, []any)
	FindAppUser(account_store.AppUserFilter, int64, int) (string, []any)
	FindAppUserTotal(account_store.AppUserFilter) (string, []any)

	ErrorMapFind(error) error
}

func init() {
	slog.Debug("Register sql_tmpl.AppUseraccount_store")
	account_store.SetupAccountStoreImplementation(&AccountStoreImpl{})
}

var useraccount_storeTemplateImpl AccountStoreTemplate

func SetupAccountStoreTemplate(t AccountStoreTemplate) {
	useraccount_storeTemplateImpl = t
}

func (t *AccountStoreImpl) Init(ctx context.Context) error {
	if useraccount_storeTemplateImpl == nil {
		return errNoaccount_storeTmpl
	}
	return nil
}

// CreateAppUser implements account_store.AppUseraccount_store.
func (t *AccountStoreImpl) CreateAppUser(ctx context.Context, login string, password string) (*account_store.AppUser, error) {
	if db, ok := ctx.Value(store.StoreCtxDB).(*sql.DB); ok {
		user := account_store.AppUser{Login: login, Password: password}
		qry, args := useraccount_storeTemplateImpl.InsertAppUser(&user)
		slog.Debug("CreateAppUser", "qry", qry, "args", &store.JsonLogValue{V: args})
		rs, err := db.ExecContext(ctx, qry, args...)
		if err != nil {
			slog.Warn("Error insert user", "qry", qry, "error", err)
			return nil, err
		}
		user.ID, err = rs.LastInsertId()
		return &user, err
	}
	return nil, errCtxNoDB
}

// UpdateAppUser implements account_store.AppUseraccount_store.
func (t *AccountStoreImpl) UpdateAppUser(ctx context.Context, user account_store.AppUser) error {
	if db, ok := ctx.Value(store.StoreCtxDB).(*sql.DB); ok {
		qry, args := useraccount_storeTemplateImpl.UpdateAppUser(&user)
		slog.Debug("UpdateAppUser", "qry", qry, "args", &store.JsonLogValue{V: args})
		_, err := db.ExecContext(ctx, qry, args...)
		return err
	}
	return errCtxNoDB
}

// DeleteAppUserByID implements account_store.AppUseraccount_store.
func (t *AccountStoreImpl) DeleteAppUserByID(ctx context.Context, id int64) error {
	if db, ok := ctx.Value(store.StoreCtxDB).(*sql.DB); ok {
		qry, args := useraccount_storeTemplateImpl.DeleteAppUserByID(id)
		slog.Debug("DeleteAppUserByID", "qry", qry, "args", &store.JsonLogValue{V: args})
		_, err := db.ExecContext(ctx, qry, args...)
		return err
	}
	return errCtxNoDB
}

// GetAppUserByID implements account_store.AppUseraccount_store.
func (t *AccountStoreImpl) GetAppUserByID(ctx context.Context, id int64) (*account_store.AppUser, error) {
	if db, ok := ctx.Value(store.StoreCtxDB).(*sql.DB); ok {
		user := account_store.AppUser{}
		qry, args := useraccount_storeTemplateImpl.GetAppUserByID(id)
		slog.Debug("GetAppUserByID", "qry", qry, "args", &store.JsonLogValue{V: args})
		err := db.QueryRowContext(ctx, qry, args...).Scan(&user.ID, &user.Login, &user.Password)
		if err != nil {
			err = useraccount_storeTemplateImpl.ErrorMapFind(err)
		}
		return &user, err
	}
	return nil, errCtxNoDB
}

// FindAppUser implements account_store.AppUseraccount_store.
func (*AccountStoreImpl) FindAppUser(ctx context.Context, filter account_store.AppUserFilter, skip int64, count int) ([]*account_store.AppUser, int64, error) {
	if db, ok := ctx.Value(store.StoreCtxDB).(*sql.DB); ok {
		total := int64(0)
		qry, args := useraccount_storeTemplateImpl.FindAppUserTotal(filter)
		slog.Debug("FindAppUserTotal", "qry", qry, "args", &store.JsonLogValue{V: args})
		err := db.QueryRowContext(ctx, qry, args...).Scan(&total)
		if err != nil {
			err = useraccount_storeTemplateImpl.ErrorMapFind(err)
			return nil, total, err
		}
		qry, args = useraccount_storeTemplateImpl.FindAppUser(filter, skip, count)
		slog.Debug("FindAppUser", "qry", qry, "args", &store.JsonLogValue{V: args})
		rows, err := db.QueryContext(ctx, qry, args...)
		if err != nil {
			err = useraccount_storeTemplateImpl.ErrorMapFind(err)
			return nil, total, err
		}
		defer rows.Close()
		users := []*account_store.AppUser{}
		for rows.Next() {
			user := account_store.AppUser{}
			err = rows.Scan(&user.ID, &user.Login, &user.Password)
			if err != nil {
				err = useraccount_storeTemplateImpl.ErrorMapFind(err)
				return nil, total, err
			}
			users = append(users, &user)
		}
		return users, total, nil
	}
	return nil, 0, errCtxNoDB
}
