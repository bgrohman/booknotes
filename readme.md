# booknotes

A command-line utility for parsing book notes.

## Usage

    Usage:
    ------
    booknotes <command> [file]

    Commands:
    ---------
    list      Prints title, subtitle, author, and metadata for each book
    full      Same as "list" but includes the full notes, too
    authors   Prints all authors in alphabetical order
    titles    Prints all titles in alphabetical order
    words     Prints words and word counts
    random    Prints a random note
    help      Prints this help message

    Options:
    --------
    file      Optional file path to process instead of all books

    Environment Variables:
    ----------------------
    BOOKNOTES_DIRECTORY        Path to directory containing book note files
    BOOKNOTES_MAX_LINE_LENGTH  Maximum number of columns when printing notes

## Conventions

Book notes / highlights must be stored in a single text file per book with
a specific format. For example, the file the\_signal\_and\_the\_noise.txt:

    ===========================================
    The Signal and the Noise
    Why So Many Predictions Fail but Some Don't
    by Nate Silver
    ===========================================

    "The major difference between a thing that might go wrong and a thing that
    cannot possibly go wrong is that when a thing that cannot possibly go wrong goes
    wrong it usually turns out to be impossible to get at or repair," wrote Douglas
    Adams in The Hitchhiker's Guide to the Galaxy series.37

    In complex systems, however, mistakes are not measured in degrees but in whole
    orders of magnitude.
    
    If our ideas are worthwhile, we ought to be willing to test them by establishing
    falsifiable hypotheses and subjecting them to a prediction.

The format must follow these conventions:

1. The header section is surrounded by rows of "=" characters. 
2. The first content line of the header is the book's title.
3. The second line is the book's subtitle. This line can be omitted.
4. The last content line of the header is the author and should be of the form
   "by <author>".
5. The header section is followed by a blank line before the notes/highlights
   start.
6. Each note/highlight should be either in a single line or on multiple adjacent
   lines.
7. Each note/highlight should be separated by a blank line.

## Why did you make this?

I have a fair number of book notes and highlights organized as plain text files
and wanted a way to parse them for metadata, pull out random notes, etc. And
I needed a project to try out the Go programming language.
