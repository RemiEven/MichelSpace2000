package audio

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

var (
	audioContext *audio.Context

	assetLibrary *assets.Library
)

func Init(assetLibraryToSet *assets.Library) {
	audioContext = audio.NewContext(sampleRate)
	assetLibrary = assetLibraryToSet
}

// NewMP3Player creates a new player for the given sound in the given audio context.
// The caller is responsible of closing the player.
func NewMP3Player(sound []byte) (*audio.Player, error) {
	audioStream, err := mp3.DecodeWithSampleRate(sampleRate, bytes.NewReader(sound))
	if err != nil {
		return nil, err
	}
	player, err := audioContext.NewPlayer(audioStream)
	if err != nil {
		return nil, err
	}
	player.SetVolume(0.35) // TODO: find a better way to adjust music volume
	player.Play()
	return player, nil
}

// NewWavPlayer creates a new player for the given sound in the fiven audio context.
// The caller is responsible of closing the player.
func NewWavPlayer(sound []byte) (*audio.Player, error) {
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

func PlaySound(soundName string) {
	if err := playSound(soundName); err != nil {
		fmt.Println("failed to play sound [" + soundName + "]: " + err.Error())
	}
}

func playSound(soundName string) error {
	if mp3Sound, ok := assetLibrary.MP3Sounds[soundName]; ok {
		player, err := NewMP3Player(mp3Sound)
		if err != nil {
			return fmt.Errorf("failed to create MP3 player for sound [%v]: %w", soundName, err)
		}
		player.Play()
	} else if wavSound, ok := assetLibrary.WavSounds[soundName]; ok {
		player, err := NewWavPlayer(wavSound)
		if err != nil {
			return fmt.Errorf("failed to create wav player for sound [%v]: %w", soundName, err)
		}
		player.Play()
	}
	return nil
}
