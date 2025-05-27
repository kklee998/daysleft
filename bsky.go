package main

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/bluesky-social/indigo/api/atproto"
	"github.com/bluesky-social/indigo/api/bsky"
	"github.com/bluesky-social/indigo/atproto/syntax"
	"github.com/bluesky-social/indigo/lex/util"
	"github.com/bluesky-social/indigo/xrpc"
)

// PostToBluesky posts a text message to Bluesky using the provided username and password.
func PostToBluesky(text string, username, password string) (*atproto.RepoCreateRecord_Output, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*3))
	defer cancel()
	xrpcClient := xrpc.Client{
		Client: &http.Client{},
		Host:   "https://bsky.social",
	}
	session, err := atproto.ServerCreateSession(ctx, &xrpcClient, &atproto.ServerCreateSession_Input{
		Identifier: username,
		Password:   password,
	})
	if err != nil {
		return nil, errors.New("failed to create session: " + err.Error())
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
	resp, err := atproto.RepoCreateRecord(ctx, &xrpcClient, &atproto.RepoCreateRecord_Input{
		Repo:       xrpcClient.Auth.Did,
		Collection: "app.bsky.feed.post",
		Record:     &util.LexiconTypeDecoder{Val: &post},
	})

	if err != nil {
		return nil, errors.New("failed to post to Bluesky: " + err.Error())
	}

	return resp, nil
}
