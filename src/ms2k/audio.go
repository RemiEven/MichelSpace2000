package ms2k

import (
	"bytes"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
)

const (
	bytesPerSample = 4

	sampleRate = 44100
)

// NewAudioContext creates a new audio context which can then host sound players
func NewAudioContext() *audio.Context {
	return audio.NewContext(sampleRate)
}

// NewPlayer creates a new player for the given sound in the given audio context.
// The caller is responsible of closing the player.
func NewPlayer(audioContext *audio.Context, sound []byte) (*audio.Player, error) {
	audioStream, err := mp3.Decode(audioContext, bytes.NewReader(sound))
	if err != nil {
		return nil, err
	}
	player, err := audio.NewPlayer(audioContext, audioStream)
	if err != nil {
		return nil, err
	}
	// player.Play()
	return player, nil
}
