// Package xkcd is a simple wrapper around the https://xkcd.com JSON interface.
package xkcd

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	protocol = "https://"
	host     = "xkcd.com"
	api      = "/info.0.json"
	random   = "/random/comic"
)

// Comic is the parsed version of the JSON returned by the xkcd API.
type Comic struct {
	/*
	 * Num is the comic number.
	 * Day is the day the comic was published.
	 * Month is the month the comic was published.
	 * Year is the year the comic was published.
	 * Title is the comic title.
	 * SafeTitle is the same as Title (but safer?) ¯\_(ツ)_/¯.
	 * Transcript is the textual description of the comic
	 * Alt is the text content you seen when you hover over the comic image.
	 * Img is the URL to the comic image
	 * News is for announcements (not sure, it's usually empty).
	 * DateTime is the Day, Month and Year parsed into type Time
	 */
	Num        int    `json:"num"`
	Day        string `json:"day"`
	Month      string `json:"month"`
	Year       string `json:"year"`
	Title      string `json:"title"`
	SafeTitle  string `json:"safe_title"`
	Transcript string `json:"transcript"`
	Alt        string `json:"alt"`
	Img        string `json:"img"`
	News       string `json:"news"`
	DateTime   time.Time
}

// newReq returns a new Request given a method, URL, and optional body.
func newReq(method, url string) (req *http.Request, err error) {
	req, err = http.NewRequest(strings.ToUpper(method), url, nil)
	return
}

// disableRedirect returns a policy to a Request object to not follow
// redirections.
func disableRedirect(req *http.Request, via []*http.Request) error {
	return http.ErrUseLastResponse
}

// ParseComicResponse parses the JSON body of the HTTP response made to the
// xkcd API, and returns a Comic.
func ParseComicResponse(body []byte) (Comic, error) {
	var (
		err   error
		comic Comic
	)

	err = json.Unmarshal(body, &comic)
	if err != nil {
		return comic, err
	}

	year, err := strconv.Atoi(comic.Year)
	if err != nil {
		return comic, err
	}

	month, err := strconv.Atoi(comic.Month)
	if err != nil {
		return comic, err
	}

	day, err := strconv.Atoi(comic.Day)
	if err != nil {
		return comic, err
	}

	comic.DateTime = time.Date(
		year, time.Month(month), day, 0, 0, 0, 0, time.UTC,
	)
	return comic, err
}

// FetchRandomComicNum gets a random comic number form the xkcd API.
func FetchRandomComicNum() (int, error) {
	var (
		client = &http.Client{CheckRedirect: disableRedirect}
		url    = protocol + "c." + host + random
		resp   *http.Response
		num    int
		err    error
	)

	req, err := newReq("get", url)
	if err != nil {
		return num, err
	}

	resp, err = client.Do(req)
	if err != nil {
		return num, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 302 {
		num, err = strconv.Atoi(
			strings.Split(resp.Header["Location"][0], "/")[3],
		)
	}
	return num, err
}

// FetchComic makes a request to the xkcd HTTP API, parses the JSON response
// into a Comic and returns it.
func FetchComic(num int) (Comic, error) {
	var (
		client = &http.Client{CheckRedirect: disableRedirect}
		url    = protocol + host + "/" + strconv.Itoa(num) + api
		resp   *http.Response
		body   []byte
		comic  Comic
		err    error
	)

	if num < 1 {
		return comic, fmt.Errorf("invalid comic number (%d)\n", num)
	}

	req, err := newReq("get", url)
	if err != nil {
		return comic, err
	}

	resp, err = client.Do(req)
	if err != nil {
		return comic, err
	}
	defer func() {
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
	}()

	if resp.StatusCode == 200 {
		body, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return comic, err
		}
		comic, err = ParseComicResponse(body)
	}
	return comic, err
}

// FetchRandomComic fetches a random comic number from FetchRandomComicNumber
// and passes it on to FetchComic and returns a Comic.
func FetchRandomComic() (comic Comic, err error) {
	num, err := FetchRandomComicNum()
	if err != nil {
		return
	}

	return FetchComic(num)
}
