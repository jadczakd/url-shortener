package shortener

import (
	"errors"
	"time"

	errs "github.com/pkg/errors"
	"github.com/teris-io/shortid"
	"gopkg.in/dealancer/validate.v2"
	"log"
)

var (
	ErrRedirectNotFound = errors.New("Redirect not foud")
	ErrRedirectInvalid  = errors.New("Redirect invalid")
)

type redirectService struct {
	rr RedirectRepository
}

func NewRedirectService(repo RedirectRepository) RedirectService {
	return &redirectService{
		rr: repo,
	}
}

func (r *redirectService) Find(code string) (*Redirect, error) {
	return r.rr.Find(code)
}

func (r *redirectService) Store(redirect *Redirect) error {
	if err := validate.Validate(redirect); err != nil {
		log.Println(err)
		return errs.Wrap(ErrRedirectInvalid, "service.Redirect.Store")
	}
	redirect.Code = shortid.MustGenerate()
	redirect.CreatedAt = time.Now().Unix()
	return r.rr.Store(redirect)
}
