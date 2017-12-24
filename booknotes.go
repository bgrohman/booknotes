package main

import (
    "bufio"
    "fmt"
    "io/ioutil"
    "log"
    "os"
    "regexp"
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
    fmt.Println("list      Prints title, subtitle, author, and metadata for each book")
    fmt.Println("full      Same as \"list\" but includes the full notes, too")
    fmt.Println("authors   Prints all authors in alphabetical order")
    fmt.Println("titles    Prints all titles in alphabetical order")
    fmt.Println("words     Prints words and word counts")
    fmt.Println("help      Prints this help message")
    fmt.Println("")
}

type Meta struct {
    title string
    subtitle string
    author string
    path string
}

func getMetaInfo() []Meta {
    var info []Meta

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
                info = append(info, Meta{title, subtitle, author, path})
            }
        }
    }

    return info
}

func sanitizeWordForCount(word string) string {
    lower := strings.ToLower(word)
    punctuation := "()[]:;,.?!\"'"
    sanitized := strings.TrimRight(lower, punctuation)
    sanitized = strings.TrimLeft(sanitized, punctuation)
    return sanitized
}

type WordCountPair struct {
    word string
    count int
}

type WordCountPairList []WordCountPair
func (p WordCountPairList) Swap(i, j int) {
    p[i], p[j] = p[j], p[i]
}
func (p WordCountPairList) Len() int {
    return len(p)
}
func (p WordCountPairList) Less(i, j int) bool {
    return p[i].count < p[j].count
}

var SKIP_WORDS = [...]string {"this", "that", "what", "have", "with", "your", "from", "they", "which", "when", "their", "there", "than", "it's", "were", "them"}

func isSkipWord(word string) bool {
    result := false

    for _, s := range SKIP_WORDS {
        if s == word {
            result = true
            break
        }
    }

    return result
}

func getNotesFromFile(path string) []string {
    var notes []string

    file, fileErr := os.Open(path)
    defer file.Close()

    if fileErr != nil {
        log.Fatal(fileErr)
    }

    currentBlock := ""
    metaCount := 0
    isMeta := false

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        line := scanner.Text()

        isMetaSeparator, matchErr := regexp.MatchString("=+", line)
        if matchErr != nil {
            log.Fatal(matchErr)
        }

        if isMetaSeparator == true {
            metaCount = metaCount + 1

            if metaCount <= 1 {
                isMeta = true
            } else {
                isMeta = false
            }

            continue
        }

        if isMeta == false {
            isBlank, matchErr := regexp.MatchString("^\\s*$", line)
            if matchErr != nil {
                log.Fatal(matchErr)
            }

            if isBlank == false {
                currentBlock = currentBlock + line
            } else {
                if len(currentBlock) > 0 {
                    currentBlock = regexp.MustCompile("\r?\n").ReplaceAllString(currentBlock, " ")
                    notes = append(notes, currentBlock)
                    currentBlock = ""
                }
            }
        }
    }

    if len(currentBlock) > 0 {
        currentBlock = regexp.MustCompile("\r?\n").ReplaceAllString(currentBlock, " ")
        notes = append(notes, currentBlock)
    }

    if err := scanner.Err(); err != nil {
        log.Fatal(err)
    }

    return notes
}

func getAllWordsFromFile(path string) []string {
    var words []string
    byteContents, fileErr := ioutil.ReadFile(path)

    if fileErr != nil {
        log.Fatal(fileErr)
    } else {
        contents := string(byteContents)
        allWords := regexp.MustCompile("\\s").Split(contents, -1)

        for _, word := range allWords {
            match, matchErr := regexp.MatchString("=+", word)

            if matchErr != nil {
                log.Fatal(matchErr)
            }

            if match == false && len(word) >= 4 && isSkipWord(word) == false {
                sanitized := sanitizeWordForCount(word)
                words = append(words, sanitized)
            }
        }
    }

    return words
}

func wordCount() {
    wordCounts := make(map[string]int)

    files, err := ioutil.ReadDir(DIRECTORY)
    if err != nil {
        log.Fatal(err)
    }

    for _, f := range files {
        if strings.HasSuffix(f.Name(), ".txt") {
            words := getAllWordsFromFile(DIRECTORY + f.Name())

            for _, word := range words {
                wordCounts[word] = wordCounts[word] + 1
            }
        }
    }

    pairs := make(WordCountPairList, len(wordCounts))
    i := 0
    for k, v := range wordCounts {
        pairs[i] = WordCountPair{k, v}
        i = i + 1
    }
    sort.Sort(pairs)

    for _, value := range pairs {
        fmt.Println(value.word, value.count)
    }
}

func list(full bool) {
    for _, meta := range getMetaInfo() {
        fmt.Println("Title:", meta.title)

        if utf8.RuneCountInString(meta.subtitle) > 0 {
            fmt.Println("Subtitle:", meta.subtitle)
        }

        fmt.Println("Author:", meta.author)

        notes := getNotesFromFile(meta.path)
        fmt.Println("Note Count:", len(notes))

        wordCount := len(getAllWordsFromFile(meta.path))
        fmt.Println("Word Count:", wordCount)
        fmt.Println("Path:", meta.path)

        if full {
            fmt.Println("Notes:")
            for _, note := range notes {
                fmt.Println(note)
                fmt.Println("")
            }
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
        case "words":
            wordCount()
        case "help":
            help()
        }
    } else {
        help()
    }
}
