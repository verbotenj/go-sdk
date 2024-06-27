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
	"github.com/utxorpc/go-codegen/utxorpc/v1alpha/cardano"
	"github.com/utxorpc/go-codegen/utxorpc/v1alpha/query"
	utxorpc "github.com/utxorpc/go-sdk"
	"golang.org/x/net/http2"
)

func main() {
	// Set mode to "readParams", "readUtxos", "searchUtxos"
	var mode string = "searchUtxos"

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
	case "readParams":
		readParams(ctx, &client)
	case "readUtxos":
		readUtxos(ctx, &client)
	case "searchUtxos":
		searchUtxos(ctx, &client)
	default:
		fmt.Println("Unknown mode:", mode)
	}
}

func readParams(ctx context.Context, client *utxorpc.UtxorpcClient) {
	req := connect.NewRequest(&query.ReadParamsRequest{})

	// req := connect.NewRequest(&submit.SubmitTxRequest{})
	fmt.Println("Connecting to utxorpc host:", client.URL())
	resp, err := client.Query.ReadParams(ctx, req)
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

func readUtxos(ctx context.Context, client *utxorpc.UtxorpcClient) {
	txHash, err := hex.DecodeString("3394533cb02fb71b062690d85bbe9d79a7b6f8f4c1b92b0e728fe7b93a1440c9")
	if err != nil {
		log.Fatalf("failed to decode hex string: %v", err)
	}
	txoRef := &query.TxoRef{
		Hash: txHash,
	}

	req := connect.NewRequest(&query.ReadUtxosRequest{
		Keys: []*query.TxoRef{txoRef},
	})
	fmt.Println("connecting to utxorpc host:", client.URL())
	resp, err := client.Query.ReadUtxos(ctx, req)
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

func searchUtxos(ctx context.Context, client *utxorpc.UtxorpcClient) {

	rawAddress := "806c262cdd383ad3af938eed1c949b31bdfc83db07935c569320bd34343fddde"
	exactAddress, err := hex.DecodeString(rawAddress)
	if err != nil {
		log.Fatalf("failed to decode hex string address: %v", err)
	}
	// Create the AddressPattern
	addressPattern := &cardano.AddressPattern{
		ExactAddress: exactAddress,
	}

	// Create the AssetPattern (if needed)
	assetPattern := &cardano.AssetPattern{
		// Populate the fields as necessary
	}

	// Create the TxOutputPattern with the AddressPattern and AssetPattern
	txOutputPattern := &cardano.TxOutputPattern{
		Address: addressPattern,
		Asset:   assetPattern,
	}

	// Create the AnyUtxoPattern_Cardano with the TxOutputPattern
	anyOutputPattern := &query.AnyUtxoPattern_Cardano{
		Cardano: txOutputPattern,
	}

	// Create the AnyUtxoPattern with the AnyUtxoPattern_Cardano
	anyUtxoPattern := &query.AnyUtxoPattern{
		UtxoPattern: anyOutputPattern,
	}

	// Create the UtxoPredicate with the AnyUtxoPattern
	utxoPredicate := &query.UtxoPredicate{
		Match: anyUtxoPattern,
	}

	// Create the SearchUtxosRequest with the UtxoPredicate
	req := connect.NewRequest(&query.SearchUtxosRequest{
		Predicate: utxoPredicate,
	})

	fmt.Println("connecting to utxorpc host:", client.URL())
	resp, err := client.Query.SearchUtxos(ctx, req)
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
