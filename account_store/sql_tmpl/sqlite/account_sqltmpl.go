package sqlite

import (
	"fmt"
	"strings"

	"github.com/senomas/gosvc_account/account_store"
	"github.com/senomas/gosvc_account/account_store/sql_tmpl"
	"github.com/senomas/gosvc_store/store"
	"github.com/senomas/gosvc_store/store/sql_tmpl/sqlite"
)

type AccountStoreTemplateImpl struct{}

func init() {
	sql_tmpl.SetupAccountStoreTemplate(&AccountStoreTemplateImpl{})
}

// InsertAppUser implements sql_tmpl.AccountStoreTemplate.
func (s *AccountStoreTemplateImpl) InsertAppUser(t *account_store.AppUser) (string, []any) {
	return `INSERT INTO app_user (login, password) VALUES ($1, $2)`, []any{t.Login, t.Password}
}

// UpdateAppUser implements sql_tmpl.AccountStoreTemplate.
func (s *AccountStoreTemplateImpl) UpdateAppUser(t *account_store.AppUser) (string, []any) {
	return `UPDATE app_user SET login = $1, password = $2 WHERE id = $3`, []any{t.Login, t.Password, t.ID}
}

// DeleteAppUserByID implements sql_tmpl.AccountStoreTemplate.
func (s *AccountStoreTemplateImpl) DeleteAppUserByID(id any) (string, []any) {
	return `DELETE FROM app_user WHERE id = $1`, []any{id}
}

// GetAppUserByID implements sql_tmpl.AccountStoreTemplate.
func (s *AccountStoreTemplateImpl) GetAppUserByID(id any) (string, []any) {
	return `SELECT id, login, password FROM app_user WHERE id = $1`, []any{id}
}

func (s *AccountStoreTemplateImpl) findAppUserWhere(filter account_store.AppUserFilter) ([]string, []any) {
	where := []string{}
	args := []any{}

	where, args = sqlite.FilterToString(where, args, "id", filter.ID)
	where, args = sqlite.FilterToString(where, args, "login", filter.Login)
	where, args = sqlite.FilterToString(where, args, "password", filter.Password)

	return where, args
}

// FindAppUser implements sql_tmpl.AccountStoreTemplate.
func (s *AccountStoreTemplateImpl) FindAppUser(filter account_store.AppUserFilter, skip int64, limit int) (string, []any) {
	where, args := s.findAppUserWhere(filter)
	sl := ""
	if limit > 0 {
		sl += fmt.Sprintf(" LIMIT %d", limit)
	} else {
		sl += " LIMIT 1000"
	}
	if skip > 0 {
		sl += fmt.Sprintf(" OFFSET %d", skip)
	}
	if len(where) > 0 {
		return `SELECT id, login, password FROM app_user WHERE ` + strings.Join(where, " AND ") + sl, args
	}
	return `SELECT id, login, password FROM app_user` + sl, args
}

// FindAppUserTotal implements sql_tmpl.AccountStoreTemplate.
func (s *AccountStoreTemplateImpl) FindAppUserTotal(filter account_store.AppUserFilter) (string, []any) {
	where, args := s.findAppUserWhere(filter)
	if len(where) > 0 {
		return `SELECT COUNT(*) FROM app_user WHERE ` + strings.Join(where, " AND "), args
	}
	return `SELECT COUNT(*) FROM app_user`, args
}

// ErrorMapFind implements sql_tmpl.AccountStoreTemplate.
func (*AccountStoreTemplateImpl) ErrorMapFind(err error) error {
	if err.Error() == "sql: no rows in result set" {
		return store.ErrNoData
	}
	return err
}
