package caching

var drivers = make(map[string]func(config any) (Driver, error))

func Register(name string, driverCreator func(config any) (Driver, error)) {
	drivers[name] = driverCreator
}

func New(name string, config any) (Driver, error) {
	driverCreator, ok := drivers[name]
	if !ok {
		panic("Driver not found: " + name)
	}

	return driverCreator(config)
}
