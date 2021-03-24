package main

import (
	"fmt"
	"github.com/veandco/go-sdl2/mix"
	"github.com/veandco/go-sdl2/sdl"
	"time"
)

const (
	WAIT_TIME      = (2 * time.Minute) + (30 * time.Second) // wait time to turn off tea pot
	SONG_PLAY_TIME = (2 * time.Minute) + (17 * time.Second) // length of alarm song
	OPEN_SOUND     = "sounds/herewego.ogg"
	READY_MUSIC    = "sounds/artOfWar.ogg"
)

func main() {
	if err := sdl.Init(sdl.INIT_AUDIO); err != nil {
		fmt.Println("Unable to initialize sdl")
		return
	} else {
		fmt.Println("SDL initialized.")
	}

	defer sdl.Quit()

	if err := mix.Init(mix.INIT_OGG); err != nil {
		fmt.Println("Unable to initialize audio mixer")
		return
	} else {
		fmt.Println("Audio mixer initialized.")
	}

	defer mix.Quit()
	err := mix.OpenAudio(44100, mix.DEFAULT_FORMAT, mix.DEFAULT_CHANNELS, mix.DEFAULT_CHUNKSIZE)

	if err != nil {
		fmt.Println("mix.OpenAudio() call failed")
		return
	}

	startSound, err := mix.LoadMUS(OPEN_SOUND)

	if err != nil {
		fmt.Println("LoadMusic:", OPEN_SOUND, "failed", err)
		return
	}

	readySound, err := mix.LoadMUS(READY_MUSIC)

	if err != nil {
		fmt.Println("LoadMusic:", READY_MUSIC, "failed", err)
		return
	}

	defer startSound.Free()

	defer readySound.Free()

	fmt.Println("All audio loaded.")

	time.Sleep(3 * time.Second)

	fmt.Println("Start making your tea now.")
	startSound.Play(1)

	readyChan := make(chan bool)

	go updateUserOnStatus(readyChan, WAIT_TIME)

	select {
	case <-readyChan:
		fmt.Println("You can turn off the tea pot now.")
	}

	readySound.Play(1)

	go updateUserOnStatus(readyChan, SONG_PLAY_TIME)

	select {
	case <-readyChan:
		fmt.Println("Alarm song is done. Your tea should be pretty how now.")
	}
}

func updateUserOnStatus(ready chan bool, duration time.Duration) {
	now := time.Now()
	fut := now.Add(duration)
	timeSinceLast := now
	const UPDATE_TIME = time.Second * 11 // update the user every 11 seconds

	for ; now.Before(fut); now = time.Now() {
		if timeSinceLast.Add(UPDATE_TIME).Before(now) {
			timeSinceLast = now
			duration -= UPDATE_TIME
			fmt.Println(duration, "seconds left.")
		}
	}

	ready <- true
}
