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

func (o Operation) String() string {
	switch o {
	case Select:
		return "SELECT"
	case Update:
		return "UPDATE"
	case Delete:
		return "DELETE"
	case Insert:
		return "INSERT"
	case Drop:
		return "DROP"
	case Create:
		return "CREATE"
	default:
		return "UNKNOWN"
	}
}
