# xc-screenshot
Save XCUITest screenshots (or other XCAttachments) automatically.

XCode's TestPlan and UITests go a long way to mitigate the tedium of taking app screenshots, but more can be done. Saving and organizing the screenshot files still requires waiting for the tests to finish and then 1. clicking and saving them manually in XCode from the test reports. or 2. dealing with Xcode's shifting and unweildy DerivedData folders.

A simpler solution that works for me, is to setup a simple local http file server and then have the test code upload the screenshot data to a local file. 

#### 1. Build and run the server locally

The server is written in Go. To build it you need the Go compiler, then just `go build ./http_upload.go`. This will create the `http_upload` binary. Alternatively download the pre-built binary for macOS under releases.

Run the binary from the directory you want the screenshots saved. 

NOTE: This is a quick and simple solution for a specific job. 
The server is as simple as it gets. It isn't daemonized, so stays attached to the terminal and logs there. Also, it is plain http and does no request validation, so it is wholly unsuitable to be used outside of localhost or to handle untrusted requests.

#### 2. Take screenshots while testing

Just add the `Screenshot.swift` file to the UITest target. It extends `XCTestCase` class with the function `takeScreenshot(name: String)` which when called inside a test will take a screenshot of the current app state and upload it wherever `./upload_binary` is running.

If it's not there, you may need to add `Allows Local Networking` to `App Transport Security Settings` in the UITest build target's `Info.plist` (not the apps target).





