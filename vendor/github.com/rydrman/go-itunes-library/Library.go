package itunes

import (
    "encoding/xml"
    "fmt"
    "io"
    "time"
)

// Library represents an entire iTunes library data structure
type Library struct {
    MajorVersion        int
    MinorVersion        int
    ApplicationVersion  string
    Date                time.Time
    Features            int
    ShowContentRatings  bool
    LibraryPersistentID string
    MusicFolder         string
    Tracks              []*Track
    tracksByID          map[int]*Track
    Playlists           []*Playlist
}

// String returns a nice string representation of this Library
func (lib *Library) String() string {
    return fmt.Sprintf("iTunes Library: %d Tracks, %d Playlists (%s)",
        len(lib.Tracks), len(lib.Playlists), lib.LibraryPersistentID)
}

// sortTracksByID sorts the tracks in the Tracks slice by id into
// the tracksByID slice for easy lookup later
func (lib *Library) sortTracksByID() {

    byID := make(map[int]*Track)

    for _, t := range lib.Tracks {

        byID[t.TrackID] = t

    }

    lib.tracksByID = byID

}

// mapAllPlaylistTracks updates the tracks in each playlist to ensure
// that they map to the valid track struct in the Library.Tracks slice
func (lib *Library) mapAllPlaylistTracks() {

    for _, p := range lib.Playlists {

        for i, t := range p.PlaylistItems {

            p.PlaylistItems[i] = lib.tracksByID[t.TrackID]

        }

    }

}

// UnmarshalXML is a custom unmarshaller function for itunes library xml format
func (lib *Library) UnmarshalXML(decoder *xml.Decoder, start xml.StartElement) error {

    for {
        token, err := decoder.Token()
        if err == io.EOF {
            lib.sortTracksByID()
            lib.mapAllPlaylistTracks()
            return nil
        }
        if err != nil {
            return err
        }

        switch t := token.(type) {

        case xml.StartElement:

            if t.Name.Local == "key" {

                key, err := readOpenTagAsString(decoder)
                if err != nil {
                    return err
                }

                err = resolveKeyOnStruct(lib, key, decoder)
                if err != nil {
                    return err
                }

            } else {

                fmt.Printf("skip %s (Library)\n", t.Name.Local)
                err = decoder.Skip()
                if err != nil {
                    return err
                }

            }

        }

    }

}
