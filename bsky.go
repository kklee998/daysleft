package main

import (
	"context"
	"log"
	"net/http"

	"github.com/bluesky-social/indigo/api/atproto"
	"github.com/bluesky-social/indigo/api/bsky"
	"github.com/bluesky-social/indigo/atproto/syntax"
	"github.com/bluesky-social/indigo/lex/util"
	"github.com/bluesky-social/indigo/xrpc"
)

func PostToBluesky(text string, username, password string) {
	ctx := context.Background()
	xrpcClient := xrpc.Client{
		Client: &http.Client{},
		Host:   "https://bsky.social",
	}
	session, err := atproto.ServerCreateSession(ctx, &xrpcClient, &atproto.ServerCreateSession_Input{
		Identifier: username,
		Password:   password,
	})
	if err != nil {
		log.Fatal("Failed to create session: " + err.Error())
	}

	xrpcClient.Auth = &xrpc.AuthInfo{
		AccessJwt:  session.AccessJwt,
		RefreshJwt: session.RefreshJwt,
		Did:        session.Did,
		Handle:     session.Handle,
	}

	post := bsky.FeedPost{
		Text:      text,
		CreatedAt: syntax.DatetimeNow().String(),
	}
	resp, err := atproto.RepoCreateRecord(context.Background(), &xrpcClient, &atproto.RepoCreateRecord_Input{
		Repo:       xrpcClient.Auth.Did,
		Collection: "app.bsky.feed.post",
		Record:     &util.LexiconTypeDecoder{Val: &post},
	})

	if err != nil {
		log.Fatal("Failed to post to Bluesky: " + err.Error())
	}
	log.Printf("Posted to Bluesky: %s\n", resp.Uri)
}
