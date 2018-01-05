package main

import (
    "fmt"
    "os"
    "booknotes/core"
)

func help() {
    fmt.Println("Usage:")
    fmt.Println("------")
    fmt.Println("booknotes <command> [file]")
    fmt.Println("")
    fmt.Println("Commands:")
    fmt.Println("---------")
    fmt.Println("list      Prints title, subtitle, author, and metadata for each book")
    fmt.Println("full      Same as \"list\" but includes the full notes, too")
    fmt.Println("authors   Prints all authors in alphabetical order")
    fmt.Println("titles    Prints all titles in alphabetical order")
    fmt.Println("words     Prints words and word counts")
    fmt.Println("random    Prints a random note")
    fmt.Println("help      Prints this help message")
    fmt.Println("")
    fmt.Println("Options:")
    fmt.Println("--------")
    fmt.Println("file      Optional file path to process instead of all books")
    fmt.Println("")
    fmt.Println("Environment Variables:")
    fmt.Println("----------------------")
    fmt.Println("BOOKNOTES_DIRECTORY        Path to directory containing book note files")
    fmt.Println("BOOKNOTES_MAX_LINE_LENGTH  Maximum number of columns when printing notes")
    fmt.Println("")
}

func main() {
    args := os.Args[1:]

    if len(args) > 0 {
        file := ""
        if len(args) >= 2 {
            file = args[1]
        }

        switch args[0] {
        case "list":
            core.List(false, file)
        case "full":
            core.List(true, file)
        case "authors":
            core.PrintSortedProperties("author")
        case "titles":
            core.PrintSortedProperties("title")
        case "words":
            core.WordCount(file)
        case "random":
            core.RandomNote(file)
        case "help":
            help()
        }
    } else {
        help()
    }
}
