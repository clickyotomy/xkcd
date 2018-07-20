package xkcd

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

var (
	comicNum  = 1024         // comic number to test
	testDir   = "./testdata" // location of static test data
	testComic Comic          // the Comic to be fetched and tested
	testImg   []byte         // the image of the Comic to be tested
	curComic  Comic          // test data (in JSON) parsed into type Comic
	curImg    []byte         // test image loaded into memory for the test
)

// loadTest loads the data for testing from `testDir'.
func loadTest(comic *Comic, img *[]byte) (err error) {
	testJSON, err := ioutil.ReadFile(filepath.Join(testDir, "error_code.json"))
	if err != nil {
		fmt.Errorf("Unable to read test data!\n")
		panic(err)
	}
	*comic, err = ParseComicResponse(testJSON)
	if err != nil {
		fmt.Errorf("Unable to parse data!\n")
		panic(err)
	}
	*img, err = ioutil.ReadFile(filepath.Join(testDir, "error_code.png"))
	if err != nil {
		fmt.Errorf("Unable to read test image!\n")
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
	curComic, err := FetchComic(comicNum)
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
	_, err := FetchComic(99999)
	if err.Error() != "error: 404 Not Found" {
		t.Errorf(
			"Expected an HTTP status of `404 Not Found', but got %s.",
			err,
		)
	}
}

// TestFetchRandomComicNum tests the functionality of FetchRandomComicNum.
func TestFetchRandomComicNum(t *testing.T) {
	var num, err = FetchRandomComicNum()

	if err != nil {
		t.Fatalf("Could not fetch a random comic.\nError: %+v\n", err)
	}

	if num < 0 {
		t.Fatalf("Fetched comic number: `%d' cannot be less than 0.\n",
			num,
		)
	}
}
