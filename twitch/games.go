package twitch

import (
	"errors"

	"fmt"
	"github.com/nicklaw5/helix"
)

var ErrQueryAmbiguous = errors.New("game query ambiguous")

func GameID(name string) (string, error) {
	resp, err := Helix.GetGames(&helix.GamesParams{
		Names: []string{name},
	})
	if err != nil {
		return "", err
	}

	if len(resp.Data.Games) != 1 {
		fmt.Println(resp)
		return "", ErrQueryAmbiguous
	}

	return resp.Data.Games[0].ID, nil
}
