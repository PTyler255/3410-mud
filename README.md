# 3410-mud

MUD is a multi-user dungeon that runs a central server for users to drop into and interact with the world and other users. Main program runs the interaction and players can join and interact within the spaces including navigate the world and chat. Can easily be expected to include greater functionality and actions as well different worldmaps.

## mud.go
Contains all of the functions and goroutines and the many parts that interact directly with the user.

## commands.go
The file that contains all of the commands used to interact with the world as well as initializing them to be used by the maint routine.

## zones.go
This file holds all of the functions used to initialize the worldmap, including reading the database, but also involves parsing the various world objects such as Zones, Rooms, and Exits.

## world.sql
The extensive schema of the world database, downloaded directly from the assignment page.

## world.db
The world database itself that the main program reads from.

## MUD Screencast
A recording of me demonstrating and then explaining the structure of the the mud code.
