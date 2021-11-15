package main

func main() {
	server := NewServer("0.0.0.0", 7050) //nc ip/localhost 7050
	server.Start()
}

//  nc 0.0.0.0 7050
