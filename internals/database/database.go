package database

import (
	"context"
	"mini-paas/ent"
)

func Connect(databaseURL string) (*ent.Client, error){
	client, err := ent.Open(
		"pgx",
		databaseURL,
	)

	if err != nil {
		return nil, err
	}

	if err := client.Schema.Create(context.Background()); err != nil {
		return nil, err
	}

	return client, nil
}