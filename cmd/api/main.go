package main

import (
	"fmt"
	"net/http"

	"github.com/baywiggins/qIt-backend/internal/api/handlers"
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
													
	handlers.HandleRoutes()

	http.ListenAndServe("localhost:3000", nil)
}