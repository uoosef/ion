package sfu

import (
	nprotoo "github.com/cloudwebrtc/nats-protoo"
	"github.com/pion/ion/pkg/log"
	"github.com/pion/ion/pkg/proto"
	"github.com/pion/ion/pkg/rtc"
	"github.com/pion/ion/pkg/util"
)

var (
	dc          = "default"
	nid         = "sfu-unkown-node-id"
	protoo      *nprotoo.NatsProtoo
	broadcaster *nprotoo.Broadcaster
)

// Init func
func Init(dcID, nodeID, rpcID, eventID, natsURL string) {
	dc = dcID
	nid = nodeID
	protoo = nprotoo.NewNatsProtoo(natsURL)
	broadcaster = protoo.NewBroadcaster(eventID)
	handleRequest(rpcID)
	checkRTC()
}

// checkRTC send `stream-remove` msg to islb when some pub has been cleaned
func checkRTC() {
	log.Infof("SFU.checkRTC start")
	go func() {
		for mInfo := range rtc.CleanChannel {
			mediaInfo, err := proto.ParseMediaInfo(mInfo)
			if err != nil {
				log.Errorf("Error parsing media info %v", mInfo)
				continue
			}
			broadcaster.Say(proto.SFUStreamRemove, util.Map("mid", mediaInfo.MID))
		}
	}()
}
