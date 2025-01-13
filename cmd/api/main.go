package main

import (
	"fmt"
	"log"

	"github.com/baywiggins/qIt-backend/internal/config"
	"github.com/baywiggins/qIt-backend/internal/db"
	"github.com/baywiggins/qIt-backend/internal/server"
	"github.com/baywiggins/qIt-backend/internal/server/pubsub"
)

func main() {
	var err error;
	// Initialize database
	database, err := db.Connect(config.DBName)
	if err != nil {
		log.Fatalf("failed to create server: %s", err)
	}
	defer database.Close()
	
	// Apply migrations
	err = db.Migrate(database)
	if err != nil {
		log.Fatalf("failed to execute migrations: %s", err)
	}

	fmt.Println("Starting API...")
	fmt.Println(`
 ______     ______        ______     ______   __    
/\  ___\   /\  __ \      /\  __ \   /\  == \ /\ \   
\ \ \__ \  \ \ \/\ \     \ \  __ \  \ \  _-/ \ \ \  
 \ \_____\  \ \_____\     \ \_\ \_\  \ \_\    \ \_\ 
  \/_____/   \/_____/      \/_/\/_/   \/_/     \/_/ 
                                                    `)
		
	// Call our function to start the server
	go pubsub.RunWebSocket()
	server.StartServer(database)
}