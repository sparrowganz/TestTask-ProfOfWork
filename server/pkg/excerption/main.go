package excerption

import (
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type Service interface {
	Get() (string, error)
}

type Excerption struct {
}

func New() Service {
	return &Excerption{}
}

func (e *Excerption) Get() (string, error) {

	r := rand.New(rand.NewSource(time.Now().Unix())).Intn(100) // initialize local pseudorandom generator

	resp, err := http.PostForm("http://api.forismatic.com/api/1.0/", url.Values{
		"method": []string{"getQuote"},
		"format": []string{"text"},
		"lang":   []string{"en"},
		"key":    []string{strconv.Itoa(r)},
	},
	)
	defer resp.Body.Close()

	bts, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", nil
	}
	return string(bts), nil
}
