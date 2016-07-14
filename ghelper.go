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
	"strings"

	"github.com/justyntemme/Ghelper/bot"
	"github.com/justyntemme/Ghelper/email"
	"github.com/thoj/go-ircevent"
)

func main() {
	//ircobj is the irc object which we can manipulate as a irc bot
	Config := new(Bot.Configuration)
	Bot.ReadConfig(Config)
	ircobj := irc.IRC("ghelper", "ghelper")
	Note := new(Bot.Notepad)
	ircobj.Connect(Config.IrcServers[0])
	ircobj.SendRawf("/msg nickserv identify %s", Config.Password)
	ircobj.Join(Config.IrcChannels[0])
	//PrivMSG callback function for when bot recives a private message
	ircobj.AddCallback("PRIVMSG", func(event *irc.Event) {
		ircobj.Privmsg("silentmoose", event.Message())
		fmt.Println(event.Message())
		switch strings.Split(event.Message(), " ")[0] {

		case "email":
			Email.CheckEmail(ircobj, event)
		case "addNote":
			Bot.AddEntry(Note, event.Message())
		case "readNotes":
			Bot.ReadNotes(Note)
			ircobj.Privmsg("silentmoose", Note.Data)
		case "chupdate":
			Bot.Chupdate(ircobj, event)
		case "startWorkDay":
			Bot.StartWorkDay()
		case "weather":
			ircobj.Privmsg(event.Nick, Bot.CheckWeather(67002)) //TODO dynamically call zip code
		case "help":
			ircobj.Privmsg(event.Nick, Bot.Help())

		}
	})
	ircobj.Loop()
}
