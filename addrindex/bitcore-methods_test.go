package addrindex

import (
	"testing"
)

var (
	testExpectedAddrTxns    = 4983
	testExpectedAddrDeltas  = 9965
	startBlock              = 300000
	endBlock                = 500000
	testBlockEndTime        = 1519237038
	testBlockStartTime      = 1519064238
	testExpectedBlockHashes = 311
	testTransactionIndex    = 1
	testExpectedTxid        = "b3922b88ba526df9cab9634785892de245004c96c36ede9b5b50f68abe584e98"
)

func bitcoreTestSetup() *AddrServer {
	return NewAddrServer(&AddrServerConfig{
		Host:    testServer,
		Usr:     testUser,
		Pass:    testPass,
		SSL:     testSSL,
		Port:    testPort,
		Version: "test",
		Commit:  "test",
		Branch:  "test",
	})
}

func TestGetAddressTxIDs(t *testing.T) {
	t.Parallel()
	as := bitcoreTestSetup()
	out, err := as.GetAddressTxIDs([]string{testAddress}, startBlock, endBlock)

	if err != nil {
		t.Fatal(err)
	}

	if len(out.Result) != testExpectedAddrTxns {
		t.Fatalf("Expected '%d' transactions, got '%d'\n", testExpectedAddrTxns, len(out.Result))
	}
}

func TestGetAddressDeltas(t *testing.T) {
	t.Parallel()
	as := bitcoreTestSetup()
	out, err := as.GetAddressDeltas([]string{testAddress}, startBlock, endBlock)

	if err != nil {
		t.Fatal(err)
	}

	if len(out.Result) != testExpectedAddrDeltas {
		t.Fatalf("Expected '%d' deltas, got '%d'\n", testExpectedAddrDeltas, len(out.Result))
	}
}

func TestGetAddressBalance(t *testing.T) {
	t.Parallel()
	as := bitcoreTestSetup()
	out, err := as.GetAddressBalance([]string{testAddress})

	if err != nil {
		t.Fatal(err)
	}

	if out.Result.Balance != testExpectedBalance {
		t.Fatalf("Expected '%d' satoshi, got '%d'\n", testExpectedBalance, out.Result.Balance)
	}
}

func TestGetAddressUTXOs(t *testing.T) {
	t.Parallel()
	as := bitcoreTestSetup()
	out, err := as.GetAddressUTXOs([]string{testAddress})

	if err != nil {
		t.Fatal(err)
	}

	if len(out.Result) != testExpectedUTXO {
		t.Fatalf("Expected '%d' utxo, got '%d'\n", testExpectedUTXO, len(out.Result))
	}

	if out.Result[0].Satoshis != testExpectedUTXOSatoshi {
		t.Fatalf("Expected '%d' satoshi, got '%d'\n", testExpectedUTXO, len(out.Result))
	}
}

func TestGetAddressMempool(t *testing.T) {
	t.Parallel()
	as := bitcoreTestSetup()

	// Fetch the raw mempool from the server
	mp, err := as.Client.GetRawMempool()
	if err != nil {
		t.Fatal(err)
	}

	// Fetch details for one of those transactions
	tx := mp[0].String()
	raw, err := as.GetRawTransaction(tx)
	if err != nil {
		t.Fatal(err)
	}

	// Run GetAddressMempool with the resulting address
	address := raw.Result.Vout[0].ScriptPubKey.Addresses[0]
	out, err := as.GetAddressMempool([]string{address})
	if err != nil {
		t.Fatal(err)
	}

	// Check results for instances of that address
	ok := false
	for _, tx := range out.Result {
		if tx.Address == address {
			ok = true
		}
	}

	if !ok {
		t.Fatalf("Expected to find address %s from transaction %s in mempool. That address was not found.", address, tx)
	}
}

func TestGetBlockHashes(t *testing.T) {
	t.Parallel()
	as := bitcoreTestSetup()

	hashes, err := as.GetBlockHashes(testBlockEndTime, testBlockStartTime)
	if err != nil {
		t.Fatal(err)
	}

	if len(hashes.Result) != testExpectedBlockHashes {
		t.Fatalf("Expected to find %d block hashes between %d and %d, but found %d", testExpectedBlockHashes, testBlockEndTime, testBlockStartTime, len(hashes.Result))
	}
}

func TestGetSpentInfo(t *testing.T) {
	t.Parallel()
	as := bitcoreTestSetup()

	spent, err := as.GetSpentInfo(testTransaction, testTransactionIndex)
	if err != nil {
		t.Fatal(err)
	}

	if spent.Result.Txid != testExpectedTxid {
		t.Fatalf("Expected output 0 from tx %s to have been spent by %s, but was spent by %s", testTransaction, testExpectedTxid, spent.Result.Txid)
	}
}

func TestGetRawTransaction(t *testing.T) {
	t.Parallel()
	as := bitcoreTestSetup()
	txn, err := as.GetRawTransaction(testTransaction)
	if err != nil {
		t.Fatal(err)
	}

	if txn.Result.Vin[0].ValueSat != testExpectedTxSatoshi {
		t.Errorf("Expected '%d' satoshi in vin 0, got '%d'\n", testExpectedTxSatoshi, txn.Result.Vin[0].ValueSat)
	}
}
