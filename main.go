package main

import (
	"fmt"
	"time"

	ds "github.com/jbenet/go-ipfs/Godeps/_workspace/src/github.com/jbenet/go-datastore"
	dsync "github.com/jbenet/go-ipfs/Godeps/_workspace/src/github.com/jbenet/go-datastore/sync"
	ma "github.com/jbenet/go-ipfs/Godeps/_workspace/src/github.com/jbenet/go-multiaddr"
	basic "github.com/jbenet/go-ipfs/p2p/host/basic"
	"github.com/jbenet/go-ipfs/p2p/net/swarm"
	"github.com/jbenet/go-ipfs/p2p/peer"
	"github.com/jbenet/go-ipfs/routing/dht"
	tu "github.com/jbenet/go-ipfs/util/testutil"
	"golang.org/x/net/context"
)

func main() {
	ps := peer.NewPeerstore()
	id, err := tu.RandIdentity()
	if err != nil {
		panic(err)
	}

	// Set own keys in peerstore
	ps.AddPrivKey(id.ID(), id.PrivateKey())
	ps.AddPubKey(id.ID(), id.PrivateKey().GetPublic())

	addr := ma.StringCast("/ip4/0.0.0.0/tcp/9000")
	listenaddrs := []ma.Multiaddr{addr}

	fmt.Printf("Hello! My peer id is %s\n", id.ID().Pretty())
	fmt.Printf("My multiaddr is %s/ipfs/%s\n", addr.String(), id.ID().Pretty())
	network, err := swarm.NewNetwork(context.Background(), listenaddrs, id.ID(), ps)
	if err != nil {
		panic(err)
	}

	mars, err := peer.IDB58Decode("QmaCpDMGvV2BGHeYERUEnRQAwe3N8SzbUtfsmvsqQLuvuJ")
	if err != nil {
		panic(err)
	}

	fmt.Println("Bootstrapping with mars...")
	marsaddr := ma.StringCast("/ip4/104.131.131.82/tcp/4001")
	ps.AddAddr(mars, marsaddr, time.Hour)
	network.DialPeer(context.Background(), mars)

	dstore := ds.NewMapDatastore()
	tsds := dsync.MutexWrap(dstore)

	host := basic.New(network, basic.NATPortMap)

	fmt.Println("firing up dht...")
	_ = dht.NewDHT(context.Background(), host, tsds)

	fmt.Println("Ready for anything!")
	<-context.Background().Done()
}
