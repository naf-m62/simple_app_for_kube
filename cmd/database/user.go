package localstorage

type (
	IUser interface {
		Get(userID int64) string
		Save(name string)
	}
	user struct {
	}
)





