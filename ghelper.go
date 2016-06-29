/*
			Copyright (C) 2016 Justyn Temme
This program is free software: you can redistribute it and/or modify
    it under the terms of the GNU General Public License as published by
    the Free Software Foundation, either version 3 of the License, or
    (at your option) any later version.

    This program is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU General Public License for more details.

    You should have received a copy of the GNU General Public License
    along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/justyntemme/Ghelper/email"
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

func readNotes(Note *Notepad) string {
	return Note.Data
}

func addEntry(Note *Notepad, Content string) {
	var sb string
	SplitContent := strings.Split(Content, " ")
	for i, cha := range SplitContent {
		if i > 0 {
			sb += cha
		}
	}
	Note.Data += sb
	return
}

func startWorkDay() {
	c1 := exec.Command("/home/user/.scripts/startDay.run")
	_, err := c1.Output()
	if err != nil {
		return
	}
}

func chupdate(ircobj *irc.Connection, event *irc.Event) {
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

func readConfig(Config *Configuration) {
	file, err := os.Open("/home/sir/.config/ghelper.conf")
	if err != nil {
		fmt.Println("could not open file", err)
	}
	_, err = toml.DecodeReader(file, &Config)
	if err != nil {
		fmt.Println("error:", err)
	}
}

func main() {
	//ircobj is the irc object which we can manipulate as a irc bot
	Config := new(Configuration)
	readConfig(Config)
	ircobj := irc.IRC("ghelper", "ghelper")
	Note := new(Notepad)
	ircobj.Connect(Config.IrcServers[0])
	ircobj.SendRawf("/msg nickserv identify %s", Config.Password)
	ircobj.Join(Config.IrcChannels[0])
	//PrivMSG callback function for when bot recives a private message
	ircobj.AddCallback("PRIVMSG", func(event *irc.Event) {
		ircobj.Privmsg("silentmoose", event.Message())
		fmt.Println(event.Message())

		//Checks if username is equal to owners username, and if message is email. Check 3 is an optional password check
		if strings.Split(event.Message(), " ")[0] == "email" && event.Nick == "silentmoose" && strings.Split(event.Arguments[1], " ")[1] == "password" {
			Email.CheckEmail(ircobj, event)
		}
		if strings.Split(event.Message(), " ")[0] == "addNote" && event.Nick == "silentmoose" {
			addEntry(Note, event.Message())
		}
		if strings.Split(event.Message(), " ")[0] == "readNote" && event.Nick == "silentmoose" {
			ircobj.Privmsg("silentmoose", readNotes(Note))
		}
		if strings.Split(event.Message(), " ")[0] == "chupdate" && event.Nick == "silentmoose" {
			chupdate(ircobj, event)
		}
		if strings.Split(event.Message(), " ")[0] == "statWorkDay" && event.Nick == "silentmoose" {
			startWorkDay()
		}
	})
	ircobj.Loop()
}
