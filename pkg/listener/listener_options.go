package listener

type Options func(t *TcpListener)

func WithPort(port int) Options {
	return func(l *TcpListener) {
		l.Port = port
	}
}
