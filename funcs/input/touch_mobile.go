// +build android ios darwin,arm darwin,arm64
// +build !js

// mobileビルドすると、isTouchPrimaryInputがtrueを返す

package input

func isTouchEnabled() bool {
	return true
}
