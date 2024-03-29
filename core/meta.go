package core

import (
    "github.com/bgrohman/booknotes/config"
    "github.com/bgrohman/booknotes/printing"
    "fmt"
    "io/ioutil"
    "log"
    "math/rand"
    "sort"
    "strings"
    "time"
    "unicode/utf8"
)

type Meta struct {
    Title string
    Subtitle string
    Author string
    Path string
}

func GetMetaInfo(whichFile string) []Meta {
    var info []Meta

    files, err := ioutil.ReadDir(config.GetConfig().Directory)
    if err != nil {
        log.Fatal(err)
    }

    for _, f := range files {
        if strings.HasSuffix(f.Name(), ".txt") {
            path := config.GetConfig().Directory + f.Name()

            if len(whichFile) < 1 || whichFile == path {
                byteContents, fileErr := ioutil.ReadFile(path)

                if fileErr != nil {
                    log.Fatal(fileErr)
                } else {
                    contents := string(byteContents)
                    lines := strings.Split(contents, "\n")

                    title := lines[1]
                    subtitle := lines[2]
                    author := lines[2]

                    if strings.HasPrefix(lines[2], config.GetConfig().AuthorPrefix) {
                        author = lines[2]
                        subtitle = ""
                    } else {
                        subtitle = lines[2]
                        author = lines[3]
                    }

                    author = strings.TrimPrefix(author, config.GetConfig().AuthorPrefix)
                    info = append(info, Meta{title, subtitle, author, path})
                }
            }
        }
    }

    return info
}

func PrintSortedProperties(property string) {
    resultMap := make(map[string]bool)

    for _, meta := range GetMetaInfo("") {
        switch property {
        case "author":
            resultMap[meta.Author] = true
        case "title":
            resultMap[meta.Title] = true
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

func List(full bool, whichFile string) {
    for _, meta := range GetMetaInfo(whichFile) {
        fmt.Println("Title:", meta.Title)

        if utf8.RuneCountInString(meta.Subtitle) > 0 {
            fmt.Println("Subtitle:", meta.Subtitle)
        }

        fmt.Println("Author:", meta.Author)

        notes := GetNotesFromFile(meta.Path)
        fmt.Println("Note Count:", len(notes))

        wordCount := len(GetAllWordsFromFile(meta.Path))
        fmt.Println("Word Count:", wordCount)
        fmt.Println("Path:", meta.Path)

        if full {
            fmt.Println("Notes:")
            fmt.Println("")
            for _, note := range notes {
                printing.PrintNote(note)
            }
        } else {
            fmt.Println("")
        }
    }
}

func RandomNote(whichFile string) {
    rand.Seed(time.Now().Unix())
    file := whichFile

    if len(whichFile) < 1 {
        files, err := ioutil.ReadDir(config.GetConfig().Directory)
        if err != nil {
            log.Fatal(err)
        }

        randomFileIndex := rand.Intn(len(files) - 1)
        file = config.GetConfig().Directory + files[randomFileIndex].Name()
    }

    notes := GetNotesFromFile(file)
    randomNoteIndex := 0
    if len(notes) > 1 {
        randomNoteIndex = rand.Intn(len(notes) - 1)
    }
    note := notes[randomNoteIndex]
    meta := GetMetaInfo(file)[0]

    fmt.Println("Title:", meta.Title)
    if utf8.RuneCountInString(meta.Subtitle) > 0 {
        fmt.Println("Subtitle:", meta.Subtitle)
    }
    fmt.Println("Author:", meta.Author)
    fmt.Println("")
    printing.PrintNote(note)
}
