package pwmon

import (
	"bufio"
	"bytes"
	"io"
	"os/exec"
	"strconv"
	"strings"
)

type Info struct {
	Volume int
	Mute   bool
}

func getInfo() (*Info, error) {
	var (
		cmd    *exec.Cmd
		buf    bytes.Buffer
		info   *Info
		fields []string
		vol    float64
		err    error
	)

	cmd = exec.Command("wpctl", "get-volume", "@DEFAULT_AUDIO_SINK@")
	cmd.Stdout = &buf

	if cmd.Err != nil {
		return nil, cmd.Err
	}

	err = cmd.Run()
	if err != nil {
		return nil, err
	}

	info = new(Info)
	fields = strings.Fields(buf.String())

	vol, err = strconv.ParseFloat(fields[1], 64)
	if err != nil {
		return nil, err
	}

	info.Volume = int(vol * 100)

	if len(fields) == 3 {
		info.Mute = true
	}

	return info, nil
}

func monitorDump(infoChan chan<- *Info, errChan chan<- error) {
	var (
		cmd           *exec.Cmd
		stdout        io.ReadCloser
		scanner       *bufio.Scanner
		info, oldInfo *Info
		err           error
	)

	cmd = exec.Command("pactl", "subscribe")
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

	scanner = bufio.NewScanner(stdout)

	for scanner.Scan() {
		if !strings.HasPrefix(scanner.Text(), "Event 'change' on sink #") {
			continue
		}

		info, err = getInfo()
		if err != nil {
			errChan <- err

			return
		}

		if oldInfo == nil || oldInfo.Volume != info.Volume || oldInfo.Mute != info.Mute {
			infoChan <- info
		}

		oldInfo = info
	}

	err = scanner.Err()
	if err != nil {
		errChan <- err

		return
	}

	err = cmd.Wait()
	if err != nil {
		errChan <- err
	}
}

func Monitor() (<-chan *Info, <-chan error, error) {
	var (
		infoChan chan *Info
		errChan  chan error
	)

	infoChan = make(chan *Info)
	errChan = make(chan error)
	go monitorDump(infoChan, errChan)

	return infoChan, errChan, nil
}
