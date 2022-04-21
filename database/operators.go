package database

type Operator int

const (
	Equals Operator = iota
	NotEquals
	GreaterThan
	GreaterThanOrEquals
	LessThan
	LessThanOrEquals
	Like
	NotLike
	In
	NotIn
	IsNull
	IsNotNull
)
