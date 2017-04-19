package backend

import "github.com/laincloud/streamrouter/model"

var testApp = model.StreamApp{
	Name: "testStreamApp",
	StreamProcs: []model.StreamProc{
		{
			Name: "testStreamProc1",
			Upstreams: []model.StreamUpstream{
				{
					InstanceNo: 1,
					Host:       "192.168.0.1",
				},
				{
					InstanceNo: 2,
					Host:       "192.168.0.2",
				},
			},
			Services: []model.StreamService{
				{
					ListenPort:   8080,
					UpstreamPort: 9080,
				},
				{
					ListenPort:   8081,
					UpstreamPort: 9081,
				},
			},
		},
		{
			Name: "testStreamProc2",
			Upstreams: []model.StreamUpstream{
				{
					InstanceNo: 1,
					Host:       "192.168.1.1",
				},
				{
					InstanceNo: 2,
					Host:       "192.168.1.2",
				},
			},
			Services: []model.StreamService{
				{
					ListenPort:   8180,
					UpstreamPort: 9180,
				},
				{
					ListenPort:   8181,
					UpstreamPort: 9181,
				},
			},
		},
	},
}
