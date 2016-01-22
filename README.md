# Go Booklet
A small go program that converts pdf-files to booklets using LaTeX.

## Installation
```go install github.com/jjdekker/gobooklet```

## Usage
```gobooklet [inputFile] ([outputFile])?```

## Restrictions
- Needs a working latex installation with pdfpages and latexmk to work.
- Filenames can not contain any latex special characters (escaping characters in filenames in LaTeX seems to be impossible)
- Number of pages is found with a trick. This might not work for every document. If you have one that you can send me, I'll gladly have a look.
