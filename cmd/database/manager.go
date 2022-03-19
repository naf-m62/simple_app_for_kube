package localstorage

type (
	ILocalStore interface {
	}

	localStore struct {
		user IUser
	}
)

func (ls *localStore) User() IUser {
	return ls.user
}
