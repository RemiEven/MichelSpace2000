package main

import (
	"bytes"
	"io"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
)

const bytesPerSample = 4

// Player represents the current audio state
type Player struct {
	audioContext *audio.Context
	audioPlayer  *audio.Player
}

// NewPlayer creates a new player for the given sound in the given audio context
func NewPlayer(audioContext *audio.Context, sound []byte) (*Player, error) {
	type audioStream interface {
		io.ReadSeeker
		Length() int64
	}
	var s audioStream
	var err error
	s, err = mp3.Decode(audioContext, bytes.NewReader(sound))
	if err != nil {
		return nil, err
	}
	p, err := audio.NewPlayer(audioContext, s)
	if err != nil {
		return nil, err
	}
	player := &Player{
		audioContext: audioContext,
		audioPlayer:  p,
	}
	player.audioPlayer.Play()
	return player, nil
}

// Close closes the player
func (p *Player) Close() error {
	return p.audioPlayer.Close()
}
