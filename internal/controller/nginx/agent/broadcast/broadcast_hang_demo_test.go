package broadcast_test

import (
	"testing"
	"time"

	. "github.com/onsi/gomega"

	"github.com/nginx/nginx-gateway-fabric/v2/internal/controller/nginx/agent/broadcast"
)

func TestSend_SubscriberNeverResponds(t *testing.T) {
	t.Parallel()
	g := NewWithT(t)

	broadcaster := broadcast.NewDeploymentBroadcaster(t.Context())
	subscriber := broadcaster.Subscribe()

	// Give time for subscription to be processed by the subscriber goroutine
	time.Sleep(10 * time.Millisecond)

	message := broadcast.NginxAgentMessage{
		ConfigVersion: "v1",
		Type:          broadcast.ConfigApplyRequest,
	}

	sendReturned := make(chan struct{})
	go func() {
		broadcaster.Send(message)
		close(sendReturned)
	}()

	// Subscriber receives the message but never responds on ResponseCh
	g.Eventually(subscriber.ListenCh).Should(Receive(Equal(message)))

	// Send should not block indefinitely when a subscriber never responds
	g.Eventually(sendReturned).WithTimeout(2 * time.Second).Should(BeClosed())
}
