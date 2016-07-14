package Bot

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/thoj/go-ircevent"
)

type Configuration struct {
	IrcServers  []string
	IrcChannels []string
	Password    string
}

type Notepad struct { //circular linked list
	Data string
}

func ReadNotes(Note *Notepad) string {
	return Note.Data
}

func AddEntry(Note *Notepad, Content string) {
	var sb string
	SplitContent := strings.Split(Content, " ")
	for i, cha := range SplitContent {
		if i > 0 {
			sb += cha
		}
	}
	Note.Data += (sb + "\n")
	return
}

func StartWorkDay() {
	c1 := exec.Command("/home/user/.scripts/startDay.run")
	_, err := c1.Output()
	if err != nil {
		return
	}
}

func Chupdate(ircobj *irc.Connection, event *irc.Event) {
	c1 := exec.Command("/usr/bin/zypper", "lp")
	out, err := c1.Output()
	out2 := strings.Split(string(out), "\n")
	if err != nil {
		return
	}
	for i := 0; i < 10; i++ {
		ircobj.Privmsg(event.Nick, out2[i])
		time.Sleep(500 * time.Millisecond)
	}
}

func ReadConfig(Config *Configuration) {
	file, err := os.Open("/home/sir/.config/ghelper.conf")
	if err != nil {
		fmt.Println("could not open file", err)
	}
	_, err = toml.DecodeReader(file, &Config)
	if err != nil {
		fmt.Println("error:", err)
	}
}
