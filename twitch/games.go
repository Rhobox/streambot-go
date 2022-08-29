package twitch

import (
	"errors"

	"fmt"
	"github.com/nicklaw5/helix"
)

var ErrQueryAmbiguous = errors.New("game query ambiguous")

func queryGame(params *helix.GamesParams) (helix.Game, error) {
	resp, err := Helix.GetGames(params)
	if err != nil {
		return helix.Game{}, err
	}

	if len(resp.Data.Games) != 1 {
		fmt.Println(resp)
		return helix.Game{}, ErrQueryAmbiguous
	}

	return resp.Data.Games[0], nil
}

func GameID(name string) (string, error) {
	game, err := queryGame(&helix.GamesParams{
		Names: []string{name},
	})
	if err != nil {
		return "", err
	}

	return game.ID, nil
}

func GameName(id string) (string, error) {
	game, err := queryGame(&helix.GamesParams{
		IDs: []string{id},
	})
	if err != nil {
		return "", err
	}

	return game.Name, nil
}
