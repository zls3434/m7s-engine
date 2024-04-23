package common

import (
	"time"

	"github.com/zls3434/m7s-engine/v4/config"
	"github.com/zls3434/m7s-engine/v4/log"
	"github.com/zls3434/m7s-engine/v4/util"
)

type IStream interface {
	AddTrack(Track) *util.Promise[Track]
	RemoveTrack(Track)
	Close()
	IsClosed() bool
	SSRC() uint32
	log.Zap
	Receive(any) bool
	SetIDR(Track)
	GetPublisherConfig() *config.Publish
	GetStartTime() time.Time
	GetType() string
}
