package main

import (
	"errors"

	"github.com/bwmarrin/discordgo"
	"github.com/shkh/lastfm-go/lastfm"
)

var (
	lastFMApi                  *lastfm.Api
	errUserNotRegisteredLastFM = errors.New("User not registered with LastFm")
	errNoTracksFound           = errors.New("No recent tracks found")
)

type lastFMSong struct {
	Name           string
	URL            string
	Artist         string
	Album          string
	ImageThumbnail string
}

func createLastFMAPIInstance() {
	lastFMApi = lastfm.New(loadedConfigData.LastFMKey, loadedConfigData.LastFMSecret)
}

func registerUserLastFM(user *discordgo.User, lastFmUsername string) error {
	currentUserStruct, err := getUserStruct(user)
	if err != nil {
		return err
	}
	currentUserStruct.LastFmAccount = lastFmUsername
	loadedUserData.Users[user.ID] = currentUserStruct
	return nil
}

func getUserLastListened(user userStruct) (lastFMSong, error) {
	tracks, err := lastFMApi.User.GetRecentTracks(lastfm.P{
		"limit": 1,
		"user":  user.LastFmAccount,
	})
	if err != nil {
		return lastFMSong{}, err
	}
	if len(tracks.Tracks) < 1 {
		return lastFMSong{}, errNoTracksFound
	}

	mostRecentSong := lastFMSong{
		Name:           tracks.Tracks[0].Name,
		URL:            tracks.Tracks[0].Url,
		Artist:         tracks.Tracks[0].Artist.Name,
		Album:          tracks.Tracks[0].Album.Name,
		ImageThumbnail: tracks.Tracks[0].Images[3].Url,
	}

	return mostRecentSong, nil
}

func getUserLastLoved(user userStruct) (lastFMSong, error) {
	// terrible code reuse should really start passing functions around for formatting but lastfm api is garbage
	tracks, err := lastFMApi.User.GetLovedTracks(lastfm.P{
		"limit": 1,
		"user":  user.LastFmAccount,
	})
	if err != nil {
		return lastFMSong{}, err
	}
	if len(tracks.Tracks) < 1 {
		return lastFMSong{}, errNoTracksFound
	}

	mostRecentSong := lastFMSong{
		Name:           tracks.Tracks[0].Name,
		URL:            tracks.Tracks[0].Url,
		Artist:         tracks.Tracks[0].Artist.Name,
		ImageThumbnail: tracks.Tracks[0].Images[3].Url,
	}

	return mostRecentSong, nil
}
