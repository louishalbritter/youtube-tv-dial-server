# YouTube TV DIAL Server

A reverse-engineered DIAL server for launching YouTube TV in a browser via the YouTube mobile app.

This project allows you to connect your YouTube mobile app to your server, which then automatically launches YouTube TV in an incognito browser tab. Given the lack of official documentation from YouTube on the DIAL protocol, this implementation is somewhat hacky and was developed through reverse engineering. Consequently, the server may lack certain protocol features and robust error handling. However, it has been tested successfully on `linux x64` and `linux arm` platforms in conjunction with the Android YouTube app.

## Prerequisites

- **Golang**: Ensure Go is installed on your system.
- **Browser**: The server is configured to launch the `brave-beta` browser by default, chosen for its ad-blocking and privacy features. However, any browser with configurable command-line arguments should work. You can easily switch browsers by modifying the `Service.start()` method in `service.go`. I use the `beta` version for avoiding messing up the windows and processes of my normal browser.

## Usage

1. Build the project:

```bash
go build .
```

2. Run the server:

```bash
./youtube-tv-dial-server
```
