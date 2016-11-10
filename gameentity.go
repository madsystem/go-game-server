package main

type GameEntity struct {
	id         uint32 //-1 = invalid
	entityType uint32 // 0 = player, 1 = monster
	alive      bool
	pos        [2]float32
	vel        [2]float32
}
