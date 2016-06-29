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

type notepad struct { //circular linked list
	data      string
	nextEntry *notepad //next
	headEntry *notepad //top
}

func readLastEntry(notePad *notepad) string {
	return notePad.data
}

func readFromTop(notePad *notepad) string { //cieves HEAD entry
	notes := ""
	for l := notePad.data; &l != nil; l += "\n" {
		notes += (notePad.data)
	}
	return notes
}

func addEntry(notePad *notepad, content string) {
	notePad.data += content
	return
}

func startWorkDay() {
	c1 := exec.Command("/home/user/.scripts/startDay.run")
	_, err := c1.Output()
	if err != nil {
		return
	}
}
func sendFeed(ircobj *irc.Connection, event *irc.Event) {

}

func Readfeed(ircobj *irc.Connection, event *irc.Event) {
	IrcFeeds := strings.Split(event.Message(), " ")
	for _, feed := range IrcFeeds {
		fmt.Println(feed)
		if strings.Contains(feed, "#") == true {
			fmt.Println(feed)
			ircobj.Join(feed)
			ircobj.AddCallback("*", func(event *irc.Event) {
				ircobj.Privmsg("silentmoose", (event.User + ": " + event.Message()))
			})
		}
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

func main() {
	//ircobj is the irc object which we can manipulate as a irc bot
	var Config Configuration
	toml.DecodeFile("/home/sir/.config/ghelper.json", &Config)
	ircobj := irc.IRC("ghelp", "ghelp")

	ircobj.Connect(Config.IrcServers[0])

	ircobj.Join(Config.IrcChannels[0])
	ircobj.SendRawf("/msg nickserv identify %s", Config.Password)
	//PrivMSG callback function for when bot recives a private message
	ircobj.AddCallback("PRIVMSG", func(event *irc.Event) {
		ircobj.Privmsg("silentmoose", event.Message())
		fmt.Println(event.Message())

		//Checks if username is equal to owners username, and if message is email. Check 3 is an optional password check
		if strings.Split(event.Message(), " ")[0] == "email" && event.Nick == "silentmoose" && strings.Split(event.Arguments[1], " ")[1] == "password" {
			email.CheckEmail(ircobj, event)
		}
		if strings.Split(event.Message(), " ")[0] == "feed" && event.Nick == "Silentmoose" {
			Readfeed(ircobj, event)
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
