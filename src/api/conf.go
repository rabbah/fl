package api

const (
	GitHubClientID = "Ov23liak5XRTpeHgGDtx"

	baseURL                 = "https://macho-dropkick-stun.flows.pstmn.io/api"
	LoginGuestAPI           = baseURL + "/default/fl-login-guest"
	LoginGitHubAPI          = baseURL + "/default/fl-login-github"
	GenerateCmdAPI          = baseURL + "/default/fl-generate"
	StartSubscriptionAPI    = baseURL + "/default/fl-subscription-start"
	CancelSubscriptionAPI   = baseURL + "/default/fl-subscription-cancel"
	StatusOfSubscriptionAPI = baseURL + "/default/fl-subscription-status"
)
