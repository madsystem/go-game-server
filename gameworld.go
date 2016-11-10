package main

type GameWorld struct {
	players  []GameEntity
	monsters []GameEntity
}


func tickGameworld() {
	updateClient()
}

func updateClient() {
	// do some ai shit
}

func execGameEntityAction() {
	// do some action shit
}
