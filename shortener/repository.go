package shortener

// RedirectRepository used to store and find the entities
type RedirectRepository interface {
	Find(code string) (*Redirect, error)
	Store(redirect *Redirect) error
}
