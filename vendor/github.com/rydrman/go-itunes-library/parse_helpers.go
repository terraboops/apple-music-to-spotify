package itunes

import (
    "encoding/xml"
    "fmt"
    "reflect"
    "strconv"
    "strings"
    "time"
)

// readOpenTagAsString reads the current open tag in decoder and returns the
// string upon success
func readOpenTagAsString(decoder *xml.Decoder) (string, error) {

    value := ""
    for {

        token, err := decoder.Token()
        if nil != err {
            return "", err
        }

        switch t := token.(type) {

        case xml.CharData:
            value = string(t)

        case xml.EndElement:
            return value, nil

        default:
            return "", NewInvalidFormatError(
                fmt.Sprintf("Unexpected token type in readOpenTagAsString: %s", reflect.TypeOf(t)))

        }
    }

}

// resolveKeyOnStruct resolves the given key by reading the expected data from decoder
func resolveKeyOnStruct(item interface{}, key string, decoder *xml.Decoder) error {

    token, err := findNextStartElement(decoder)
    if err != nil {
        return err
    }

    // get information about posiible matching
    // struct field
    fieldName := strings.Replace(strings.Title(key), " ", "", -1)
    field := reflect.ValueOf(item).Elem().FieldByName(fieldName)
    if !field.IsValid() {
        fmt.Printf("skipping missing field: %s\n", fieldName)
        _ = decoder.Skip()
        return nil
    }

    // convert the type based on the field type
    switch field.Interface().(type) {

    case string:
        v, err := readOpenTagAsString(decoder)
        if nil != err {
            return err
        }
        field.SetString(v)

    case int:
        v, err := readOpenTagAsString(decoder)
        if nil != err {
            return err
        }
        i, err := strconv.Atoi(v)
        if nil != err {
            return err
        }
        field.SetInt(int64(i))

    case bool:
        v := token.Name.Local
        err = decoder.Skip()
        if nil != err {
            return err
        }
        field.SetBool(v == "true")

    case time.Time:
        v, err := readOpenTagAsString(decoder)
        if nil != err {
            return err
        }
        t, err := time.Parse(time.RFC3339, v)
        if nil != err {
            return err
        }
        field.Set(reflect.ValueOf(t))

    case []byte:
        v, err := readOpenTagAsString(decoder)
        if nil != err {
            return err
        }
        field.Set(reflect.ValueOf([]byte(v)))

    case []*Track:
        tracks, err := decodeTracks(decoder, token)
        if err != nil {
            return err
        }
        field.Set(reflect.ValueOf(tracks))

    case []*Playlist:

        tracks, err := decodePlaylists(decoder, token)
        if err != nil {
            return err
        }
        field.Set(reflect.ValueOf(tracks))

    default:
        return NewInvalidFormatError(
            fmt.Sprintf("unknown field type for %s\n", key))

    }

    return nil

}

func findNextStartElement(decoder *xml.Decoder) (xml.StartElement, error) {

    for {

        token, err := decoder.Token()
        if err != nil {
            return xml.StartElement{}, err
        }

        switch t := token.(type) {

        case xml.StartElement:

            return t, nil

        }

    }

}

func decodeTracks(decoder *xml.Decoder, start xml.StartElement) ([]*Track, error) {

    tracks := make([]*Track, 0)

    for {
        token, err := decoder.Token()
        if err != nil {
            return nil, err
        }

        switch t := token.(type) {

        case xml.StartElement:

            if t.Name.Local == "key" {
                err := decoder.Skip()
                if nil != err {
                    return nil, err
                }
                continue
            }

            if t.Name.Local == "dict" {

                track := &Track{}
                err := track.UnmarshalXML(decoder, t)
                if err != nil {
                    return nil, err
                }

                tracks = append(tracks, track)

            } else {

                return nil, NewInvalidFormatError(
                    "expected dict element after key in track list")
            }

        case xml.EndElement:
            if t == start.End() {
                return tracks, nil
            }
            return nil, NewInvalidFormatError(
                fmt.Sprintf("unexpected end elemet in track list: %s", t.Name.Local))

        case xml.CharData:
            continue

        default:
            return nil, NewInvalidFormatError(
                fmt.Sprintf("Unexpected token type in track list: %s", reflect.TypeOf(t)))

        }

    }

}

func decodePlaylists(decoder *xml.Decoder, start xml.StartElement) ([]*Playlist, error) {

    playlists := make([]*Playlist, 0)

    for {
        token, err := decoder.Token()
        if err != nil {
            return nil, err
        }

        switch t := token.(type) {

        case xml.StartElement:

            if t.Name.Local == "dict" {

                playlist := &Playlist{}
                err := playlist.UnmarshalXML(decoder, t)
                if err != nil {
                    return nil, err
                }

                playlists = append(playlists, playlist)

            } else {

                return nil, NewInvalidFormatError(
                    "expected dict element after key in playlist list")
            }

        case xml.EndElement:
            if t == start.End() {
                return playlists, nil
            }
            return nil, NewInvalidFormatError(
                fmt.Sprintf("unexpected end elemet in playlist list: %s", t.Name.Local))

        case xml.CharData:
            continue

        default:
            return nil, NewInvalidFormatError(
                fmt.Sprintf("Unexpected token type in playlist list: %s", reflect.TypeOf(t)))

        }

    }

}
