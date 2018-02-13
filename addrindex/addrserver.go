package addrindex

import (
	"fmt"
	"net/http"

	"github.com/btcsuite/btcd/rpcclient"
	"github.com/gorilla/mux"
)

// AddrServer is the struct where all methods are defined
type AddrServer struct {
	Host       string
	User       string
	Pass       string
	DisableTLS bool
	Port       int
	Client     *rpcclient.Client
}

// AddrServerConfig configures the AddrServer
type AddrServerConfig struct {
	Host string `json:"host"`
	Usr  string `json:"usr"`
	Pass string `json:"pass"`
	SSL  bool   `json:"ssl"`
	Port int    `json:"port"`
}

// NewAddrServer returns a new AddrServer instance
func NewAddrServer(cfg *AddrServerConfig) *AddrServer {
	out := &AddrServer{
		Host:       cfg.Host,
		User:       cfg.Usr,
		Pass:       cfg.Pass,
		DisableTLS: !cfg.SSL,
		Port:       cfg.Port,
	}
	client, err := rpcclient.New(out.connCfg(), nil)
	if err != nil {
		panic(err)
	}
	out.Client = client
	return out
}

// URL returns the backend server's URL
func (as *AddrServer) URL() string {
	if as.DisableTLS {
		return fmt.Sprintf("http://%s:%s@%v", as.User, as.Pass, as.Host)
	}
	return fmt.Sprintf("https://%s:%s@%v", as.User, as.Pass, as.Host)
}

func (as *AddrServer) connCfg() *rpcclient.ConnConfig {
	return &rpcclient.ConnConfig{
		Host:         as.Host,
		User:         as.User,
		Pass:         as.Pass,
		HTTPPostMode: true,
		DisableTLS:   as.DisableTLS,
	}
}

// HANDLERS 👇

// HandleAddrUTXO handles the /addr/<addr>/utxo route
func (as *AddrServer) HandleAddrUTXO(w http.ResponseWriter, r *http.Request) {
	addr := mux.Vars(r)["addr"]

	searchaddress, err := as.SearchRawTransactions(addr, 0, 100)
	if err != nil {
		panic(err)
	}
	utxo := searchaddress.Result.UTXO(addr)
	w.Write(utxo.JSON())
}
