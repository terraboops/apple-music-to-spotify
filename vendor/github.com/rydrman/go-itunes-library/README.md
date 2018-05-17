# Go Itunes Library Parser

Due to the unique way that Apple uses the XML format to save iTunes libraries, I have created this small go library to handle parsing the data.

I will be using this library in my own tool for exporting itunes libraries to Spotify, and hopefully you will find uses for this code as well.

## Usage

```go
import (
    "fmt"
    "github.com/rydrman/go-itunes-library"
)

libraryFilePath := "path/to/library.xml"

// parse the library file and check for errors
lib, err := itunes.ParseFile(libraryFilePath)
if err != nil {
    panic(err)
}

// print a nice string for each track in the library
for _, track := range lib.Tracks {
    fmt.Println(track)
}

// print a nice string for each playlist in the library
for _, playlist := range lib.Playlists {
    fmt.Println(playlist)
}

```

## Contributing

Please do! If you could run everything through gofmt before submission that would be great, but as this is a tiny library I am not going to be picky.