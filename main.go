package main

import (
	"context"
	"flag"
	"fmt"
	"little/handlers"
	"little/packets"
	"little/resourceparser"
	"log"
	"math/rand"
	"net"
	"os"
	"sync"
	"time"

	"little/storage"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Client struct {
	packets *packets.Packets
	handler *handlers.Handler
}

type mainServer struct {
	players *handlers.Players
	world   *handlers.World
	clients map[net.Conn]Client
	storage handlers.Storage
	sync.RWMutex
}

func NewServer(storage handlers.Storage) *mainServer {
	players := handlers.NewPlayers()

	return &mainServer{
		clients: make(map[net.Conn]Client),
		storage: storage,
		RWMutex: sync.RWMutex{},
		players: players,
		world:   handlers.NewWorld(players),
	}
}

func (es *mainServer) OnOpen(conn net.Conn) []byte {
	es.Lock()
	defer es.Unlock()

	fmt.Println("open connect")

	firstDate := &packets.Date{}

	resBuf := packets.BuildPacket(firstDate)

	handler := handlers.NewHandler(firstDate.Time.Value, es.players, es.world, es.storage, conn)

	es.clients[conn] = Client{packets.NewPackets(), handler}

	return resBuf
}

func (es *mainServer) OnPacket(c net.Conn, buf []byte) {
	if len(buf) == 2 {
		c.Write([]byte{0x00, 0x02})
		return
	}

	header := packets.UnpackPacket(&buf)

	fmt.Println("Accept new packet, with opcode:", header.Opcode)

	if header.Opcode == 432 {
		es.clients[c].handler.Destroy()
		delete(es.clients, c)
		c.Close()

		return
	}

	es.Lock()
	client := es.clients[c]
	es.Unlock()

	packet := client.packets.GetPacket(header.Opcode)
	if packet == nil {
		return
	}

	packet.Process(&buf)

	resPackets := client.handler.Ev(context.Background(), client.packets.GetPacket(header.Opcode))

	for _, pkt := range resPackets {
		fmt.Printf("Sending a packet to client with opcode: %v\n", pkt.Opcode())

		resBuf := packets.BuildPacket(pkt)

		c.Write(resBuf)
	}
}

func main() {
	rand.Seed(time.Now().Unix())

	var port int
	var multicore bool

	flag.IntVar(&port, "port", 1973, "--port 1973")
	flag.BoolVar(&multicore, "multicore", false, "--multicore true")
	flag.Parse()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resourceParser, err := resourceparser.NewParser()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Item pre:", len(resourceParser.ItemPre))
	fmt.Println("Item info:", len(resourceParser.ItemInfo))

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017").SetAuth(options.Credential{
		Username: "root",
		Password: "example",
	}))
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}

	storageMongo := storage.NewMongo(client)

	storagePath := "database.json"
	storage, err := storage.NewStorage(storagePath)
	if err != nil {
		log.Fatal(err)
	}
	defer storage.Save()

	addr := fmt.Sprintf("127.0.0.1:%d", port)

	ln, _ := net.Listen("tcp", addr)
	defer ln.Close()

	main := NewServer(storageMongo)

	fmt.Printf("Running server on addr %v \n", addr)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}

		fmt.Println("Accepted a new connect with ip:", conn.RemoteAddr())

		go func(conn net.Conn) {
			_, err := conn.Write(main.OnOpen(conn))
			if err != nil {
				fmt.Println(err)
				return
			}

			for {
				buffer := make([]byte, 20196)
				ln, err := conn.Read(buffer)
				if err != nil {
					fmt.Println(err)
					break
				}

				main.OnPacket(conn, buffer[:ln])
			}
		}(conn)
	}
}
