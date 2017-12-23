package main

import (
    "fmt"
    "io/ioutil"
    "log"
    "os"
    "sort"
    "strings"
    "unicode/utf8"
)

const DIRECTORY = "/home/bryan/Documents/book_notes/"
const AUTHOR_PREFIX = "by "

func help() {
    fmt.Println("booknotes [command]")
    fmt.Println("")
    fmt.Println("Commands:")
    fmt.Println("list      Prints title, subtitle, and author for each book")
    fmt.Println("full      Same as \"list\" but includes file path")
    fmt.Println("authors   Prints all authors in alphabetical order")
    fmt.Println("titles    Prints all titles in alphabetical order")
    fmt.Println("help      Prints this help message")
    fmt.Println("")
}

type meta struct {
    title string
    subtitle string
    author string
    path string
}

func getMetaInfo() []meta {
    var info []meta

    files, err := ioutil.ReadDir(DIRECTORY)
    if err != nil {
        log.Fatal(err)
    }

    for _, f := range files {
        if strings.HasSuffix(f.Name(), ".txt") {
            path := DIRECTORY + f.Name()
            byteContents, fileErr := ioutil.ReadFile(path)

            if fileErr != nil {
                log.Fatal(fileErr)
            } else {
                contents := string(byteContents)
                lines := strings.Split(contents, "\n")

                title := lines[1]
                subtitle := lines[2]
                author := lines[2]

                if strings.HasPrefix(lines[2], AUTHOR_PREFIX) {
                    author = lines[2]
                    subtitle = ""
                } else {
                    subtitle = lines[2]
                    author = lines[3]
                }

                author = strings.TrimPrefix(author, AUTHOR_PREFIX)
                info = append(info, meta{title, subtitle, author, path})
            }
        }
    }

    return info
}

func list(full bool) {
    for _, meta := range getMetaInfo() {
        fmt.Println(meta.title)

        if utf8.RuneCountInString(meta.subtitle) > 0 {
            fmt.Println(meta.subtitle)
        }

        fmt.Println("by", meta.author)

        if full {
            fmt.Println(meta.path)
        }

        fmt.Println("")
    }
}

func printSortedProperties(property string) {
    resultMap := make(map[string]bool)

    for _, meta := range getMetaInfo() {
        switch property {
        case "author":
            resultMap[meta.author] = true
        case "title":
            resultMap[meta.title] = true
        }
    }

    results := make([]string, 0, len(resultMap))
    for key := range resultMap {
        results = append(results, key)
    }

    sort.Strings(results)

    for _, result := range results {
        fmt.Println(result)
    }
}

func main() {
    args := os.Args[1:]

    if len(args) > 0 {
        switch args[0] {
        case "list":
            list(false)
        case "full":
            list(true)
        case "authors":
            printSortedProperties("author")
        case "titles":
            printSortedProperties("title")
        case "help":
            help()
        }
    } else {
        help()
    }
}
