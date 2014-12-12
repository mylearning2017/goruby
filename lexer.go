package main

type itemType int

const (
	itemError itemType = iota


)

type item struct {
	typ itemType
	val string
}


