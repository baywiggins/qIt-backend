package models

// Structs for demarshaling spotify responses
type PlaybackState struct {
	CurrentlyPlayingType string `json:"currently_playing_type"`
	IsPlaying bool `json:"is_playing"`
	Item Item `json:"item"`
	Progress uint32 `json:"progress_ms"`
}

type CurrentlyPlaying struct {
	Item Item `json:"item"`
	Progress uint32 `json:"progress_ms"`
}

type CurrentQueue struct {
	CurrentlyPlaying Item `json:"currently_playing"`
	Queue []Item `json:"queue"`
}

type SearchByTrack struct {
	Tracks Tracks
}

type Tracks struct {
	Items []Item `json:"items"`
	Limit uint8 `json:"limit"`
	Next string `json:"next"`
	Previous string `json:"previous"`
	Offset uint8 `json:"offset"`
	Total int `json:"total"`
}

type Item struct {
	Album Album `json:"album"`
	Artists []Artist `json:"artists"`
	Duration uint32 `json:"duration_ms"`
	ID string `json:"id"`
	Name string `json:"name"`
	URI string `json:"uri"`
}

type Album struct {
	ExternalURLs ExternalURLs `json:"external_urls"`
	ID string `json:"id"`
	Name string `json:"name"`
	Images []Images `json:"images"`
}

type Artist struct {
	ExternalURLs ExternalURLs `json:"external_urls"`
	ID string `json:"id"`
	Name string `json:"name"`
}

type ExternalURLs struct {
	URL string `json:"spotify"`
}

type Images struct {
	Height uint16 `json:"height"`
	Width uint16 `json:"width"`
	URL string `json:"url"`
}