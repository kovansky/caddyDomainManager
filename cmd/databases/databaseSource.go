package databases

type DatabaseSource interface {
	Connect() bool
	CreateUser(name string, userHost string, password string) bool
	CreateDatabase(name string) bool
	Close()
}
