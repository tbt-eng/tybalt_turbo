<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/time_claim.go","sha256":"85a06db218f9d3e488140a742072e5431a0fb03e3e58d916b9d31852da40d09d","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/time_claim.go

## Time claim required for time-related actions

Lines 11-20

Turbo refuses a user who does not hold the time claim ("time"). The message is "time claim required". Turbo decides whether the user holds the claim with utilities.HasClaim, which is in another file. If that lookup fails, Turbo returns that error and does not show the "time claim required" message.
