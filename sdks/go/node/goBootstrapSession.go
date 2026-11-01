package node

import "context"

type goBootstrapSessionKey struct{}
type goBootstrapSession struct {
	path    string
	program *goProgram
	release func()
}

func withGoBootstrapSession(ctx context.Context, path string) (context.Context, *goBootstrapSession) {
	session := &goBootstrapSession{path: path}
	return context.WithValue(ctx, goBootstrapSessionKey{}, session), session
}
func (session *goBootstrapSession) close() {
	if session.release != nil {
		session.release()
		session.release = nil
	}
}
