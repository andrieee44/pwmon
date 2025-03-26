package pwmon

import (
	"bytes"
	"encoding/json"
	"io"
	"math"
	"os/exec"
	"strconv"
	"time"
)

type Info struct {
	Volume int
	Mute   bool
}

func defaultAudioSinkID() (uint32, error) {
	var (
		cmd   *exec.Cmd
		buf   bytes.Buffer
		idStr string
		id    uint64
		err   error
	)

	cmd = exec.Command("wpctl", "inspect", "@DEFAULT_AUDIO_SINK@")
	cmd.Stdout = &buf

	if cmd.Err != nil {
		return 0, cmd.Err
	}

	err = cmd.Run()
	if err != nil {
		return 0, err
	}

	idStr, err = buf.ReadString(',')
	if err != nil {
		return 0, err
	}

	id, err = strconv.ParseUint(idStr[3:len(idStr)-1], 10, 32)
	if err != nil {
		return 0, err
	}

	return uint32(id), nil
}

func monitorDump(id uint32, infoChan chan<- *Info, errChan chan<- error) {
	var (
		cmd     *exec.Cmd
		stdout  io.ReadCloser
		decoder *json.Decoder
		idx     int
		err     error

		infoJson []struct {
			Id uint32 `json:"id"`

			Info struct {
				Params struct {
					Props []struct {
						Mute           bool      `json:"mute"`
						ChannelVolumes []float64 `json:"channelVolumes"`
					} `json:"Props"`
				} `json:"params"`
			} `json:"info"`
		}
	)

	cmd = exec.Command("pw-dump", "-mN")
	stdout, err = cmd.StdoutPipe()
	if err != nil {
		errChan <- err

		return
	}

	err = cmd.Start()
	if err != nil {
		errChan <- err

		return
	}

	decoder = json.NewDecoder(stdout)

	for {
		err = decoder.Decode(&infoJson)
		if err != nil {
			break
		}

		for idx = range infoJson {
			if infoJson[idx].Id != id {
				continue
			}

			infoChan <- &Info{
				Volume: int(math.Cbrt(infoJson[idx].Info.Params.Props[0].ChannelVolumes[0]) * 100),
				Mute:   infoJson[idx].Info.Params.Props[0].Mute,
			}
		}
	}

	errChan <- err

	err = cmd.Wait()
	if err != nil {
		errChan <- err

		return
	}
}

func Monitor() (<-chan *Info, <-chan error, error) {
	var (
		infoChan chan *Info
		errChan  chan error
		id       uint32
		err      error
	)

	infoChan = make(chan *Info)
	errChan = make(chan error)

	for range 10 {
		id, err = defaultAudioSinkID()
		if err != nil {
			time.Sleep(time.Second)

			continue
		}

		go monitorDump(id, infoChan, errChan)

		return infoChan, errChan, nil
	}

	return nil, nil, err
}
