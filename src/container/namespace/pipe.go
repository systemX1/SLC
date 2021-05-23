package namespace

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func sendCommand(comArray []string, writePipe *os.File) {
	cmd := strings.Join(comArray, " ")
	_, err := writePipe.WriteString(cmd)
	if err != nil {
		fmt.Println(err)
	}

	err = writePipe.Close()
	if err != nil {
		fmt.Println(err)
	}
}

func readCommand() []string {
	pipe := os.NewFile(uintptr(3), "pipe")
	defer func(pipe *os.File) {
		err := pipe.Close()
		if err != nil {

		}
	}(pipe)
	msg, err := ioutil.ReadAll(pipe)
	if err != nil {
		fmt.Printf("init read pipe error %v", err)
		return nil
	}
	msgStr := string(msg)
	return strings.Split(msgStr, " ")
}
