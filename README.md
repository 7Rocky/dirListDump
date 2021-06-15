# dirListDump

`dirListDump` is a tool to download all the files from a website that has directory listing enabled.

### Usage


### Running the tool

First step, clone or download this repository:

```console
git clone https://github.com/7Rocky/dirListDump
cd dirListDump

Then, you could run the Go source code using `go run` as shown below:

```console
go run dirListDump.go --url <the-url>
```

The preferred way is tu build the source code into a binary executable file. You might be tempted to use `go build dirListDump.go`, but it is not recommended because of the huge file it produces. Instead, it is possible to use:

```console
go build --ldflags='-s -w' dirListDump.go
upx --ultra-brute dirListDump
```

**Note:** You need to have `upx` installed in your machine.

The above process can last a few minutes, but generates a lighter file. Then, you can execute:

```console
./dirListDump --url <the-url>
```

You can also include the binary executable file at `/usr/local/bin` or other similar path to be ablle to run it from every directory on your machine just as 

```console
dirListDump --url <the-url>
```

Hope it is useful! :smile:
