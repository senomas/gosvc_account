package account_store

import "github.com/senomas/gosvc_store/store"

type AppUser struct {
	Login    string
	Password string
	ID       int64
}

type AppUserFilter struct {
	Login    store.FilterString
	Password store.FilterString
	ID       store.FilterInt64
}
