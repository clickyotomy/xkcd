package xkcd

import "fmt"

// ExampleFetchComic shows the usage of FetchComic.
func ExampleFetchComic() {
    comic, err := FetchComic(1024)

    if err != nil {
        fmt.Errorf("%s\n", err)
    }

    fmt.Printf(
        "[%s/%s/%s]: \"%s\"\n%s\n",
        comic.Month, comic.Day, comic.Year, comic.Title, comic.Img,
    )
    // Output:
    // [3/2/2012]: "Error Code"
    // https://imgs.xkcd.com/comics/error_code.png
}
