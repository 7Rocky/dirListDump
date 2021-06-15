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

	"github.com/fatih/color"
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
		color.Red("\nFlag --url not set.\n")
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
		color.Red(err.Error() + "\n\n")
		os.Exit(1)
	}

	body, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()

	if err != nil {
		panic(err)
	}

	return body
}

var pattern = regexp.MustCompile(`<a href="(.*?)">`)

func getDirLine(depth int, isLast bool) string {
	if isLast {
		return color.CyanString(strings.Repeat("│   ", depth) + "└── ")
	}

	return color.CyanString(strings.Repeat("│   ", depth) + "├── ")
}

func getDirList(dirList chan string, url string, depth int, printTree bool) {
	body := getBodyUrl(url)
	matches := pattern.FindAllStringSubmatch(string(body), -1)

	if printTree && depth == 0 {
		fmt.Println("/")
	}

	for i, m := range matches {
		if m[1] == "/" {
			continue
		}

		if printTree {
			fmt.Println(getDirLine(depth, i == len(matches)-1) + m[1])
		}

		if strings.HasSuffix(m[1], "/") {
			getDirList(dirList, url+m[1], depth+1, printTree)
		} else {
			dirList <- url + m[1]
		}
	}

	if depth == 0 {
		close(dirList)
	}
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

	getDirList(dirList, *rootUrl, 0, *printTree)

	if *printTree {
		fmt.Println("")
	}

	wg.Wait()
	fmt.Printf("Done!! ")

	timeLapse := time.Since(start)

	if timeLapse.Seconds() < 3 {
		color.Green(timeLapse.String())
	} else if timeLapse.Seconds() < 20 {
		color.Yellow(timeLapse.String())
	} else {
		color.Red(timeLapse.String())
	}
}
