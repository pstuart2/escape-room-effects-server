package sound

import (
	"fmt"
	"os"
	"github.com/faiface/beep/wav"
	"github.com/faiface/beep/speaker"
	"time"
	"github.com/faiface/beep"
	"math/rand"
)

const (
	BunnyGrowl        = "./sounds/bunny-growl.wav"
	ChainDoorShut     = "./sounds/chain-door-shut.wav"
	ChainDrag         = "./sounds/chain-drag.wav"
	CreeperExplosion  = "./sounds/creeper-explosion.wav"
	DoorSlam          = "./sounds/door-slam.wav"
	EnderDeath        = "./sounds/ender-death-1.wav"
	Explosion         = "./sounds/explosion.wav"
	MusicLoop         = "./sounds/short-12second-music-loop.wav"
	RandomDoorSlam    = "./sounds/random-door-slam.wav"
	UndergroundEffect = "./sounds/underground-sound-effect.wav"

	Pause       = "./sounds/pause.wav"
	Unpause     = "./sounds/unpause.wav"
	LightToggle = "./sounds/light-toggle.wav"

	CorrectAnswer = "./sounds/correct-answer.wav"
	WrongAnswer   = "./sounds/wrong-answer.wav"
	Clapping   = "./sounds/clapping.wav"
)

func Play(sound string) {
	fmt.Printf("Starting sound: %s\n", sound)

	f, _ := os.Open(sound)
	s, format, _ := wav.Decode(f)
	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))

	done := make(chan struct{})

	speaker.Play(beep.Seq(s, beep.Callback(func() {
		fmt.Printf("Ending sound: %s\n", sound)
		close(done)
	})))

	<-done
}

var effects = []string{BunnyGrowl, ChainDrag, CreeperExplosion, EnderDeath, Explosion, RandomDoorSlam, MusicLoop, UndergroundEffect}

func StartRandomEffects() chan bool {
	stop := make(chan bool)

	go func() {
		r := rand.New(rand.NewSource(time.Now().UnixNano()))

		shouldStop := false
		for !shouldStop {

			nextEffect := effects[r.Intn(len(effects))]
			nextEffectTime := r.Intn(120) + 60
			fmt.Printf("Next effect [%s] playing in [%d]\n", nextEffect, nextEffectTime)

			select {
			case shouldStop = <- stop:
				fmt.Println("Stopping effects")
			case <-time.After(time.Second * time.Duration(nextEffectTime)):
				Play(nextEffect)
			}
		}
	}()

	return stop
}
