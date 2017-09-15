package gohelprabbitmq

import (
	"fmt"
	"testing"
	"time"

	. "github.com/Cappta/gofixture"
	"github.com/Cappta/gohelpgabs"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	rabbitMQ             = ConnectLocallyWithDefaultUser()
	testServiceConsumer  = NewAsyncConsumer(rabbitMQ, "TestForRPC")
	testServicePublisher = NewSimplePublisher(rabbitMQ, fmt.Sprintf("@%s", testServiceConsumer.QueueSettings.GetName()))
	rpc                  = NewRPC(NewSimpleConsumer(rabbitMQ, "CallbackToRPC"), time.Second*5)
	messageRouter        = NewMessageRouter("TestForRPC.ForwardTo", rabbitMQ)
)

func TestRPC(t *testing.T) {
	done := make(chan bool)
	go func() {
		testServiceConsumer.Consume()
		done <- true
	}()
	go rpc.Consume()

	err := NewQueueObserver(rabbitMQ, testServiceConsumer.QueueSettings).AwaitConsumer(time.Second * 5)
	if err != nil {
		panic(err)
	}
	err = NewQueueObserver(rabbitMQ, rpc.QueueSettings).AwaitConsumer(time.Second * 5)
	if err != nil {
		panic(err)
	}

	Convey("Given a running async consumer", t, func() {
		Convey("Given any string with length between 10 and 100", func() {
			anyString := AnyString(AnyIntBetween(10, 100))
			Convey("Given a json container with the string set as value", func() {
				container := gohelpgabs.New()
				path := "Value"
				container.SetP(anyString, path)
				Convey("When preparing the rpc", func() {
					callback := rpc.Prepare(container)
					Convey("Then callback should not be nil", func() {
						So(callback, ShouldNotBeNil)
						Convey("When publishing to the test service", func() {
							testServicePublisher.Publish(container.Bytes())
							Convey("Then the test service should receive the message within five seconds", func() {
								message, err := testServiceConsumer.AwaitDelivery(time.Second * 5)
								Convey("Then error should be nil", func() {
									So(err, ShouldBeNil)
									Convey("Then message should not be nil", func() {
										So(message, ShouldNotBeNil)
										Convey("Then message should be a valid json", func() {
											messageContainer, err := gohelpgabs.ParseJSON(message.Body)
											So(err, ShouldBeNil)
											Convey("Then value's path should exist", func() {
												So(messageContainer.ExistsP(path), ShouldBeTrue)
												Convey("Then returned value should equal provided value", func() {
													So(messageContainer.Path(path).Data().(string), ShouldEqual, anyString)
													Convey("When routing the message", func() {
														err = messageRouter.Route(messageContainer)
														Convey("Then error should be nil", func() {
															So(err, ShouldBeNil)
															Convey("Then rpc should receive the message within 5 seconds", func() {
																select {
																case <-time.After(time.Second * 5):
																	t.Fail()
																case returnedMessage := <-callback:
																	Convey("Then message should not be nil", func() {
																		So(returnedMessage, ShouldNotBeNil)
																		Convey("Then message should be a valid json", func() {
																			returnedMessageContainer, err := gohelpgabs.ParseJSON(returnedMessage.Body)
																			So(err, ShouldBeNil)
																			Convey("Then value's path should exist", func() {
																				So(returnedMessageContainer.ExistsP(path), ShouldBeTrue)
																				Convey("Then returned value should equal provided value", func() {
																					So(returnedMessageContainer.Path(path).Data().(string), ShouldEqual, anyString)
																					Convey("Then RPC's path should not exist", func() {
																						So(returnedMessageContainer.Exists(rpcPath), ShouldBeFalse)
																					})
																				})
																			})
																		})
																	})
																}
															})
														})
													})
												})
											})
										})
									})
								})
							})
						})
					})
				})
			})
		})
	})

	err = testServiceConsumer.StopConsuming()
	if err != nil {
		panic(err)
	}
	err = rpc.StopConsuming()
	if err != nil {
		panic(err)
	}
	<-done
}
