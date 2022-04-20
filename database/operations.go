package database

type Operation int

const (
	Select Operation = iota
	Update
	Delete
	Insert
	Drop
	Create
)
