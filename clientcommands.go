package main

type clientCmd interface {
	getCmdType() string
}

type attackInfo struct {
	attackerID int32
	targetID   int32
}

func (attackInfo *attackInfo) getCmdType() string {
	return "AttackCmd"
}
