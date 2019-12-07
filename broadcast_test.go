package socket

import (
	"github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestBroadcast_Join(t *testing.T) {
	bc := newDefaultBroadcast()

	convey.Convey("testing broadcast join", t, func() {
		err := bc.Join(DefaultBroadcastRoomName, &fakeServerClient{id: "tests"})
		convey.So(err, convey.ShouldBeNil)
		convey.So(bc.Len(DefaultBroadcastRoomName), convey.ShouldEqual, 1)
	})

	convey.Convey("testing broadcast join, duplicate client id", t, func() {
		err := bc.Join(DefaultBroadcastRoomName, &fakeServerClient{"tests"})
		convey.So(err, convey.ShouldEqual, ErrClientInRoomAlreadyExist)
		convey.So(bc.Len(DefaultBroadcastRoomName), convey.ShouldEqual, 1)
	})

	convey.Convey("testing broadcast join", t, func() {
		err := bc.Join(DefaultBroadcastRoomName, &fakeServerClient{"tests1"})
		convey.So(err, convey.ShouldBeNil)
		convey.So(bc.Len(DefaultBroadcastRoomName), convey.ShouldEqual, 2)
	})

	convey.Convey("testing broadcast join", t, func() {
		err := bc.Join("test", &fakeServerClient{id: "tests"})
		convey.So(err, convey.ShouldBeNil)
		convey.So(bc.Len("test"), convey.ShouldEqual, 1)
	})
}

func TestBroadcast_Leave(t *testing.T) {
	bc := newDefaultBroadcast()

	convey.Convey("testing broadcast leave", t, func() {
		err := bc.Join(DefaultBroadcastRoomName, &fakeServerClient{id: "tests"})
		convey.So(err, convey.ShouldBeNil)
		convey.So(bc.Len(DefaultBroadcastRoomName), convey.ShouldEqual, 1)

		err = bc.Leave(DefaultBroadcastRoomName, &fakeServerClient{id: "error"})
		convey.So(err, convey.ShouldEqual, ErrClientInRoomNotExist)
		convey.So(bc.Len(DefaultBroadcastRoomName), convey.ShouldEqual, 1)

		err = bc.Leave(DefaultBroadcastRoomName, &fakeServerClient{id: "tests"})
		convey.So(err, convey.ShouldBeNil)
		convey.So(bc.Len(DefaultBroadcastRoomName), convey.ShouldEqual, -1)

		err = bc.Leave("qwe", &fakeServerClient{id: "tests"})
		convey.So(err, convey.ShouldEqual, ErrRoomNotExist)
		convey.So(bc.Len(DefaultBroadcastRoomName), convey.ShouldEqual, -1)
	})
}

func TestBroadcast_Send(t *testing.T) {
	bc := newDefaultBroadcast()

	convey.Convey("testing broadcast send", t, func() {
		err := bc.Join(DefaultBroadcastRoomName, &fakeServerClient{id: "tests1"})
		convey.So(err, convey.ShouldBeNil)
		convey.So(bc.Len(DefaultBroadcastRoomName), convey.ShouldEqual, 1)

		err = bc.Join(DefaultBroadcastRoomName, &fakeServerClient{id: "tests2"})
		convey.So(err, convey.ShouldBeNil)
		convey.So(bc.Len(DefaultBroadcastRoomName), convey.ShouldEqual, 2)

		err = bc.Join(DefaultBroadcastRoomName, &fakeServerClient{id: "tests3"})
		convey.So(err, convey.ShouldBeNil)
		convey.So(bc.Len(DefaultBroadcastRoomName), convey.ShouldEqual, 3)

		err = bc.Send(&fakeServerClient{id: "tests1"}, DefaultBroadcastRoomName, "test", nil)
		convey.So(err, convey.ShouldBeNil)

		err = bc.Send(&fakeServerClient{id: "tests1"}, "test", "test", nil)
		convey.So(err, convey.ShouldEqual, ErrRoomNotExist)

		err = bc.Send(nil, DefaultBroadcastRoomName, "test", nil)
		convey.So(err, convey.ShouldBeNil)
	})
}
