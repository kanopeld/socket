package socket

import (
	"github.com/smartystreets/goconvey/convey"
	"testing"
)

var testRoom Room = getRoom()

func TestRoom_SetClient(t *testing.T) {
	convey.Convey("testing room set client. regular case", t, func() {
		err := testRoom.SetClient(&fakeServerClient{Id: "tests"})
		convey.So(err, convey.ShouldBeNil)
		convey.So(testRoom.Len(), convey.ShouldEqual, 1)
	})

	convey.Convey("testing room set client, duplicate client id", t, func() {
		err := testRoom.SetClient(&fakeServerClient{Id: "tests"})
		convey.So(err, convey.ShouldEqual, ErrClientInRoomAlreadyExist)
		convey.So(testRoom.Len(), convey.ShouldEqual, 1)
	})

	convey.Convey("testing broadcast join", t, func() {
		err := testRoom.SetClient(&fakeServerClient{Id: "tests1"})
		convey.So(err, convey.ShouldBeNil)
		convey.So(testRoom.Len(), convey.ShouldEqual, 2)
	})
}

func TestRoom_RemoveClient(t *testing.T) {
	_ = testRoom.SetClient(&fakeServerClient{Id: "tests"})
	_ = testRoom.SetClient(&fakeServerClient{Id: "tests1"})

	convey.Convey("testing room set client. regular case", t, func() {
		err := testRoom.RemoveClient(&fakeServerClient{Id: "tests"})
		convey.So(err, convey.ShouldBeNil)
		convey.So(testRoom.Len(), convey.ShouldEqual, 1)
	})

	convey.Convey("testing room set client, duplicate client id", t, func() {
		err := testRoom.RemoveClient(&fakeServerClient{Id: "tests"})
		convey.So(err, convey.ShouldEqual, ErrClientInRoomNotExist)
		convey.So(testRoom.Len(), convey.ShouldEqual, 1)
	})

	convey.Convey("testing broadcast join", t, func() {
		err := testRoom.RemoveClient(&fakeServerClient{Id: "tests1"})
		convey.So(err, convey.ShouldBeNil)
		convey.So(testRoom.Len(), convey.ShouldEqual, 0)
	})
}
