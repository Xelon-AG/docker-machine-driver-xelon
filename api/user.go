package api

// User represents a Xelon user.
type User struct {
	APIToken         string `json:"api_token"`
	FirstName        string `json:"firstname"`
	ID               int    `json:"id"`
	Surname          string `json:"surname"`
	TenantIdentifier string `json:"tenantIdentifier"`
}

type userRoot struct {
	User User `json:"user"`
}
