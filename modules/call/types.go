package call

type Offer struct {
    SDP string `json:"sdp"`
}

type Answer struct {
    SDP string `json:"sdp"`
}

type ICECandidate struct {
    Candidate string `json:"candidate"`
}
