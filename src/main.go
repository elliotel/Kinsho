package main

func main() {

	complete := make(chan struct{})
	inOut := make(chan string)
	displayGUI(inOut,complete)
}
