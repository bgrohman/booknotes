package config

import (
    "log"
    "os"
    "strconv"
)

type BooknoteConfig struct {
    Directory string
    AuthorPrefix string
    MaxLineLength int
}

const DEFAULT_DIRECTORY = "/home/bryan/Documents/book_notes/"
const DEFAULT_AUTHOR_PREFIX = "by "
const DEFAULT_MAX_LINE_LENGTH = 80

var loadedConfig bool = false
var cachedConfig BooknoteConfig

func GetConfig() BooknoteConfig {
    if loadedConfig == true {
        return cachedConfig
    } else {
        config := BooknoteConfig{DEFAULT_DIRECTORY, DEFAULT_AUTHOR_PREFIX, DEFAULT_MAX_LINE_LENGTH}

        var directory string = os.Getenv("BOOKNOTES_DIRECTORY")
        if len(directory) > 0 {
            config.Directory = directory
        }

        maxLineLengthStr := os.Getenv("BOOKNOTES_MAX_LINE_LENGTH")
        if len(maxLineLengthStr) > 0 {
            maxLineLength, parseErr := strconv.Atoi(maxLineLengthStr)
            if parseErr != nil {
                log.Fatal(parseErr)
            }
            if maxLineLength > 0 {
                config.MaxLineLength = maxLineLength
            }
        }

        cachedConfig = config
        loadedConfig = true
        return config
    }
}
