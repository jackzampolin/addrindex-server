package addrindex

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/btcsuite/btcd/btcjson"
)

var (
	testServer              = "35.227.54.185:8332"
	testPass                = "blockstacksystem"
	testUser                = "blockstack"
	testSSL                 = false
	testPort                = 18332
	testAddress             = "1E2vd1baZLrhNRgU6FPZv1zxxuRNhj4btG"
	testExpectedUTXO        = 1
	testExpectedUTXOSatoshi = 1127408
	testExpectedBalance     = 1127408
	testExpectedReceived    = 9151137107
	testExpectedSent        = 9150009699
	testTransaction         = "a75aceccabf4b44d82cd72fef667a394dc898f59553ddd4356edf4f4bdc267ae"
	testExpectedTxSatoshi   = 990420
	testBlock               = "000000000000000004ec466ce4732fe6f1ed1cddc2ed4b328fff5224276e3f6f"
	testExpectedBlockTx     = 1660
	testExpectedBlockNonce  = uint32(657220870)
	testExpectedBlockHash   = "0100000001ab3e33c13344dde0794a3a1ea03d9f1335a2a08d762b70d943527e9ce9421bc5020000006b483045022100db58090a26c9e3fcbd000c3000210ba7050df65b68e1745b54b695579e2615e00220109b2225eda4c12accf537d16023008121587e61d942494bf02ef3c8667f0a270121023beb2a616698c02d6f64c1ca11a5cccfd722563caccc6a7b5eabd9c6c4158feeffffffff020000000000000000296a2769642bc318e73d4f47f8474028190f1561d9b5319564146d3d4ab1283228a26e82641e96b290ee50b50e00000000001976a9148ef6cfc4c8ee2e6142001c94275ff47d3f05886588ac00000000"
	testBlockIndex          = 400000
)

func handlerTestSetup() (*AddrServer, *httptest.Server) {
	as := NewAddrServer(&AddrServerConfig{
		Host:    testServer,
		Usr:     testUser,
		Pass:    testPass,
		SSL:     testSSL,
		Port:    testPort,
		Version: "test",
		Commit:  "test",
		Branch:  "test",
	})
	return as, httptest.NewServer(as.Router())
}

func TestHandleAddrUTXO(t *testing.T) {
	t.Parallel()
	_, server := handlerTestSetup()
	defer server.Close()
	tempString := "%s/addr/%s/utxo"

	// Make HTTP request to route
	resp, err := http.Get(fmt.Sprintf(tempString, server.URL, testAddress))
	if resp.StatusCode != 200 {
		t.Fatalf("Received non-200 response: %d\n", resp.StatusCode)
	} else if err != nil {
		t.Fatalf("Error making HTTP call to server: %s\n", err.Error())
	}

	// Read the body
	actual, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed Reading Response Body: %s\n", err.Error())
	}

	// Unmarshal the response into proper struct
	resStruct := []UTXOIns{}
	err = json.Unmarshal(actual, &resStruct)
	if err != nil {
		t.Fatalf("Failed Unmarshalling response: %s\n", err.Error())
	}

	// Check response values for accuracy
	expectedUTXO := testExpectedUTXO
	if len(resStruct) != expectedUTXO {
		t.Errorf("Expected '%d' utxo, got '%d'\n", expectedUTXO, len(resStruct))
	}

	// Check response values for accuracy
	if resStruct[0].Satoshis != testExpectedUTXOSatoshi {
		t.Errorf("Expected '%d' satoshi, got '%d'\n", testExpectedUTXOSatoshi, resStruct[0].Satoshis)
	}
}

func TestHandleAddrBalance(t *testing.T) {
	t.Parallel()
	_, server := handlerTestSetup()
	defer server.Close()
	tempString := "%s/addr/%s/balance"

	// Make HTTP request to route
	resp, err := http.Get(fmt.Sprintf(tempString, server.URL, testAddress))
	if resp.StatusCode != 200 {
		t.Fatalf("Received non-200 response: %d\n", resp.StatusCode)
	} else if err != nil {
		t.Fatalf("Error making HTTP call to server: %s\n", err.Error())
	}

	// Read the body
	actual, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed Reading Response Body: %s\n", err.Error())
	}

	// Unmarshal the response into proper struct
	retVal, err := strconv.ParseInt(string(actual), 10, 64)
	if err != nil {
		t.Fatalf("Failed Parsing Balance: %s\n", err.Error())
	}

	// Check response values for accuracy
	if int(retVal) != testExpectedBalance {
		t.Errorf("Expected balance '%d', got '%d'\n", testExpectedBalance, retVal)
	}
}

func TestHandleAddrRecieved(t *testing.T) {
	t.Parallel()
	_, server := handlerTestSetup()
	defer server.Close()
	tempString := "%s/addr/%s/totalReceived"

	// Make HTTP request to route
	resp, err := http.Get(fmt.Sprintf(tempString, server.URL, testAddress))
	if resp.StatusCode != 200 {
		t.Fatalf("Received non-200 response: %d\n", resp.StatusCode)
	} else if err != nil {
		t.Fatalf("Error making HTTP call to server: %s\n", err.Error())
	}

	// Read the body
	actual, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed Reading Response Body: %s\n", err.Error())
	}

	// Unmarshal the response into proper struct
	retVal, err := strconv.ParseInt(string(actual), 10, 64)
	if err != nil {
		t.Fatalf("Failed Parsing Balance: %s\n", err.Error())
	}

	// Check response values for accuracy
	if int(retVal) != testExpectedReceived {
		t.Errorf("Expected balance '%d', got '%d'\n", testExpectedReceived, retVal)
	}
}

func TestHandleAddrSent(t *testing.T) {
	t.Parallel()
	_, server := handlerTestSetup()
	defer server.Close()
	tempString := "%s/addr/%s/totalSent"

	// Make HTTP request to route
	resp, err := http.Get(fmt.Sprintf(tempString, server.URL, testAddress))
	if resp.StatusCode != 200 {
		t.Fatalf("Received non-200 response: %d\n", resp.StatusCode)
	} else if err != nil {
		t.Fatalf("Error making HTTP call to server: %s\n", err.Error())
	}

	// Read the body
	actual, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed Reading Response Body: %s\n", err.Error())
	}

	// Unmarshal the response into proper struct
	retVal, err := strconv.ParseInt(string(actual), 10, 64)
	if err != nil {
		t.Fatalf("Failed Parsing Balance: %s\n", err.Error())
	}

	// Check response values for accuracy
	if int(retVal) != testExpectedSent {
		t.Errorf("Expected balance '%d', got '%d'\n", testExpectedSent, retVal)
	}
}

func TestHandleTxGet(t *testing.T) {
	t.Parallel()
	_, server := handlerTestSetup()
	defer server.Close()
	tempString := "%s/tx/%s"

	// Make HTTP request to route
	resp, err := http.Get(fmt.Sprintf(tempString, server.URL, testTransaction))
	if resp.StatusCode != 200 {
		t.Fatalf("Received non-200 response: %d\n", resp.StatusCode)
	} else if err != nil {
		t.Fatalf("Error making HTTP call to server: %s\n", err.Error())
	}

	// Read the body
	actual, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed Reading Response Body: %s\n", err.Error())
	}

	// Unmarshal the response into proper struct
	resStruct := TransactionIns{}
	err = json.Unmarshal(actual, &resStruct)
	if err != nil {
		t.Fatalf("Failed Unmarshalling response: %s\n", err.Error())
	}

	// Check response values for accuracy
	if resStruct.Vin[0].ValueSat != testExpectedTxSatoshi {
		t.Errorf("Expected '%d' satoshi, got '%d'\n", testExpectedTxSatoshi, resStruct.Vin[0].ValueSat)
	}
}

func TestHandleRawTxGet(t *testing.T) {
	t.Parallel()
	_, server := handlerTestSetup()
	defer server.Close()
	tempString := "%s/rawtx/%s"

	// Make HTTP request to route
	resp, err := http.Get(fmt.Sprintf(tempString, server.URL, testTransaction))
	if resp.StatusCode != 200 {
		t.Fatalf("Received non-200 response: %d\n", resp.StatusCode)
	} else if err != nil {
		t.Fatalf("Error making HTTP call to server: %s\n", err.Error())
	}

	// Read the body
	actual, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed Reading Response Body: %s\n", err.Error())
	}

	// Unmarshal the response into proper struct
	resStruct := map[string]string{}
	err = json.Unmarshal(actual, &resStruct)
	if err != nil {
		t.Fatalf("Failed Unmarshalling response: %s\n", err.Error())
	}

	// Check response values for accuracy
	if resStruct["rawtx"] != testExpectedBlockHash {
		t.Errorf("Expected '%s' satoshi, got '%s'\n", testExpectedBlockHash, resStruct["rawtx"])
	}
}

func TestHandleGetBlock(t *testing.T) {
	t.Parallel()
	_, server := handlerTestSetup()
	defer server.Close()
	tempString := "%s/block/%s"

	// Make HTTP request to route
	resp, err := http.Get(fmt.Sprintf(tempString, server.URL, testBlock))
	if resp.StatusCode != 200 {
		t.Fatalf("Received non-200 response: %d\n", resp.StatusCode)
	} else if err != nil {
		t.Fatalf("Error making HTTP call to server: %s\n", err.Error())
	}

	// Read the body
	actual, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed Reading Response Body: %s\n", err.Error())
	}

	// Unmarshal the response into proper struct
	resStruct := btcjson.GetBlockVerboseResult{}
	err = json.Unmarshal(actual, &resStruct)
	if err != nil {
		t.Fatalf("Failed Unmarshalling response: %s\n", err.Error())
	}

	// Check response values for accuracy
	if len(resStruct.Tx) != testExpectedBlockTx {
		t.Errorf("Expected '%d' utxo, got '%d'\n", testExpectedBlockTx, len(resStruct.Tx))
	}

	// Check response values for accuracy
	if resStruct.Nonce != testExpectedBlockNonce {
		t.Errorf("Expected '%d' satoshi, got '%d'\n", testExpectedBlockNonce, resStruct.Nonce)
	}
}

func TestHandleGetBlockHash(t *testing.T) {
	t.Parallel()
	_, server := handlerTestSetup()
	defer server.Close()
	tempString := "%s/block-index/%d"

	// Make HTTP request to route
	resp, err := http.Get(fmt.Sprintf(tempString, server.URL, testBlockIndex))
	if resp.StatusCode != 200 {
		t.Fatalf("Received non-200 response: %d\n", resp.StatusCode)
	} else if err != nil {
		t.Fatalf("Error making HTTP call to server: %s\n", err.Error())
	}

	// Read the body
	actual, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed Reading Response Body: %s\n", err.Error())
	}

	// Unmarshal the response into proper struct
	resStruct := map[string]string{}
	err = json.Unmarshal(actual, &resStruct)
	if err != nil {
		t.Fatalf("Failed Unmarshalling response: %s\n", err.Error())
	}

	// Check response values for accuracy
	if resStruct["blockHash"] != testBlock {
		t.Errorf("Expected '%s' satoshi, got '%s'\n", testBlock, resStruct["blockHash"])
	}
}
