package ports

type Cloudflare interface {
	Install() error
	Login() error
}
