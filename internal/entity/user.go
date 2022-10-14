package entity

type UserID int64

func NewUserID(userID int64) UserID {
	return UserID(userID)
}

func (u UserID) ToInt64() int64 {
	return (int64)(u)
}

// ----

type User struct {
	id              UserID
	defaultCurrency string
}

func NewUser(userID UserID) User {
	return User{
		id:              userID,
		defaultCurrency: "",
	}
}

func (u User) GetID() UserID {
	return u.id
}

func (u User) GetDefaultCurrency() string {
	return u.defaultCurrency
}

func (u *User) SetDefaultCurrency(currency string) {
	u.defaultCurrency = currency
}
