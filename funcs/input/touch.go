// +build darwin,!arm,!arm64 freebsd linux windows js
// +build !android
// +build !ios

// PCでビルドすると、isTouchPrimaryInputがfalseを返す

package input

func isTouchEnabled() bool {
	return false
}
