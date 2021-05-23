package anfrage

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"io"
	"io/ioutil"
	"strings"
	"golang.org/x/net/html/charset"
)

const (
	GET  = "GET"
	POST = "POST"
	PUT  = "PUT"

	CreateRequestObjectException      = 10010
	GetResponseContentException       = 10011
	ReaderResponseContentException    = 10012
	PostDataJsonMarshalException      = 10013
	CreatePostRequestOjectException   = 10014
	SendPostRequestException          = 10015
	GetPostResponseContentException   = 10016
	CreatePutRequestObjectException   = 10017
	GetPutResponseCotnentException    = 10018
	ReaderPutResponseContentException = 10019
)

var (
	defaultClient = &http.Client{}
	Headers       = make(map[string]string)
	Cookies       = make(map[string]string)
)

type ErrorType int

type Error struct {
	TypeError ErrorType
	msg       string
}

func (e Error)Error() string {
	return e.msg
}

func newError(t ErrorType, msg string) Error {
	return Error{TypeError: t, msg: msg}
}

func setHeader(key, value string) {
	Headers[key] = value
}

func setHttpHeaderAndCookie(request *http.Request) {
	for k, v := range Headers {
		request.Header.Set(k, v)
	}

	for k, v := range Cookies {
		request.AddCookie(&http.Cookie{
			Name: k,
			Value: v,
		})
	}
}

func getWithClient(target string, client *http.Client) (string, error) {
	request, err := http.NewRequest(GET, target, nil)
	if err != nil {
		return "", newError(CreateRequestObjectException, "create request object exception : " + target)
	}

	setHttpHeaderAndCookie(request)

	response, err := client.Do(request)
	if err != nil {
		return "", newError(GetResponseContentException, "Can't Get response body exception : " + target)
	}

	defer response.Body.Close()

	utf8Reader, err := charset.NewReader(response.Body, response.Header.Get("Content-Type"))
	if err != nil {
		return "", err
	}

	content, err := ioutil.ReadAll(utf8Reader)

	if err != nil {
		return "", newError(ReaderResponseContentException, "Can't to read the response body")
	}
	return string(content), nil
}

func getBody(rawContent interface{}) (io.Reader, error) {
	var bodyReader io.Reader

	if rawContent != nil {
		switch body := rawContent.(type) {
		case map[string]string:
			jsonBody, err := json.Marshal(body)
			if err != nil {
				return nil, newError(PostDataJsonMarshalException, "Can't serialize map of string to json")
			}
			bodyReader = bytes.NewBuffer(jsonBody)
		case url.Values:
			bodyReader = strings.NewReader(body.Encode())
		case []byte:
			bodyReader = bytes.NewBuffer(body)
		case string:
			bodyReader = strings.NewReader(body)
		default:
			return nil, newError(PostDataJsonMarshalException, fmt.Sprintf("Can't handle body type %T", rawContent))
		}
	}
	return bodyReader, nil
}

func postWithClient(target, bodyType string, body interface{}, client *http.Client) (string, error) {
	bodyReader, err := getBody(body)
	if err != nil {
		return "", err
	}

	request, err := http.NewRequest(POST, target, bodyReader)
	setHeader("Content-Type", bodyType)
	setHttpHeaderAndCookie(request)

	if err != nil {
		return "", newError(CreatePostRequestOjectException, "Can't create post request object : " + target)
	}

	response, err := client.Do(request)
	if err != nil {
		return "", newError(SendPostRequestException, "Can't send http post request" + target)
	}

	defer response.Body.Close()

	content, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", newError(GetPostResponseContentException, "Cant' get http post response content : " +target)
	}
	return string(content), nil
}

func putWithClient(target, bodyType string, body interface{}, client *http.Client) (string, error) {
	bodyReader, err := getBody(body)
	if err != nil {
		return "", err
	}

	request, err := http.NewRequest(PUT, target, bodyReader)
	if err != nil {
		return "", newError(CreatePutRequestObjectException, "Can't create put request object : " + target)
	}

	setHeader("Content-Type", bodyType)
	setHttpHeaderAndCookie(request)

	response, err := client.Do(request)
	if err != nil {
		return "", newError(GetPutResponseCotnentException, "Can't get put request response content : " + target)
	}

	defer response.Body.Close()

	content, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", newError(ReaderPutResponseContentException, "Can't read put response content : " + target)
	}
	return string(content), nil
}

func Get(target string) (string, error) {
	return getWithClient(target, defaultClient)
}

func Post(target, bodyType string, body interface{}) (string, error) {
	return postWithClient(target, bodyType, body, defaultClient)
}

func PostForm(target string, data url.Values) (string, error) {
	return postWithClient(target, "application/x-www-form-urlencoded", data, defaultClient)
}

func Put(target, bodyType string, body interface{}) (string, error) {
	return putWithClient(target, bodyType, body, defaultClient)
}
