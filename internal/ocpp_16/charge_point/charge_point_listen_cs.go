package chargepoint

import (
	"sync"
	"time"

	"github.com/Beep-Technologies/beepbeep3-ocpp/internal/ocpp_16/messaging"
)

// callMessageQueue is a simple FIFO queue
type callMessageQueue struct {
	queue []messaging.OCPP16CallMessage
	mutex *sync.Mutex // this mutex controls access to the queue
	cond  *sync.Cond  // this cond is used to signal when there are available messages
}

// enqueue a message to the back of the queue
func (cmq *callMessageQueue) enqueue(msg messaging.OCPP16CallMessage) {
	cmq.mutex.Lock()
	defer cmq.mutex.Unlock()

	cmq.queue = append(cmq.queue, msg)
	cmq.cond.Broadcast()
}

// dequeue a message from the front of the queue. the second return value is false if the queue is empty
func (cmq *callMessageQueue) dequeue() (messaging.OCPP16CallMessage, bool) {
	cmq.mutex.Lock()
	defer cmq.mutex.Unlock()

	if len(cmq.queue) == 0 {
		return messaging.OCPP16CallMessage{}, false
	}

	msg := cmq.queue[0]
	cmq.queue = cmq.queue[1:]

	return msg, true
}

// get the length of the queue
func (cmq *callMessageQueue) len() int {
	cmq.mutex.Lock()
	defer cmq.mutex.Unlock()

	return len(cmq.queue)
}

// listenCS is invoked in the OCPP16ChargePoint.Listen() as a goroutine
// that listens on channels for CS-initiated operations
func (cp *OCPP16ChargePoint) listenCS() {
	// the current outgoing call, that has not been responded to / timed out
	var currentOutgoingCall *messaging.OCPP16CallMessage

	// current outgoing call timer
	var currentOutgoingCallTimer *time.Timer

	// outgoingCallCond is used to signal when currentOutoingCall is nil
	currentOutgoingCallCond := sync.NewCond(&sync.Mutex{})

	// call messages will be placed into the outgoingCallStream when
	// there are no outgoing calls that have been responded to / timed out
	// and there are available messages to send
	outgoingCallStream := make(chan messaging.OCPP16CallMessage)

	// drain the call message queue and place calls into the outgoingCallStream when possible
	go func() {
		for {
			select {
			case <-cp.ctx.Done():
				return
			default:
				// wait until a message can be sent
				currentOutgoingCallCond.L.Lock()
				for currentOutgoingCall != nil {
					currentOutgoingCallCond.Wait()
				}

				// get a message from the callMessageQueue
				cp.callMessageQueue.cond.L.Lock()
				for cp.callMessageQueue.len() == 0 {
					cp.callMessageQueue.cond.Wait()
				}

				outgoingCallMessage, _ := cp.callMessageQueue.dequeue()
				currentOutgoingCall = &outgoingCallMessage
				currentOutgoingCallTimer = time.AfterFunc(10*time.Second, func() {
					currentOutgoingCall = nil
					currentOutgoingCallCond.Broadcast()
				})

				// write to outgoingCallStream
				outgoingCallStream <- outgoingCallMessage

				cp.callMessageQueue.cond.L.Unlock()
				currentOutgoingCallCond.L.Unlock()
			}
		}
	}()

	// while the charge point's context is not cancelled:
	// either:
	// - handle the response to an outgoing call
	// - make an outgoing call
	go func() {
		for {
			select {
			case <-cp.ctx.Done():
				return
			case msg := <-outgoingCallStream:
				cp.outCallStream <- msg
			case <-cp.inCallResultStream:
				currentOutgoingCall = nil
				if currentOutgoingCallTimer != nil {
					currentOutgoingCallTimer.Stop()
				}
				currentOutgoingCallCond.Broadcast()
			case <-cp.inCallErrorStream:
				currentOutgoingCall = nil
				if currentOutgoingCallTimer != nil {
					currentOutgoingCallTimer.Stop()
				}
				currentOutgoingCallCond.Broadcast()
			}
		}
	}()
}
