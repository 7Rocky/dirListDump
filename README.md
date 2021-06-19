# dirListDump

`dldump` is a tool to download all the files from a website that has directory listing enabled (useful in CTF for extracting files from a web server).

### Usage

The parameters that can be used in this tool are:

| Parameter | Default | Type             | Meaning                                       |
| --------- | ------- | ---------------- | --------------------------------------------- |
| --url     | -       | URL (http://...) | The URL to list the files from (required)     |
| --threads | 10      | integer          | Number of threads to use                      |
| --dump    | dump    | string           | Folder to dump contents                       |
| --tree    | true    | boolean          | Show the tree file structure                  |
| --help    | -       | -                | Display the above parameters and explanations |

### Running the tool

First step, clone or download this repository:

```console
git clone https://github.com/7Rocky/dirListDump
cd dirListDump
```

Then, you could run the Go source code using `go run` as shown below:

```console
go run *.go --url <the-url>
```

The preferred way is to build the source code into a binary executable file.

You might be tempted to use `go build -o dldump *.go`, but it is not recommended because of the huge file it produces. Instead, it is possible to use:

```console
go build --ldflags='-s -w' -o dldump *.go
upx --ultra-brute dldump
```

**Note:** You need to have `upx` installed in your machine.

The above process can last a few minutes, but generates a lighter file. Then, you can execute:

```console
./dldump --url <the-url>
```

You can also include the binary executable file at `/usr/local/bin` or other similar path to be able to run it from every directory on your machine just as:

```console
dldump --url <the-url>
```

Hope it is useful! :smile:
