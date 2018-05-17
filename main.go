package main

import (
    "fmt"
    "os"
    "flag"

    "golang.org/x/oauth2"

    "github.com/rydrman/go-itunes-library"
    "github.com/zmb3/spotify"
    "github.com/pkg/errors"
)

func main() {
    // dem flags
    accessToken := flag.String("token", "", "OAuth2 AccessToken from Spotify developer console.")
    libFile := flag.String("library", "", "Apple Music Library.xml file exported from iTunes.")
    flag.Parse()
    if *accessToken == "" || *libFile == "" {
        fmt.Fprintf(os.Stderr, "Missing required flags. Please seek help.\n")
        os.Exit(1)
    }

    // parse the library file and check for errors
    lib, err := itunes.ParseFile(*libFile)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Could not parse the library %s.\n", libFile)
        os.Exit(1)
    }

    // make me a client for to call the spotify
    token := &oauth2.Token{AccessToken: *accessToken}
    client := spotify.Authenticator{}.NewClient(token)
    client.AutoRetry = true
    // Verify we can use the client before continuing
    user, err := client.CurrentUser()
    if err != nil {
        fmt.Fprintf(os.Stderr, err.Error())
        os.Exit(2)
    }
    fmt.Println("Logged in as...")
    fmt.Println("User ID:", user.ID)
    fmt.Println("Display name:", user.DisplayName)

    // Delete most existing playlists
    // didn't look into why it's only most, not all
    resp, err := client.GetPlaylistsForUser(user.ID)
    for _, playlist := range resp.Playlists {
        err = client.UnfollowPlaylist(spotify.ID(user.ID), playlist.ID)
        if err != nil {
            fmt.Fprintf(os.Stderr, err.Error() + "\n")
            continue
        }
    }

    // Match all tracks from iTunes Library to something in Spotify
    fmt.Println("Searching Spotify for equivalent tracks.")
    tracksNotFound := make([]*itunes.Track, 0)
    tracksFound := make([]spotify.ID, 0)
    trackCache := make(map[string]spotify.ID)
    for _, track := range lib.Tracks {
        query := track.Artist + " " + track.Name
        fmt.Println("Searching for:", query)
        srchRes, err := client.Search(query, spotify.SearchTypeTrack)
        if err != nil {
            fmt.Fprintf(os.Stderr, err.Error() + "\n")
            tracksNotFound = append(tracksNotFound, track)
            continue
        }

        if len(srchRes.Tracks.Tracks) > 0 {
            tracksFound = append(tracksFound, srchRes.Tracks.Tracks[0].ID)
            trackCache[track.PersistentID] = srchRes.Tracks.Tracks[0].ID
        } else {
            tracksNotFound = append(tracksNotFound, track)
            continue
        }
    }
    chunkSize := 50
    for i := 0; i < len(tracksFound); i += chunkSize {
        end := i + chunkSize

        if end > len(tracksFound) {
            end = len(tracksFound)
        }

        fmt.Fprintf(os.Stdout, "Adding %d tracks to Spotify.\n", (end - i))
        err = client.AddTracksToLibrary(tracksFound[i:end]...)
        if err != nil {
            fmt.Fprintf(os.Stderr, errors.Wrap(err, "ERROR: Unable to add tracks to library").Error() + "\n")
        }
    }

    // Recreate all playlists from iTunes Library in Spotify
    fmt.Println("Recreate all playlists.")
    for _, applePlist := range lib.Playlists {
        spotifyPlist, err := client.CreatePlaylistForUser(user.ID, applePlist.Name, false)
        if err != nil {
            fmt.Fprintf(os.Stderr, errors.Wrap(err, "ERROR: Unable to create playlist").Error() + "\n")
        }

        chunkSize := 100
        for i := 0; i < len(applePlist.PlaylistItems); i += chunkSize {
            end := i + chunkSize

            if end > len(applePlist.PlaylistItems) {
                end = len(applePlist.PlaylistItems)
            }

            spotifyTracks := make([]spotify.ID, 0)
            for _, appleTrack := range applePlist.PlaylistItems[i:end] {
                if val, cacheHit := trackCache[appleTrack.PersistentID]; cacheHit {
                    spotifyTracks = append(spotifyTracks, val)
                } else {
                    query := appleTrack.Artist + " " + appleTrack.Name
                    srchRes, err := client.Search(query, spotify.SearchTypeTrack)
                    if err != nil {
                        fmt.Fprintf(os.Stderr, err.Error() + "\n")
                        tracksNotFound = append(tracksNotFound, appleTrack)
                        continue
                    }

                    if len(srchRes.Tracks.Tracks) > 0 {
                        spotifyTracks = append(spotifyTracks, srchRes.Tracks.Tracks[0].ID)
                        trackCache[appleTrack.PersistentID] = srchRes.Tracks.Tracks[0].ID
                    } else {
                        tracksNotFound = append(tracksNotFound, appleTrack)
                        continue
                    }
                }
            }

            fmt.Fprintf(os.Stdout, "Adding %d tracks to %s (%s).\n", (end - i), spotifyPlist.Name, spotifyPlist.ID)
            _, err := client.AddTracksToPlaylist(user.ID, spotifyPlist.ID, spotifyTracks...)
            if err != nil {
                fmt.Fprintf(os.Stderr, errors.Wrap(err, "ERROR: Unable to add tracks to playlist").Error() + "\n")
                continue
            }
        }
    }


    fmt.Println("Unable to find these tracks, try adding manually:")
    for _, track := range tracksNotFound {
     fmt.Println(track)
    }
}