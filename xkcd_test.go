package xkcd_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/clickyotomy/xkcd"

	"github.com/davecgh/go-spew/spew"
)

var (
	// comicNum is the comic number to test.
	comicNum = 1024
	// testDir is the location of static test data.
	testDir = "./testdata"
	// testComic is the `Comic' to be fetched.
	testComic xkcd.Comic
	// testImg the image of the comic to be fetched.
	testImg []byte
	// curComic test data (in JSON) parsed into type `Comic'
	// (and tested again `testComic').
	curComic xkcd.Comic
	// curImg is the test image loaded into memory for
	// testing against `testImg'.
	curImg []byte
)

// loadTest loads the data for testing from `testDir'.
func loadTest(comic *xkcd.Comic, img *[]byte) (err error) {
	testJSON, err := ioutil.ReadFile(filepath.Join(testDir, "error_code.json"))
	if err != nil {
		fmt.Printf("Unable to read test data!\n")
		panic(err)
	}
	*comic, err = xkcd.ParseComicResponse(testJSON)
	if err != nil {
		fmt.Printf("Unable to parse data!\n")
		panic(err)
	}
	*img, err = ioutil.ReadFile(filepath.Join(testDir, "error_code.png"))
	if err != nil {
		fmt.Printf("Unable to read test image!\n")
		panic(err)
	}
	return nil
}

// fetchImg downloads the comic image.
func fetchImg(url string) (content []byte, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	content, err = ioutil.ReadAll(resp.Body)
	return
}

// TestInit initalized the test data.
func TestInit(t *testing.T) {
	loadTest(&testComic, &testImg)
}

// TestFetchComic tests the functionality of FetchComic.
func TestFetchComic(t *testing.T) {
	curComic, err := xkcd.FetchComic(comicNum)
	if err != nil {
		t.Errorf("Could not fetch comic.\nError: %+v\n", err)
	}

	if curComic.Num != comicNum {
		t.Fatalf(
			"Fetched comic number: `%d' is not the requested one: `%d'.\n",
			curComic.Num, comicNum,
		)
	}

	if !reflect.DeepEqual(testComic, curComic) {
		t.Fatalf(
			"Stored and downloaded comics are not the same.\n"+
				"Test Comic:\n%s\nDownladed Comic:\n%s\n",
			testComic.ToStr(), curComic.ToStr(),
		)

	}

	curImg, err := fetchImg(curComic.Img)
	if err != nil {
		t.Fatalf("Could not download the comic image.\nError: %+v\n", err)
	}

	if !reflect.DeepEqual(testImg, curImg) {
		scs := spew.ConfigState{Indent: "\t"}
		t.Fatalf(
			"Stored and downloaded images are not the same.\n"+
				"Test Image:\n%s\nDownladed Image:\n%s\n",
			scs.Sdump(testImg), scs.Sdump(curImg),
		)
	}
}

// TestFetchComicNotFound tests for a failure to fetch a comic.
func TestFetchComicNotFound(t *testing.T) {
	_, err := xkcd.FetchComic(99999)
	if err.Error() != "error: 404 Not Found" {
		t.Errorf(
			"Expected an HTTP status of `404 Not Found', but got %s.",
			err,
		)
	}
}

// TestFetchRandomComicNum tests the functionality of FetchRandomComicNum.
func TestFetchRandomComicNum(t *testing.T) {
	var num, err = xkcd.FetchRandomComicNum()

	if err != nil {
		t.Fatalf("Could not fetch a random comic.\nError: %+v\n", err)
	}

	if num < 0 {
		t.Fatalf("Fetched comic number: `%d' cannot be less than 0.\n",
			num,
		)
	}
}
