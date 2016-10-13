package handlers

func (h *Handler) GetRoutes() []Route {
	return []Route{
		// redirects to https://github.com/quantum/discovery
		{
			"GetHome",
			[]string{"GET"},
			"/",
			h.HomeHandler,
		},
		// Returns a new token
		{
			"GetNewToken",
			[]string{"GET"},
			"/new",
			h.NewTokenHandler,
		},
		// Returns the health of discovery service
		{
			"GetHealth",
			[]string{"GET"},
			"/health",
			h.HealthHandler,
		},
		// Prevents robots from indexing
		{
			"GetRobot",
			[]string{"GET"},
			"/robots.txt",
			h.RobotsHandler,
		},
		// Returns a JSON of current cluster status
		{
			"Token",
			[]string{"GET", "PUT"},
			"/{token:[a-f0-9]{32}}",
			h.TokenHandler,
		},
		// This is used by etcd instances to form a quorum during the bootstrap process.
		{
			"Machine",
			[]string{"GET", "PUT", "DELETE"},
			"/{token:[a-f0-9]{32}}/{machine}",
			h.TokenHandler,
		},
		// get the current size of the etcd quorum
		{
			"GetClusterSize",
			[]string{"GET"},
			"/{token:[a-f0-9]{32}}/_config/size",
			h.TokenHandler,
		},
		// will be used by the current cluster leader to update the latest membership info
		{
			"PostRenew",
			[]string{"POST"},
			"/renew",
			h.RenewHandler,
		},
		// can be used to query the latest etcd quorum members in a cleaner format compare to the standard etcd
		// format provided by GET "/{token:[a-f0-9]{32}}"
		{
			"GetClusterMembers",
			[]string{"GET"},
			"/members/{token:[a-f0-9]{32}}",
			h.MemberHandler,
		},
	}
}
