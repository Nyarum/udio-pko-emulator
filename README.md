# UDIO: Pirate King Online (China) | Tales of Pirates (English) | Piratia (Russia) server-side emulator
Latest and really ready for testing / developing version of emulator MMORPG Tales of Pirates (or Piratia)

## Story of the emulator way:
- https://github.com/Nyarum/noterius - first public release and it was on Golang, implemented only auth / characters and pretty bad design for sure
- https://github.com/Nyarum/barrel - the processor you see in the emulator now it was barrel parent
- https://github.com/Nyarum/brrraaa - my try to rework architecture of noterius and pretty fast done my motivation
- https://github.com/Nyarum/avocado - pretty fully view of the emulator now but in Crystal language

:test_tube: Source codes I use:
- Thank you PkoDev.net, I took 1.3* sources from server and client (which wow bad I seen)

:clapper: Almost all the work done you can watch:
- https://www.youtube.com/watch?v=6oHN87gFmHk&list=PL91_HHBS7J9a8787S5eJTztJd-5nFzwqq (Playlist on Russian language)

So now I implemented new way for architecture and features\
What's I great for this? Just write you shit code and that's it, small thinking - more done progress lol :3\
Your intuition will do for youself everything, just believe in you

:face_exhaling: Features of architecture:
- Storages plus interface (MongoDB / JSON file-based)
- Docker | Docker-compose files to easy deploy it
- Concurrent way to handle everything - world, players, connections, exchanging between players (haha :3)
- Processor in bi-direction way to have easy way to unmarshal server-side packets in structs and check if the structs written well
- Resource parser (CSV-style format, .txt files from original files)

:face_exhaling: Features of packets:
- Auth in account | Auto-register account
- Create character
- Exit account
- Enter to the game
- Walking in game
- View of another characters
- Chat options and have view chat from another player
- Unview object if a player get out of distance

:face_exhaling: Features of mathemathics:
- Distance between objects
- Angle of character by x, y coords

## How to install

- Install docker / golang (>1.20) on your machine
- Get docker-compose version >3.1
- Run docker-compose up -d
- Run command - "go run main.go"
- Download any client 1.3* and with IPChanger change IP of client to 127.0.0.1
- Now you can login into your account and default password of any account with auto-register option is "testtest"
- Enjoy!
