package itunes

import (
    "encoding/xml"
    "io"
    "os"
)

// ParseFile parses the file at the given filepath
// as an itunes library file
func ParseFile(filename string) (*Library, error) {

    libraryFile, err := os.Open(filename)
    if err != nil {
        return nil, err
    }
    defer libraryFile.Close()

    return ParseReader(libraryFile)

}

// ParseReader parses the given readers bytes as an library xml file
func ParseReader(input io.Reader) (*Library, error) {

    lib := &Library{}

    decoder := xml.NewDecoder(input)

    for {
        token, err := decoder.Token()
        if nil != err {

            if err == io.EOF {
                return nil, NewInvalidFormatError(
                    "Unexpected end of library file")
            }

            return nil, err
        }

        switch t := token.(type) {

        case xml.StartElement:

            if t.Name.Local == "dict" {
                err = decoder.DecodeElement(lib, &t)
                if nil != err {
                    return nil, err
                }
                return lib, nil
            }

        }
    }

}
