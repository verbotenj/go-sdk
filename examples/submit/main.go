package main

import (
	"context"
	"crypto/tls"
	"encoding/hex"
	"errors"
	"fmt"
	"net"
	"net/http"

	"connectrpc.com/connect"
	"github.com/utxorpc/go-codegen/utxorpc/v1alpha/submit"
	utxorpc "github.com/utxorpc/go-sdk"
	"golang.org/x/net/http2"
)

func main() {
	// Set mode to "submitTx", "readMempool", "waitForTx", or "watchMempool" to select the desired example.
	var mode string = "submitTx"

	ctx := context.Background()
	httpClient := &http.Client{
		CheckRedirect: func(_ *http.Request, _ []*http.Request) error {
			return http.ErrUseLastResponse
		},
		Transport: &http2.Transport{
			AllowHTTP: true,
			DialTLS: func(network, addr string, _ *tls.Config) (net.Conn, error) {
				return net.Dial(network, addr)
			},
		},
	}
	baseUrl := "http://localhost:58502"
	client := utxorpc.NewClient(httpClient, baseUrl)

	switch mode {
	case "submitTx":
		submitTx(ctx, &client)
	case "readMempool":
		readMempool(ctx, &client)
	case "waitForTx":
		waitForTx(ctx, &client)
	case "watchMempool":
		watchMempool(ctx, &client)
	default:
		fmt.Println("Unknown mode:", mode)
	}
}

func submitTx(ctx context.Context, client *utxorpc.UtxorpcClient) {
	txCbor := "Replace this with the signed transaction in CBOR format."
	txRawBytes, err := hex.DecodeString(txCbor)
	if err != nil {
		panic(fmt.Errorf("failed to decode transaction hash: %v", err))
	}

	// Create a SubmitTxRequest with the transaction data
	tx := &submit.AnyChainTx{
		Type: &submit.AnyChainTx_Raw{
			Raw: txRawBytes,
		},
	}

	// Create a list with one transaction
	req := connect.NewRequest(&submit.SubmitTxRequest{
		Tx: []*submit.AnyChainTx{tx},
	})

	// req := connect.NewRequest(&submit.SubmitTxRequest{})
	fmt.Println("Connecting to utxorpc host:", client.URL())
	resp, err := client.Submit.SubmitTx(ctx, req)
	if err != nil {
		fmt.Println(connect.CodeOf(err))
		if connectErr := new(connect.Error); errors.As(err, &connectErr) {
			fmt.Println(connectErr.Message())
			fmt.Println(connectErr.Details())
		}
		panic(err)
	}
	fmt.Printf("Response: %+v\n", resp)
}

func readMempool(ctx context.Context, client *utxorpc.UtxorpcClient) {
	req := connect.NewRequest(&submit.ReadMempoolRequest{})
	fmt.Println("Connecting to utxorpc host:", client.URL())
	resp, err := client.Submit.ReadMempool(ctx, req)
	if err != nil {
		fmt.Println(connect.CodeOf(err))
		if connectErr := new(connect.Error); errors.As(err, &connectErr) {
			fmt.Println(connectErr.Message())
			fmt.Println(connectErr.Details())
		}
		panic(err)
	}
	fmt.Printf("Response: %+v\n", resp)
}

func waitForTx(ctx context.Context, client *utxorpc.UtxorpcClient) {
	req := connect.NewRequest(&submit.WaitForTxRequest{})
	fmt.Println("Connecting to utxorpc host:", client.URL())
	stream, err := client.Submit.WaitForTx(ctx, req)
	if err != nil {
		fmt.Println("Error initiating stream:", connect.CodeOf(err))
		if connectErr := new(connect.Error); errors.As(err, &connectErr) {
			fmt.Println("Error message:", connectErr.Message())
			fmt.Println("Error details:", connectErr.Details())
		}
		panic(err)
	}

	fmt.Println("Connected to utxorpc host, watching mempool...")
	for stream.Receive() {
		resp := stream.Msg()
		fmt.Printf("Stream response: %+v\n", resp)
	}

	if err := stream.Err(); err != nil {
		fmt.Println("Stream ended with error:", err)
	} else {
		fmt.Println("Stream ended normally.")
	}
}

func watchMempool(ctx context.Context, client *utxorpc.UtxorpcClient) {
	req := connect.NewRequest(&submit.WatchMempoolRequest{})
	fmt.Println("Connecting to utxorpc host:", client.URL())
	stream, err := client.Submit.WatchMempool(ctx, req)
	if err != nil {
		fmt.Println("Error initiating stream:", connect.CodeOf(err))
		if connectErr := new(connect.Error); errors.As(err, &connectErr) {
			fmt.Println("Error message:", connectErr.Message())
			fmt.Println("Error details:", connectErr.Details())
		}
		panic(err)
	}

	fmt.Println("Connected to utxorpc host, watching mempool...")
	for stream.Receive() {
		resp := stream.Msg()
		fmt.Printf("Stream response: %+v\n", resp)
	}

	if err := stream.Err(); err != nil {
		fmt.Println("Stream ended with error:", err)
	} else {
		fmt.Println("Stream ended normally.")
	}
}
