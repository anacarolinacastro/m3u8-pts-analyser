package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/grafov/m3u8"
	"gopkg.in/vansante/go-ffprobe.v2"
)

func getPodID(line string) string {
	podIDRegexp := `.*/pod/([0-9]*)/profile/.*`
	r := regexp.MustCompile(podIDRegexp)
	match := r.FindStringSubmatch(line)
	if len(match) > 1 {
		return match[1]
	}
	return ""
}

func parseStartEndPTS(filePath string) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	data, err := ffprobe.ProbeURL(ctx, filePath)
	if err != nil {
		fmt.Printf("Error getting data: %v\n", err)
		panic(err)
	}

	for _, stream := range data.Streams {
		if stream.CodecType == string(ffprobe.StreamVideo) {
			fmt.Printf("pts: %d durationTS %f", stream.StartPts, data.Format.DurationSeconds)
		}
	}
	fmt.Printf(" startTime: %fs endTime: %fs [%s]\n", data.Format.StartTimeSeconds, data.Format.DurationSeconds+data.Format.StartTimeSeconds, filepath.Base(filePath))

}

func downloadSegment(fileName string, segmentURL *url.URL) {
	resp, err := http.Get(segmentURL.String())
	if err != nil || resp.StatusCode > 399 {
		fmt.Printf("Error downloading %s\n", segmentURL)
		panic(err)
	}

	defer resp.Body.Close()
	out, err := os.Create(fileName)
	if err != nil {
		fmt.Printf("Error creating file %s: %v\n", fileName, err)
		panic(err)
	}
	defer out.Close()
	io.Copy(out, resp.Body)
}

func downloadPlaylist(playlistURL, playlistFilePath string) {
	resp, err := http.Get(playlistURL)
	if err != nil {
		fmt.Printf("Error downloading %s\n", playlistURL)
		panic(err)
	}

	defer resp.Body.Close()
	out, err := os.Create(playlistFilePath)
	if err != nil {
		fmt.Printf("Error creating file %s: %v\n", "playlist.m3u8", err)
		panic(err)
	}
	defer out.Close()
	io.Copy(out, resp.Body)
}

func parseMediaPlaylist(playlistURL *url.URL, tmpDir, filePath string) {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Error opening file %s: %v\n", filePath, err)
		panic(err)
	}

	playlistBufioReader := bufio.NewReader(file)
	playlist, _, _ := m3u8.DecodeFrom(playlistBufioReader, true)
	mediapl := playlist.(*m3u8.MediaPlaylist)

	for _, segment := range mediapl.Segments {
		if segment == nil {
			break
		}

		if segment.Discontinuity {
			fmt.Printf("DISCONTINUITY\n")
		}

		segmentURL, err := url.Parse(segment.URI)
		if err != nil {
			panic(err)
		}

		urlWithoutParams := strings.Split(segmentURL.String(), "?")[0]
		urlPath := filepath.Base(urlWithoutParams)
		podId := getPodID(urlWithoutParams)
		segmentfileName := fmt.Sprintf("%s/%s-%s", tmpDir, podId, urlPath)

		if segmentURL.Scheme == "" {
			segmentURL = playlistURL.ResolveReference(segmentURL)
			// will evaluate only Google segments

			continue
		}

		downloadSegment(segmentfileName, segmentURL)
		parseStartEndPTS(segmentfileName)
	}
}

func main() {
	if len(os.Args) <= 1 {
		fmt.Println("Error: Missing playlist URL")
		return
	}

	playlist := os.Args[1]
	playlistURL, err := url.Parse(playlist)
	if err != nil {
		panic(err)
	}

	tmpDir := fmt.Sprintf("/tmp/%d", time.Now().Unix())
	os.MkdirAll(tmpDir, os.ModePerm)
	fmt.Printf("\nAnalyzing: %s\n\nSaving files at: %s\n\n", playlist, tmpDir)

	playlistFilePath := fmt.Sprintf("%s/playlist.m3u8", tmpDir)
	downloadPlaylist(playlistURL.String(), playlistFilePath)
	parseMediaPlaylist(playlistURL, tmpDir, playlistFilePath)
}
