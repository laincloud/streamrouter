package model

type StreamAppList map[string][]StreamProc

type StreamApp struct {
	Name        string
	StreamProcs []StreamProc
}

type StreamProc struct {
	Name      string
	Upstreams []StreamUpstream
	Services  []StreamService
}

type StreamUpstream struct {
	Host       string
	InstanceNo int
}

type StreamService struct {
	UpstreamPort int
	ListenPort   int
	Send         string
	Expect       string
}
