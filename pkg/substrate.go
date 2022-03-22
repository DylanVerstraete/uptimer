package pkg

import (

	//"github.com/centrifuge/go-substrate-rpc-client/v4/config"
	//"github.com/centrifuge/go-substrate-rpc-client/v4/signature"
	"context"
	"fmt"

	"github.com/centrifuge/go-substrate-rpc-client/v4/rpc/state"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/pkg/errors"

	//"github.com/gin-gonic/gin"
	"github.com/threefoldtech/substrate-client"
	//"github.com/threefoldtech/tfchain_bridge/pkg"
	//"github.com/threefoldtech/tfchain_bridge/pkg/stellar"
	//subpkg "github.com/threefoldtech/tfchain_bridge/pkg/substrate"
)

// CONFIG STRUCT
// WSURL, DBURL, ...

type Uptimer struct {
	subClient *substrate.Substrate
	// TODO
	// MONGO CONNECTION
}

func Start(ctx context.Context) error {
	//mgrsub, err := subpkg.NewManager("wss://tfchain.dev.grid.tf")
	//fmt.Println(mgrsub)

	mgr := substrate.NewManager("wss://tfchain.dev.grid.tf")
	substrat, err := mgr.Substrate()
	if err != nil {
		return err
	}

	uptimer := Uptimer{
		subClient: substrat,
	}

	cl, _, err := uptimer.subClient.GetClient()
	if err != nil {
		return errors.Wrap(err, "failed to get client")
	}

	chainHeadsSub, err := cl.RPC.Chain.SubscribeFinalizedHeads()
	if err != nil {
		return errors.Wrap(err, "failed to subscribe to finalized heads")
	}

	for {
		select {
		case head := <-chainHeadsSub.Chan():
			err := uptimer.processEventsForHeight(uint32(head.Number))
			if err != nil {
				return err
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (uptimer *Uptimer) processEventsForHeight(height uint32) error {
	// log.Info().Msgf("fetching events for blockheight %d", height)
	records, err := uptimer.subClient.GetEventsForBlock(height)
	if err != nil {
		// log.Info().Msgf("failed to decode block with height %d", height)
		return err
	}

	for _, e := range records.TfgridModule_NodeUptimeReported {
		// TODO, SERIALIZE INTO OBJECTS -> SAVE IN MONGO
		fmt.Println(e.Node)
		fmt.Println(e.Timestamp)
		fmt.Println(e.Uptime)
	}

	return nil
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

// CLIENT STUFF
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
