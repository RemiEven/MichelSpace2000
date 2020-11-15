package main

import (
	"bytes"
	"io"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
)

// Player represents the current audio state.
type Player struct {
	audioContext *audio.Context
	audioPlayer  *audio.Player
}

func NewPlayer(audioContext *audio.Context, sound []byte) (*Player, error) {
	type audioStream interface {
		io.ReadSeeker
		Length() int64
	}
	const bytesPerSample = 4 // TODO: This should be defined in audio package
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
func (p *Player) Close() error {
	return p.audioPlayer.Close()
}
