&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp; ![alt text](https://3l-d1abl0.github.io/glogger/assets/Glogger.png)

~ A simple multi-link concurrent Downloader Written in Go

## 🚀Installation

**Make Sure you have Go(>=1.19) Installed.**

### Clone the repo and build binary
```
$ git clone https://github.com/3l-d1abl0/glogger.git
$ cd glogger/cmd/glogger/
```

##### Make binary
```
$ go build -o glogger
$ ./glogger -h
$ ./glogger -f "/path/to/links-file" -o "/path/to/output/folder/"

```

##### Alternative - run main.go
```
$ go run main.go -h
$ go run main.go -f "/path/to/links-file" -o "/path/to/output/folder/"
```


