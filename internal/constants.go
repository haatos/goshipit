package internal

type RoleID int

func (r RoleID) ToString() string {
	switch r {
	case Superuser:
		return "superuser"
	case ProMember:
		return "pro member"
	case Member:
		return "member"
	default:
		return "any user"
	}
}

const (
	Member    RoleID = 1
	ProMember RoleID = 10
	Superuser RoleID = 1000
)

const (
	CookieExpiresLayout = "2006-01-02T15:04:05Z"
	DBTimestampLayout   = "2006-01-02 15:04:05"
	SessionCookie       = "session"
	OAuthCookie         = "oauthstate"
	PasswordlessCookie  = "passwordless"
)

type Unit string

const (
	Kilos  Unit = "kg"
	Pounds Unit = "lbs"
)
