package main

import ("fmt"
	_ "github.com/mattn/go-sqlite3"
	"database/sql"
)

type Zone struct {
	ID int
	Name string
	Rooms map[int]*Room
}

type Room struct {
	ID int
	Zone *Zone
	Name string
	Description string
	Exits [6]Exit
}

type Exit struct {
	To *Room
	Description string
}

func getRoom(zones map[int]*Zone, room int) (*Room, error) {
	for _, zone := range zones {
		if val, ok := zone.Rooms[room]; ok {
			return val, nil
		}
	}
	return new(Room), fmt.Errorf("Room does not exist in map")
}

func (r *Room) getExits() map[string]*Room {
	directions := [6]string{"north","east","west","south","up","down"}
	exits := make(map[string]*Room)
	for i, v := range r.Exits {
		if v != *new(Exit) {
			exits[directions[i]] = v.To
		}
	}
	return exits
}

func (st *State) printRoom(players map[string]*State) {
	exits := st.Room.getExits()
	pre := "Exits:"
	peep := "People:"
	rid := st.Room.ID
	for key, _ := range exits {
		pre += " " + key
	}
	for key, v := range players {
		if v.Room.ID == rid && st.Name != key {
			peep += " " + key
		}
	}
	st.Output <- fmt.Sprintf("%s\n\n%s\n%s\n%s", st.Room.Name, st.Room.Description, peep, pre)
}

//-------------------------------------------------------------------------

func openDB() (*sql.DB, error) {
	path := "world.db"
	options := "?" + "_busy_timeout=10000" + "&" + "_foreign_keys=ON"
	db, err := sql.Open("sqlite3", path+options)
	if err != nil {
		return new(sql.DB), fmt.Errorf("Problems with opening database: %v", err)
	}
	return db, nil
}


func initZones(db *sql.DB, zones map[int]*Zone) error {
	rows, err := db.Query("SELECT * FROM zones")
	if err != nil {
		return fmt.Errorf("Querying zones from db: %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		var name string
		err := rows.Scan(&id, &name)
		if err != nil {
			return fmt.Errorf("Scanning row from zone table: %v", err)
		}
		newZone := Zone{id, name, map[int]*Room{}}
		zones[id] = &newZone
	}
	err = rows.Err()
	if err != nil {
		return fmt.Errorf("Error handling rows with zone db: %v", err)
	}
	return nil
}

func populateZones(db *sql.DB, zones map[int]*Zone) error {
	for key, element := range zones {
		rows, err := db.Query("SELECT * FROM rooms WHERE zone_id = ? ORDER BY id", key)
		if err != nil {
			return fmt.Errorf("Querying rooms from db: %v", err)
		}
		for rows.Next() {
			var (
				id, zone_id int
				name, description string
			)
			err := rows.Scan(&id, &zone_id, &name, &description)
			if err != nil {
				rows.Close()
				return fmt.Errorf("Scanning row from room table: %v", err)
			}
			newRoom := Room{id, element, name, description, [6]Exit{}}
			element.Rooms[id] = &newRoom
		}
		err = rows.Err()
		if err != nil {
			rows.Close()
			return fmt.Errorf("Error handling rows with room db: %v", err)
		}
		rows.Close()
	}
	return nil
}

func initExits(db *sql.DB, zones map[int]*Zone) error {
	for _, element := range zones {
		for _, v := range element.Rooms {
			rows, err := db.Query("SELECT to_room_id, direction, exits.description FROM exits WHERE from_room_id = ?" , v.ID)
			if err != nil {
				return fmt.Errorf("Querying exits from db: %v", err)
			}
			for rows.Next() {
				var (
					to_room_id int
					direction, description string
				)
				err := rows.Scan(&to_room_id, &direction, &description )
				if err != nil {
					rows.Close()
					return fmt.Errorf("Scanning row from exit table: %v", err)
				}
				var d int
				switch direction{
				case "n":
					d = 0
				case "e":
					d = 1
				case "w":
					d = 2
				case "s":
					d = 3
				case "u":
					d = 4
				case "d":
					d = 5
				}
				row, err := db.Query("SELECT zone_id FROM rooms WHERE id = ?" , to_room_id)
				if err != nil {
					row.Close()
					return fmt.Errorf("Querying zone id from rooms db: %v", err)
				}
				var zone_id int
				for row.Next() {
					err := row.Scan(&zone_id)
					if err != nil {
						row.Close()
						return fmt.Errorf("Scanning row from room table: %v", err)
					}
				}
				row.Close()
				//get zone from zones map with id
				//get room address from zones room map with id
				//put room address in exit object
				//put exit object in room object
				neededZone, ok := zones[zone_id]
				if !ok {
					rows.Close()
					return fmt.Errorf("Something is up with the db dawg")
				}
				room := neededZone.Rooms[to_room_id]
				newExit := Exit{ room, description}
				v.Exits[d] = newExit
			}
			err = rows.Err()
			if err != nil {
				rows.Close()
				return fmt.Errorf("Error handling rows from exits db: %v", err)
			}
			rows.Close()
		}
	}
	return nil
}
