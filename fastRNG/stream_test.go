///////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////
package fastRNG

import (
	"fmt"
	"gitlab.com/elixxir/crypto/csprng"
	"reflect"
	"testing"
	"time"
)

//Test the creation of a new stream generator and that it is configured correctly
func TestNewStreamGenerator(t *testing.T) {
	sg := NewStreamGenerator(csprng.NewSystemRNG(), 12, 20)
	if sg.maxStreams != 20 || sg.scalingFactor != 12 || !reflect.DeepEqual(sg.rng, csprng.NewSystemRNG() ){
		t.Errorf("Failure to initialize a stream generator correctly")
	}
}

//Test the creation of new streams and that the counters are in fact working
func TestNewStream(t *testing.T) {
	sg := NewStreamGenerator(csprng.NewSystemRNG(), 12, 3)
	sg.NewStream()
	sg.NewStream()
	sg.NewStream()
	//See if there are the appropriate amount of streams in the streams slice and the stream count
	if sg.numStreams != uint(len(sg.streams)) && sg.numStreams != 3 {
		t.Errorf("New streams bookkeeping is not working.")
	}
}

//Test that the stream generator panics when there are too many streams being made
func TestNewStream_DoesPanic(t *testing.T) {
	//The defer function will catch the panic
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("FastRNG should panic when too many streams are made!")
		}
	}()
	//Stream count is 3, but 4 streams are being created, thus it should panic
	sg := NewStreamGenerator(csprng.NewSystemRNG(), 12, 3)
	sg.NewStream()
	sg.NewStream()
	sg.NewStream()
	sg.NewStream()

	t.Errorf("Fast RNG should have panicked here, too many streams were created")

}

//Test that it does not panic when it reaches capacity
func TestNewStream_NotPanic(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("FastRNG should not panic when there are exactly max streams " +
				"made!")
		}
	}()
	//Stream count is 3, and 3 streams are being created, thus it should not panic
	sg := NewStreamGenerator(csprng.NewSystemRNG(), 12, 3)
	sg.NewStream()
	sg.NewStream()
	sg.NewStream()
}

//Test that the getStream calls newStream correctly/appropriately
func TestGetStream_NewStream(t *testing.T) {
	sg := NewStreamGenerator(csprng.NewSystemRNG(), 12, 3)
	sg.GetStream()
	sg.GetStream()
	sg.GetStream()

	if sg.numStreams != uint(len(sg.streams)) && sg.numStreams != 3 {
		t.Errorf("New streams bookkeeping is not working.")
	}
}

//Test that a blocked channel will grab a stream when it becomes available
func TestGetStream_GrabsWaitingStream(t *testing.T) {
	sg := NewStreamGenerator(csprng.NewSystemRNG(), 12, 3)
	stream0 := sg.GetStream()
	sg.GetStream()
	sg.GetStream()
	//Allow the main thread to block as streams aren't available, then close it
	go func() {
		time.Sleep(500 * time.Millisecond)
		sg.Close(stream0)
	}()
	newStream := sg.GetStream()
	if !reflect.DeepEqual(newStream, stream0) {
		t.Errorf("The next stream did not grab the correct stream")
	}
}

func TestClose_WaitingChannelLength(t *testing.T) {
	sg := NewStreamGenerator(csprng.NewSystemRNG(), 12, 3)
	stream0 := sg.GetStream()
	stream1 := sg.GetStream()
	stream2 := sg.GetStream()

	//Close all the streams created
	sg.Close(stream0)
	sg.Close(stream1)
	sg.Close(stream2)

	//Check that the waiting streams channel is the appropriate length
	if len(sg.waitingStreams) != 3 {
		t.Errorf("Waiting channel isn't the appropriate size after closing streams")
	}
}
