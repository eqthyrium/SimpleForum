package module

type HttpModule interface {
	SignUp(nickname, email, password, oauth string) error
	LogIn(email, password, oauth string) (string, error)
	//Login(email, password string) (err error)
	// Here we write what kind of services can be used in the http handler
}
