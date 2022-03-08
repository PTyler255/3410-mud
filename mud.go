package main

import (
	"bufio"
	"fmt"
	"log"
	"strings"
	"net"
)


type State struct {
	Name string
	Room *Room
	//Inventory map[string]string
	Connection net.Conn
	Output chan string
}

type Event struct {
	Player *State
	Command string
}

func main() {
	//INIT WORLD
	commands := make(map[string]func(*State, []string, map[string]*State)bool)
	initCommands(commands)
	db, err := openDB()
	if err != nil {
		log.Fatalf("%v", err)
	}
	zones := make(map[int]*Zone)
	err = initZones(db, zones)
	if err != nil {
		log.Fatalf("%v", err)
	}
	err = populateZones(db, zones)
	if err != nil {
		log.Fatalf("%v", err)
	}
	err = initExits(db, zones)
	if err != nil {
		log.Fatalf("%v", err)
	}
	players := make(map[string]*State)
	input := make(chan Event)
	startRoom, err := getRoom(zones, 3001)
	if err != nil {
		log.Fatalf("%v", err)
	}
	//-----------------------------------------------
	go manageConnections(startRoom, input)
	for action := range input {
		//notifyPlayers(0, players, "TEST")
		command := action.Command
		player := action.Player
		if player.Output != nil {
			if val, ok := commands[strings.ToLower(strings.Fields(command)[0])]; ok {
				if !(val(player, strings.Fields(command)[1:], players)) {
					player.Output <- "Huh?"
				}
			} else {
				player.Output <- "Huh?"
			}
		}
	}
}

func manageConnections(room *Room, input chan Event) {
	ln, err := net.Listen("tcp", ":3410")
	if err != nil {
		log.Fatalf("%v", err)
	}
	log.Printf("Server is running.\n")
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("%v\n", err)
			continue
		}
		var player State
		player.Room = room
		player.Connection = conn
		player.Output = make(chan string)
		log.Printf("New connection.\n")
		go player.runOutput()
		go player.runInput(input)
	}
}


func (st *State) runInput(input chan Event) {
	scan := bufio.NewScanner(st.Connection)
	st.Output <- "Name?"
	scan.Scan()
	name := scan.Text()
	st.Name = name
	input <- Event{
		Player: st,
		Command: "join",
	}
	for scan.Scan() {
		response := scan.Text()
		if len(response) != 0 {
			if strings.Fields(response)[0] != "join" {
				input <- Event{
					Player: st,
					Command: response,
				}
			} else {
				st.Output <- "Huh?"
			}
		}
	}
	input <- Event{
		Player: st,
		Command: "quit",
	}
}

func (st *State) runOutput() {
	for message := range st.Output {
		st.Printf("\n%s\n> ", message)
	}
}

func (st *State) Printf(format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	_, err := fmt.Fprint(st.Connection, msg)
	if err != nil {
		log.Printf("network error while printing: %v", err)
	}
}

func notifyPlayers(roomId int, players map[string]*State, message string) {
	//If only notify users in zones, check to see roomId <= 100 & > 0
	for _, v := range players {
		if roomId == 0 || v.Room.ID == roomId {
			v.Output <- message
		}
	}
}

/*
func runWorld(player State, commands map[string]func(*State, []string)bool, conn net.Conn) {
	player.Connection = conn
	player.Printf("Name? ")
	scan := bufio.NewScanner(conn)
	//Print intro
	scan.Scan()
	name := scan.Text()
	player.Name = name
	player.Printf("Welcome, %s.\n Joining server...\n", player.Name)
	player.Printf("**************************************************************************\n")
	player.printRoom()
	player.Printf("> ")
	for scan.Scan() {
		response := scan.Text()
		if len(response) != 0 {
			if val, ok := commands[strings.ToLower(strings.Fields(response)[0])]; ok {
				if !(val(&player, strings.Fields(response)[1:])) {
					player.Printf("Huh?\n")
				}
			}else{
				player.Printf("Huh?\n")
			}
		}
		player.Printf("> ")
	}
}
*/
