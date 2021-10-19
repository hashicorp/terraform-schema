package backend

type BackendData interface {
	Copy() BackendData
}

type UnknownBackendData struct{}

func (*UnknownBackendData) Copy() BackendData {
	return &UnknownBackendData{}
}

type Remote struct {
	Hostname string
}

func (r *Remote) Copy() BackendData {
	return &Remote{
		Hostname: r.Hostname,
	}
}
