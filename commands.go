package main

import (
	"strings"
	"fmt"
	"strconv"
)

func addCommand(c map[string]func(*State, []string, map[string]*State)bool, command string, action func(*State, []string, map[string]*State) bool) {
	var prog string
	for i := range command {
		prog += string(command[i])
		c[prog] = action
	}
}

//Directions
func doNorth(st *State, s []string, players map[string]*State) bool {
	if len(s) == 0 {
		if st.Room.Exits[0] != *new(Exit) {
			oldRoom := st.Room.ID
			newRoom := st.Room.Exits[0].To
			notifyPlayers(newRoom.ID, players, fmt.Sprintf("%s has entered the room.", st.Name))
			st.Room = newRoom
			notifyPlayers(oldRoom, players, fmt.Sprintf("%s has left the room.", st.Name))
			st.printRoom(players)
		} else {
			st.Output <- "You can't go that direction."
		}
		return true
	}
	return false
}

func doEast(st *State, s []string, players map[string]*State) bool {
	if len(s) == 0 {
		if st.Room.Exits[1] != *new(Exit) {
			oldRoom := st.Room.ID
			newRoom := st.Room.Exits[1].To
			notifyPlayers(newRoom.ID, players, fmt.Sprintf("%s has entered the room.", st.Name))
			st.Room = newRoom
			notifyPlayers(oldRoom, players, fmt.Sprintf("%s has left the room.", st.Name))
			st.printRoom(players)
		} else {
			st.Output <- "You can't go that direction."
		}
		return true
	}
	return false
}

func doSouth(st *State, s []string, players map[string]*State) bool {
	if len(s) == 0 {
		if st.Room.Exits[3] != *new(Exit) {
			oldRoom := st.Room.ID
			newRoom := st.Room.Exits[3].To
			notifyPlayers(newRoom.ID, players, fmt.Sprintf("%s has entered the room.", st.Name))
			st.Room = newRoom
			notifyPlayers(oldRoom, players, fmt.Sprintf("%s has left the room.", st.Name))
			st.printRoom(players)
		} else {
			st.Output <- "You can't go that direction."
		}
		return true
	}
	return false
}

func doWest(st *State, s []string, players map[string]*State) bool {
	if len(s) == 0 {
		if st.Room.Exits[2] != *new(Exit) {
			oldRoom := st.Room.ID
			newRoom := st.Room.Exits[2].To
			notifyPlayers(newRoom.ID, players, fmt.Sprintf("%s has entered the room.", st.Name))
			st.Room = newRoom
			notifyPlayers(oldRoom, players, fmt.Sprintf("%s has left the room.", st.Name))
			st.printRoom(players)
		} else {
			st.Output <- "You can't go that direction."
		}
		return true
	}
	return false
}

func doUp(st *State, s []string, players map[string]*State) bool {
	if len(s) == 0 {
		if st.Room.Exits[4] != *new(Exit) {
			oldRoom := st.Room.ID
			newRoom := st.Room.Exits[4].To
			notifyPlayers(newRoom.ID, players, fmt.Sprintf("%s has entered the room.", st.Name))
			st.Room = newRoom
			notifyPlayers(oldRoom, players, fmt.Sprintf("%s has left the room.", st.Name))
			st.printRoom(players)
		} else {
			st.Output <- "You can't go that direction."
		}
		return true
	}
	return false
}

func doDown(st *State, s []string, players map[string]*State) bool {
	if len(s) == 0 {
		if st.Room.Exits[5] != *new(Exit) {
			oldRoom := st.Room.ID
			newRoom := st.Room.Exits[5].To
			notifyPlayers(newRoom.ID, players, fmt.Sprintf("%s has entered the room.", st.Name))
			st.Room = newRoom
			notifyPlayers(oldRoom, players, fmt.Sprintf("%s has left the room.", st.Name))
			st.printRoom(players)
		} else {
			st.Output <- "You can't go that direction."
		}
		return true
	}
	return false
}
//Interact
func doLook(st *State, s []string, players map[string]*State) bool {
	directions := map[string]int{"north": 0,"east": 1,"west": 2,"south": 3,"up": 4, "down": 5}
	if len(s) == 0 {
		st.printRoom(players)
		return true
	} else if i, ok := directions[strings.ToLower(s[0])]; ok {
		if exit := st.Room.Exits[i]; exit != *new(Exit) {
			st.Output <- exit.Description
			return true
		}
	}
	return false
}

func doGrab(st *State, s []string, players map[string]*State) bool {
	return false
}

func doUse(st *State, s []string, players map[string]*State) bool {
	return false
}

func doApproach(st *State, s []string, players map[string]*State) bool {
	return false
}

func doRetreat(st *State, s []string, players map[string]*State) bool {
	return false
}
//Speak
func doSay(st *State, s []string, players map[string]*State) bool {
	if len(s) != 0 {
		notifyPlayers(st.Room.ID, players, fmt.Sprintf("%s: '%s'", strings.ToUpper(st.Name), strings.Join(s, " ")))
		return true
	}
	return false
}

func doTell(st *State, s []string, players map[string]*State) bool {
	if player, ok := players[s[0]]; ok {
		player.Output <- fmt.Sprintf("%s tells you: %s", st.Name, strings.Join(s[1:], " "))
		st.Output <- fmt.Sprintf("You tell %s: %s", player.Name, strings.Join(s[1:], " "))
		return true
	}
	return false
}
//System
func doCommand(st *State, s []string, players map[string]*State) bool {
	if len(s) == 0 {
		st.Output <- fmt.Sprintf("- commands\n- tell\n- say\n- look\n- north\n- east\n- south\n- west\n- up\n- down\n- quit")
		return true
	}
	return false
}

func doJoin(st *State, s []string, players map[string]*State) bool {
	var n int
	name := "newbie"
	if st.Name != "" {
		name = st.Name
	}
	for _, in := players[st.Name]; in; _, in = players[st.Name] {
		n++
		st.Name = name + strconv.Itoa(n)
	}
	st.Output <- fmt.Sprintf("Welcome %s!\nJoining server...\nTo check commands, type 'commands'\n***************************************", st.Name)
	notifyPlayers(0, players, fmt.Sprintf("%s has joined the server.", st.Name))
	players[st.Name] = st
	st.printRoom(players)
	return true
}

func doQuit(st *State, s []string, players map[string]*State) bool {
	if len(s) == 0 {
		close(st.Output)
		st.Output = nil
		delete(players, st.Name)
		st.Connection.Close()
		notifyPlayers(0, players, fmt.Sprintf("%s has left the server.", st.Name))
		return true
	}
	return false
}
//TODO: Combat system if i have time

//INIT
func initCommands(c map[string]func(*State, []string, map[string]*State)bool) {
	c["join"] = doJoin
	addCommand(c, "quit", doQuit)
	addCommand(c, "commands", doCommand)
	addCommand(c, "tell", doTell)
	addCommand(c, "approach", doApproach)
	addCommand(c, "retreat", doRetreat)
	addCommand(c, "say", doSay)
	addCommand(c, "grab", doGrab)
	addCommand(c, "use", doUse)
	addCommand(c, "look", doLook)
	addCommand(c, "down", doDown)
	addCommand(c, "up", doUp)
	addCommand(c, "west", doWest)
	addCommand(c, "south", doSouth)
	addCommand(c, "east", doEast)
	addCommand(c, "north", doNorth)
}
