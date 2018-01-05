package printing

import (
    "booknotes/config"
    "fmt"
    "log"
    "math"
    "os"
    "os/exec"
    "regexp"
    "strconv"
    "strings"
)

func GetTerminalColumnCount() int {
    columns := -1
    command := exec.Command("stty", "size")
    command.Stdin = os.Stdin
    out, commandErr := command.Output()

    if commandErr != nil {
        log.Fatal(commandErr)
    }

    sttyOutput := string(out)

    if len(sttyOutput) > 0 {
        sizes := strings.Split(sttyOutput, " ")

        if len(sizes) == 2 {
            columnCount, convErr := strconv.Atoi(strings.Replace(sizes[1], "\n", "", 1))

            if convErr != nil {
                log.Fatal(convErr)
            } else {
                columns = columnCount
            }
        }
    }

    return columns
}

func runeToAscii(r rune) string {
    if r < 128 {
        return string(r)
    } else {
        return "\\u" + strconv.FormatInt(int64(r), 16)
    }
}

func isBlank(s string) bool {
    isBlank, matchErr := regexp.MatchString("\\s", s)

    if matchErr != nil {
        log.Fatal(matchErr)
    }

    return isBlank
}

func PrintNote(note string) {
    var cols int
    terminalCols := GetTerminalColumnCount() - 1
    maxLineLength := config.GetConfig().MaxLineLength
    if terminalCols < 1 || maxLineLength < terminalCols {
        cols = maxLineLength
    } else {
        cols = terminalCols
    }

    if cols > 0 {
        noteLength := len(note)
        lastNoteIndex := noteLength - 1
        lastNoteIndex64 := float64(lastNoteIndex)
        var line string
        var start int = 0
        var end int

        for {
            end = int(math.Min(float64(start + cols), lastNoteIndex64))

            if end == lastNoteIndex {
                line = note[start:end + 1]
                start = end + 1
            } else if isBlank(string(note[end])) {
                line = note[start:end]
                start = end + 1
            } else if lastNoteIndex >= end + 1 && isBlank(string(note[end + 1])) {
                line = note[start:end + 1]
                start = end + 2
            } else {
                // search for previous blank
                for i := end; i >= start; i-- {
                    if isBlank(string(note[i])) {
                        line = note[start:i]
                        start = i + 1
                        break
                    }
                }
            }

            fmt.Println(line)

            if end > lastNoteIndex || start > lastNoteIndex {
                break
            }
        }
    } else {
        fmt.Println(note)
    }

    fmt.Println("")
}
