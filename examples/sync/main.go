package main

import (
	"context"
	"crypto/tls"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"

	"connectrpc.com/connect"
	"github.com/utxorpc/go-codegen/utxorpc/v1alpha/sync"
	utxorpc "github.com/utxorpc/go-sdk"
	"golang.org/x/net/http2"
)

func main() {
	ctx := context.Background()
	httpClient := &http.Client{
		CheckRedirect: func(_ *http.Request, _ []*http.Request) error {
			return http.ErrUseLastResponse
		},
		Transport: &http2.Transport{
			AllowHTTP: true,
			DialTLS: func(network, addr string, _ *tls.Config) (net.Conn, error) {
				// If you're also using this client for non-h2c traffic, you may want
				// to delegate to tls.Dial if the network isn't TCP or the addr isn't
				// in an allowlist.
				return net.Dial(network, addr)
			},
		},
	}
	baseUrl := "http://localhost:58502"
	client := utxorpc.NewClient(httpClient, baseUrl)
	// FollowTipRequest with no Intersect
	// req := connect.NewRequest(&sync.FollowTipRequest{})

	// FollowTipRequest with Intersect
	blockHash, err := hex.DecodeString("230eeba5de6b0198f64a3e801f92fa1ebf0f3a42a74dbd1922187249ad3038e7")
	if err != nil {
		log.Fatalf("failed to decode hex string: %v", err)
	}

	blockRef := &sync.BlockRef{
		Hash: blockHash,
	}
	req := connect.NewRequest(&sync.FollowTipRequest{
		Intersect: []*sync.BlockRef{blockRef},
	})

	fmt.Println("connecting to utxorpc host:", baseUrl)
	resp, err := client.ChainSync.FollowTip(ctx, req)
	if err != nil {
		fmt.Println(connect.CodeOf(err))
		if connectErr := new(connect.Error); errors.As(err, &connectErr) {
			fmt.Println(connectErr.Message())
			fmt.Println(connectErr.Details())
		}
		panic(err)
	}
	fmt.Println("connected to utxorpc...")
	fmt.Printf("Response: %+v\n", resp)
}
