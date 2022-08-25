package twitch

import (
	"github.com/nicklaw5/helix"
)

func Streams(gameID string) (streams []helix.Stream, err error) {
	resp, err := Helix.GetStreams(&helix.StreamsParams{
		GameIDs: []string{gameID},
	})
	if err != nil {
		return nil, err
	}

	streams = append(streams, resp.Data.Streams...)

	for resp.Data.Pagination.Cursor != "" {
		resp, err = Helix.GetStreams(&helix.StreamsParams{
			GameIDs: []string{gameID},
			After:   resp.Data.Pagination.Cursor,
		})

		streams = append(streams, resp.Data.Streams...)
	}

	return
}
