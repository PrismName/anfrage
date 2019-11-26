package anfrage

import (
	"golang.org/net/x/publicsuffix"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
)

const (
	GET    = "GET"
	POST   = "POST"
	PUT    = "PUT"
	DELETE = "DELETE"
)

type Response *http.Response

type Request *http.Request

type Anfrage struct {
	Url       string
	Method    string
	Data      map[string]interface{}
	FormData  url.Values
	QueryData url.Values
	Headers   http.Header
	Cookies   []*http.Cookie
	UserAgent string
	BasicAuth struct{ Username, Password string }
	Proxy     func(proxyUrl string) (string, error)
	IsClear   bool
}

func (a *Anfrage) NewAnfrage() {
	cookieJarOptions := cookiejar.Options{
		PublicSuffixList: publicsuffix.List,
	}

	jar, _ := cookiejar.New(&cookieJarOptions)

	a := &Anfrage{
		Data:      make(map[string]interface{}),
		FormData:  http.Values{},
		QueryData: http.Values{},
		Headers:   http.Header{},
		Cookies:   make([]*http.Cookie, 0),
		UserAgent: "",
		BasicAuth: struct{ Username, Password string }{},
	}

	return a
}

func (a *Anfrage) ClearAnfrage() {
	if a.IsClear {
		return
	}
	a.Url = ""
	a.Method = ""
	a.Headers = http.Header{}
	a.Data = make(map[string]interface{})
	a.FormData = http.Values{}
	a.QueryData = http.Values{}
	a.Cookies = make([]*http.Cookie, 0)
	a.UserAgent = ""
}

func (a *Anfrage) Get(rawUrl string) *Anfrage {
	a.ClearAnfrage()
	a.Method = GET
	a.Url = rawUrl
	return a
}

func (a *Anfrage) Post(rawUrl string) *Anfrage {
	a.ClearAnfrage()
	a.Method = POST
	a.Url = rawUrl
	return a
}

func (a *Anfrage) Delete(rawUrl string) *Anfrage {
	a.ClearAnfrage()
	a.Method = DELETE
	a.Url = rawUrl
	return a
}

func (a *Anfrage) Put(rawUrl string) *Anfrage {
	a.ClearAnfrage()
	a.Method = PUT
	a.Url = rawUrl
	return a
}

func (a *Anfrage) UrlOpen(method, rawUrl string) *Anfrage {
	switch method {
	case GET:
		return a.Get(rawUrl)
	case POST:
		return a.Post(rawUrl)
	case PUT:
		return a.Put(rawUrl)
	case DELETE:
		return a.Delete(rawUrl)
	default:
		a.Method = method
		a.Url = rawUrl
		return a
	}
}

func (a *Anfrage) SetHeaders(key, value string) *Anfrage {
	a.Headers.Add(key, value)
	return a
}

func (a *Anfrage) SetBasicAuth(username, password string) *Anfrage {
	a.BasicAuth = struct{ Username, Password string }{username, password}
	return a
}

func (a *Anfrage) SetCookie(c *http.Cookie) *Anfrage {
	s.Cookies = append(s.Cookies, c)
	return a
}

func (a *Anfrage) SetProxies(proxyUrl string) *Anfrage {
	proxiesUrl, err := url.Parse(proxyUrl)

	if err != nil {
		log.Fatal("Can't parser proxy url : ", err)
	}

	a.Proxy = http.ProxyURL(proxiesUrl)
	return a
}
