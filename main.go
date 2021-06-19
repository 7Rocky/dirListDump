package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"regexp"
	"strings"
	"sync"
	"time"

	"io/ioutil"
	"net/http"
)

var (
	rootUrl   = flag.String("url", "", "The URL to list the files from. Directory listing must be available")
	dumpDir   = flag.String("dump", "dump", "Folder to dump contents")
	threads   = flag.Int("threads", 10, "Number of threads to use")
	printTree = flag.Bool("tree", true, "Show the tree file structure")
)

func parseParams() {
	flag.Parse()

	if *rootUrl == "" {
		RedPrint("\nFlag --url not set.\n")
		fmt.Printf("\nUse %s --help for more information\n\n", os.Args[0])
		os.Exit(1)
	}

	if !strings.HasSuffix(*rootUrl, "/") {
		*rootUrl += "/"
	}

	fmt.Println("")
}

func getBodyUrl(url string) []byte {
	res, err := http.Get(url)

	if err != nil {
		RedPrint(err.Error() + "\n\n")
		os.Exit(1)
	}

	body, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()

	if err != nil {
		panic(err)
	}

	return body
}

const dirPerm = os.FileMode(0755)
const filePerm = os.FileMode(0644)

var workDir string
var mu sync.Mutex

func writeFile(dir, file string, body []byte) {
	mu.Lock()
	defer mu.Unlock()

	os.MkdirAll(dir, dirPerm)
	os.Chdir(dir)

	if err := ioutil.WriteFile(file, body, filePerm); err != nil {
		panic(err)
	}

	os.Chdir(workDir)
}

var wg sync.WaitGroup

func dumpDirList(wg *sync.WaitGroup, dirList chan string, rootUrl string) {
	wg.Add(1)
	defer wg.Done()

	for url := range dirList {
		body := getBodyUrl(url)
		dir, file := path.Split(url[len(rootUrl):])

		writeFile(dir, file, body)
	}
}

var pattern = regexp.MustCompile(`<a href="(.*?)">`)
var tree = Tree{}

func getDirList(dirList chan string, url string, depth int) {
	body := getBodyUrl(url)
	matches := pattern.FindAllStringSubmatch(string(body), -1)

	for _, m := range matches {
		if m[1] == "/" {
			continue
		}

		tree.Add(url[len(*rootUrl):] + m[1])

		if strings.HasSuffix(m[1], "/") {
			getDirList(dirList, url+m[1], depth+1)
		} else {
			dirList <- url + m[1]
		}
	}

	if depth == 0 {
		close(dirList)
	}
}

func main() {
	parseParams()

	start := time.Now()

	os.MkdirAll(*dumpDir, dirPerm)
	os.Chdir(*dumpDir)
	workDir, _ = os.Getwd()

	dirList := make(chan string)

	for i := 0; i < *threads; i++ {
		go dumpDirList(&wg, dirList, *rootUrl)
	}

	getDirList(dirList, *rootUrl, 0)

	if *printTree {
		tree.Print()
		fmt.Println("")
	}

	wg.Wait()
	fmt.Printf("Done!! ")

	timeLapse := time.Since(start)

	if timeLapse.Seconds() < 3 {
		GreenPrint(timeLapse.String())
	} else if timeLapse.Seconds() < 20 {
		YellowPrint(timeLapse.String())
	} else {
		RedPrint(timeLapse.String())
	}
}
