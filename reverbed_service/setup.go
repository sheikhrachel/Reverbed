package reverbed_service

import (
	"sync"
	"time"

	"github.com/integralist/go-elasticache/elasticache"

	aws "github.com/sheikhrachel/reverbed/api_common/aws_utils"
	"github.com/sheikhrachel/reverbed/api_common/call"
	"github.com/sheikhrachel/reverbed/api_common/utils/errutil"
)

type ReverbedService struct {
	aws *aws.AWSClient
	// Mutex to protect shared resources
	mu              sync.Mutex
	memcachedClient *elasticache.Client
}

func New(cc call.Call, awsClient *aws.AWSClient) *ReverbedService {
	service := &ReverbedService{aws: awsClient}
	service.setupMemcached(cc)
	return service
}

const (
	itemTTL int32 = 300
)

// setupMemcached sets up the memcached client by reading the ELASTICACHE_ENDPOINT environment variable
func (r *ReverbedService) setupMemcached(cc call.Call) {
	var err error
	done := make(chan bool)
	go func() {
		r.memcachedClient, err = elasticache.New()
		done <- true
	}()
	select {
	case <-done:
		if errutil.HandleError(cc, err) {
			return
		}
		cc.InfoF("Successfully set up memcached client")
	case <-time.After(3 * time.Second):
	}
}
