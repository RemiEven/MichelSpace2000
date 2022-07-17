package ms2k

import (
	"bytes"
	"fmt"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"

	"github.com/RemiEven/michelSpace2000/src/ms2k/assets"
)

const (
	bytesPerSample = 4

	sampleRate = 44100
)

// NewAudioContext creates a new audio context which can then host sound players
func NewAudioContext() *audio.Context {
	return audio.NewContext(sampleRate)
}

// NewMP3Player creates a new player for the given sound in the given audio context.
// The caller is responsible of closing the player.
func NewMP3Player(audioContext *audio.Context, sound []byte) (*audio.Player, error) {
	audioStream, err := mp3.DecodeWithSampleRate(sampleRate, bytes.NewReader(sound))
	if err != nil {
		return nil, err
	}
	player, err := audioContext.NewPlayer(audioStream)
	if err != nil {
		return nil, err
	}
	player.Play()
	return player, nil
}

// NewWavPlayer creates a new player for the given sound in the fiven audio context.
// The caller is responsible of closing the player.
func NewWavPlayer(audioContext *audio.Context, sound []byte) (*audio.Player, error) {
	audioStream, err := wav.DecodeWithSampleRate(sampleRate, bytes.NewReader(sound))
	if err != nil {
		return nil, err
	}
	player, err := audioContext.NewPlayer(audioStream)
	if err != nil {
		return nil, err
	}
	player.Play()
	return player, nil
}

func PlaySound(audioContext *audio.Context, assetLibrary *assets.Library, soundName string) error {
	if mp3Sound, ok := assetLibrary.MP3Sounds[soundName]; ok {
		player, err := NewMP3Player(audioContext, mp3Sound)
		if err != nil {
			return fmt.Errorf("failed to create MP3 player for sound [%v]: %w", soundName, err)
		}
		player.Play()
	} else if wavSound, ok := assetLibrary.WavSounds[soundName]; ok {
		player, err := NewWavPlayer(audioContext, wavSound)
		if err != nil {
			return fmt.Errorf("failed to create wav player for sound [%v]: %w", soundName, err)
		}
		player.Play()
	}
	return nil
}
