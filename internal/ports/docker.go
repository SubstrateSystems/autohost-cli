package ports

type Docker interface {
	Install() error
	Login() error
}
