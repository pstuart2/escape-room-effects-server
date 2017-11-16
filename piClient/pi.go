package piClient

import (
	"encoding/json"
	"net/http"
	"bytes"
	"fmt"
)

const (
	PiServer = "http://192.168.86.101:8080"
)

/*
	1 = Secret light
	2 = Nothing
	3 = Wall Lights
	4 = Nothing
	5 = Game Room Lights
	6 = Nothing
	7 = Hall Lights
	8 = Nothing
 */

func LightsOn() {
	post([]uint{0, 0, 1, 0, 1, 0, 1, 0})
}

func LightsOff() {
	post([]uint{0, 0, 0, 0, 0, 0, 0, 0})
}

func WallLightsOnly() {
	post([]uint{0, 0, 1, 0, 0, 0, 0, 0})
}

func GameRoomLightsOnly() {
	post([]uint{0, 0, 0, 0, 1, 0, 0, 0})
}

func SecretLight() {
	post([]uint{1, 0, 0, 0, 0, 0, 0, 0})
}

func post(m []uint) {
	go func() {
		jsonString, _ := json.Marshal(m)

		req, err := http.NewRequest("POST", PiServer+"/lights", bytes.NewBuffer(jsonString))
		if err != nil {
			fmt.Print(err)
			return
		}

		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Print(err)
			return
		}
		defer resp.Body.Close()
	}()
}