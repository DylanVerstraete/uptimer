package main

import (
	"fmt"

	gsrpc "github.com/centrifuge/go-substrate-rpc-client/v4"
	//"github.com/centrifuge/go-substrate-rpc-client/v4/config"
	//"github.com/centrifuge/go-substrate-rpc-client/v4/signature"
	"github.com/centrifuge/go-substrate-rpc-client/v4/rpc/state"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	//"github.com/gin-gonic/gin"
	"github.com/threefoldtech/substrate-client"
	//"github.com/threefoldtech/tfchain_bridge/pkg"
	//"github.com/threefoldtech/tfchain_bridge/pkg/stellar"
	//subpkg "github.com/threefoldtech/tfchain_bridge/pkg/substrate"
)

func main() {

	//mgrsub, err := subpkg.NewManager("wss://tfchain.dev.grid.tf")
	//fmt.Println(mgrsub)

	mgr := substrate.NewManager("wss://tfchain.dev.grid.tf")
	substrat,err := mgr.Substrate()
	//node, err := substrat.GetNode(1)
	//fmt.Println(node.FarmID)
	//lastnode, err := substrat.GetLastNodeID()
	//fmt.Println(lastnode)
	//farm, err := substrat.GetFarm(uint32(node.FarmID))
	//fmt.Println(farm)

	//fmt.Println(substrat.GetCurrentHeight())
	// The following example shows how to instantiate a Substrate API and use it to connect to a node

	storagekey,storagechange, err := substrat.FetchEventsForBlockRange(1,10)
	fmt.Println(storagekey)
	fmt.Println(storagechange)
	fmt.Println(err)
	//mes := substrat.SubscribeToEvents(1,10)

	cl,_, err := substrat.GetClient()
	chainHeadsSub, err := cl.RPC.Chain.SubscribeFinalizedHeads()
	fmt.Println(chainHeadsSub)

	substratcl,err := NewSubstrateClient("wss://tfchain.dev.grid.tf")
	storagesubs,storagekey, err := substratcl.SubscribeEvents()
	fmt.Println(storagesubs)

	//events,err  := substrat.GetEventsForBlock(4240) 
	//fmt.Println(events)

	//hash,err := substrat.GetUser()

	// "wss://poc3-rpc.polkadot.io"
	api, err := gsrpc.NewSubstrateAPI("wss://tfchain.dev.grid.tf")
	if err != nil {
		panic(err)
	}

	chain, err := api.RPC.System.Chain()
	if err != nil {
		panic(err)
	}
	nodeName, err := api.RPC.System.Name()
	if err != nil {
		panic(err)
	}
	nodeVersion, err := api.RPC.System.Version()
	if err != nil {
		panic(err)
	}

	hash, err := api.RPC.Chain.GetBlockHashLatest()
	if err != nil {
		panic(err)
	}

	block, err := api.RPC.Chain.GetBlock(hash)
	fmt.Println(block)

	fmt.Printf("You are connected to chain %v using %v v%v\n", chain, nodeName, nodeVersion)
	fmt.Printf("Latest block hash: %v\n", hash)


	//node := substr.GetNode()
	//fmt.Printf("Node: %v\n", node)

	// Output: You are connected to chain Development using Substrate Node v3.0.0-dev-1b646b2-x86_64-linux-gnu


}
type SubstrateClient struct {
	*substrate.Substrate
}

func NewSubstrateClient(url string) (*SubstrateClient, error) {
	mngr := substrate.NewManager(url)
	cl, err := mngr.Substrate()
	if err != nil {
		return nil, err
	}

	return &SubstrateClient{
		cl,
	}, nil
}

func (s *SubstrateClient) SubscribeEvents() (*state.StorageSubscription, types.StorageKey, error) {
	cl, meta, err := s.GetClient()
	if err != nil {
		return nil, nil, err
	}

	// Subscribe to system events via storage
	key, err := types.CreateStorageKey(meta, "System", "Events", nil)
	if err != nil {
		return nil, nil, err
	}

	sub, err := cl.RPC.State.SubscribeStorageRaw([]types.StorageKey{key})
	if err != nil {
		return nil, nil, err
	}
	// defer unsubscribe(sub)
	return sub, key, err
}