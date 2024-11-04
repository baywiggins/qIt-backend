package main

import (
	"fmt"
	"github.com/baywiggins/qIt-backend/internal/server"
)

func main() {
	fmt.Println("Starting API...")
	fmt.Println(`
 ______     ______        ______     ______   __    
/\  ___\   /\  __ \      /\  __ \   /\  == \ /\ \   
\ \ \__ \  \ \ \/\ \     \ \  __ \  \ \  _-/ \ \ \  
 \ \_____\  \ \_____\     \ \_\ \_\  \ \_\    \ \_\ 
  \/_____/   \/_____/      \/_/\/_/   \/_/     \/_/ 
                                                    `)
		
	// Call our function to start the server
	server.StartServer()
}