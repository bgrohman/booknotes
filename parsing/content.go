package parsing

import (
    "bufio"
    "fmt"
    "io/ioutil"
    "log"
    "os"
    "regexp"
    "sort"
    "strings"
)

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

func GetNotesFromFile(path string) []string {
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
                if len(currentBlock) > 0 {
                    currentBlock = currentBlock + " " + line
                } else {
                    currentBlock = currentBlock + line
                }
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

func GetAllWordsFromFile(path string) []string {
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

func WordCount(whichFile string) {
    wordCounts := make(map[string]int)

    files, err := ioutil.ReadDir(DIRECTORY)
    if err != nil {
        log.Fatal(err)
    }

    for _, f := range files {
        if strings.HasSuffix(f.Name(), ".txt") {
            path := DIRECTORY + f.Name()

            if len(whichFile) < 1 || whichFile == path {
                words := GetAllWordsFromFile(path)

                for _, word := range words {
                    wordCounts[word] = wordCounts[word] + 1
                }
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
