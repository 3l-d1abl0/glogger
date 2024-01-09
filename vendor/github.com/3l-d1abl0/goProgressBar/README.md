# GoProgressbar

A simple progressbar for teminal app,
written in Golang.

### How to Use ?

1. Execute the following in your Project:
```
go get -u github.com/3l-d1abl0/goProgressBar/
```

2. Import the packge in your code
```
import github.com/3l-d1abl0/go-progressbar
```

3. Checkout the following Example:

```
func main() {

	var N int64 = 100
	var barSize int64 = 50
	var barSymbol string = "%"
	bar := goProgressBar.GetNewBar(N, 0, barSymbol, barSize)

	for i := 0; i <= int(N); i++ {

		time.Sleep(100 * time.Millisecond)
		bar.Display(int64(i))
	}

	bar.End()

}
    
```