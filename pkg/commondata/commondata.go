// commondata.go
package commondata

type UrlObject struct {
	Url      string
	Filename string
	Size     int64
}

// A type for the target Urls
type TargetUrls struct {
	ValidUrls   []UrlObject
	InvalidUrls []UrlObject
}

type Message struct {
	Idx  int
	Size int64
}
