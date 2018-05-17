package itunes

import (
    "encoding/xml"
    "fmt"
    "time"
)

// Track represents an entire iTunes Track data structure
type Track struct {
    TrackID             int
    Size                int
    TotalTime           int
    DiscNumber          int
    DiscCount           int
    TrackNumber         int
    TrackCount          int
    Year                int
    ReleaseDate         time.Time
    DateModified        time.Time
    DateAdded           time.Time
    BitRate             int
    SampleRate          int
    PlayCount           int
    PlayDate            int
    PlayDateUTC         time.Time
    SkipCount           int
    SkipDate            time.Time
    Rating              int
    RatingComputed      bool
    AlbumRating         int
    AlbumRatingComputed bool
    ArtworkCount        int
    PersistentID        string
    TrackType           string
    FileFolderCount     int
    LibraryFolderCount  int
    Compilation         bool
    Name                string
    Artist              string
    AlbumArtist         string
    Composer            string
    Album               string
    Genre               string
    Kind                string
    Location            string
    Podcast             bool
    Unplayed            bool
    SortArtist          string
    SortName            string
    SortAlbum           string
    SortAlbumArtist     string
    SortComposer        string
    Purchased           bool
    Explicit            bool
    Comments            string
    Loved               bool
    Grouping            string
    BPM                 int
    Protected           bool
    HasVideo            bool
    VolumeAdjustment    int
    StartTime           int
    StopTime            int
    Movie               bool
    ITunesU             bool
    Disabled            bool
    PartOfGaplessAlbum  bool
    Series              string
    Episode             string
    EpisodeOrder        int
    TVShow              bool
    Season              string
    ContentRating       string
    Equalizer           string
}

func (t *Track) String() string {

    return fmt.Sprintf("Track: %s by %s (%d)", t.Name, t.Artist, t.TrackID)

}

// UnmarshalXML is a custom unmarshaller function for itunes Track xml format
func (t *Track) UnmarshalXML(decoder *xml.Decoder, start xml.StartElement) error {

    for {
        token, err := decoder.Token()
        if err != nil {
            return err
        }

        switch tok := token.(type) {

        case xml.StartElement:

            if tok.Name.Local == "key" {

                key, err := readOpenTagAsString(decoder)
                if err != nil {
                    return err
                }

                err = resolveKeyOnStruct(t, key, decoder)
                if err != nil {
                    return err
                }

            } else {

                fmt.Printf("skip %s (Track)\n", tok.Name.Local)
                _ = decoder.Skip()

            }

        case xml.EndElement:

            if start.End() == tok {
                return nil
            }

            return NewInvalidFormatError(
                fmt.Sprintf("Unexpected end element in track: %s", tok.Name.Local))

        }

    }

}
