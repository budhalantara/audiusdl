package models

type Track struct {
	// LatestChainBlock int `json:"latest_chain_block"`
	Data TrackData
}

type TrackData struct {
	Title         string
	TrackSegments []TrackSegment `json:"track_segments"`
	User          TrackUser
}

type TrackSegment struct {
	Duration  float32
	Multihash string
}

type TrackUser struct {
	CreatorNodeEndpoint string `json:"creator_node_endpoint"`
}
