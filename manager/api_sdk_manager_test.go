package manager

import (
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"testing"
)

func TestGetSDKWithOauth2Authentication(t *testing.T) {

	Convey("Given I have a request with a passthrough token in the header", t, func() {

		req, _ := http.NewRequest("Get", "foo", nil)
		req.Header.Add("Eric-Access-Token", `{"token_type":"bearer","access_token":"bar","expires_in":1234}`)

		Convey("When I try to retrieve an instance of the SDK", func() {
			sdkManager := Get()
			service, err := sdkManager.GetSDK(req, true)

			Convey("Then no errors should be returned", func() {

				So(err, ShouldBeNil)

				Convey("And the SDK service should not be nil", func() {

					So(service, ShouldNotBeNil)
				})
			})
		})
	})
}

func TestGetSDKWithApiKeyAuthentication(t *testing.T) {

	Convey("Given I have a request without a passthrough token in the header ", t, func() {

		req, _ := http.NewRequest("Get", "foo", nil)

		Convey("When I try to retrieve an instance of the SDK", func() {
			sdkManager := Get()
			service, err := sdkManager.GetSDK(req, false)

			Convey("Then no errors should be returned", func() {

				So(err, ShouldBeNil)

				Convey("And the SDK service should not be nil", func() {

					So(service, ShouldNotBeNil)
				})
			})
		})
	})
}

func TestGetPrivateSDKWithOauth2Authentication(t *testing.T) {

	Convey("Given I have a request with a passthrough token in the header", t, func() {

		req, _ := http.NewRequest("Get", "foo", nil)
		req.Header.Add("Eric-Access-Token", `{"token_type":"bearer","access_token":"bar","expires_in":1234}`)

		Convey("When I try to retrieve an instance of the SDK", func() {
			sdkManager := Get()
			service, err := sdkManager.GetPrivateSDK(req, true)

			Convey("Then no errors should be returned", func() {

				So(err, ShouldBeNil)

				Convey("And the SDK service should not be nil", func() {

					So(service, ShouldNotBeNil)
				})
			})
		})
	})
}

func TestGetPrivateSDKWithApiKeyAuthentication(t *testing.T) {

	Convey("Given I have a request without a passthrough token in the header ", t, func() {

		req, _ := http.NewRequest("Get", "foo", nil)

		Convey("When I try to retrieve an instance of the SDK", func() {
			sdkManager := Get()
			service, err := sdkManager.GetPrivateSDK(req, false)

			Convey("Then no errors should be returned", func() {

				So(err, ShouldBeNil)

				Convey("And the SDK service should not be nil", func() {

					So(service, ShouldNotBeNil)
				})
			})
		})
	})
}
