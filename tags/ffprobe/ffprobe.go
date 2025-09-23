package ffprobe

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"go.senan.xyz/gonic/tags/tagcommon"
)

type FFProbe struct{}

func (FFProbe) CanRead(absPath string) bool {
	switch ext := strings.ToLower(filepath.Ext(absPath)); ext {
	case ".mp3", ".flac", ".aac", ".m4a", ".m4b", ".ogg", ".opus", ".wma", ".wav", ".wv":
		return true
	}
	return false
}

func (FFProbe) Read(absPath string) (tagcommon.Info, error) {
	out, err := exec.Command("ffprobe", "-hide_banner", "-v", "0", "-i", absPath, "-show_entries", "format", "-of", "json").Output()
	if err != nil {
		return nil, fmt.Errorf("output: %w", err)
	}

	var format map[string]any
	err = json.Unmarshal(out, &format)
	if err != nil {
		return nil, fmt.Errorf("json: %w", err)
	}
	format = format["format"].(map[string]any)
	tags, ok := format["tags"]
	if !ok {
		return nil, fmt.Errorf("no tags found")
	}

	return &info{format, tags.(map[string]any)}, nil
}

type info struct {
	ffprobe map[string]any
	tags    map[string]any
}

func queryFirstFoundTags(tags map[string]any, queries []string) string {
	for _, query := range queries {
		res, ok := tags[query]
		if ok {
			return res.(string)
		}
	}
	return ""
}

// Album implements tagcommon.Info.
func (i *info) Album() string {
	return queryFirstFoundTags(i.tags, []string{"album", "ALBUM"})
}

// AlbumArtist implements tagcommon.Info.
func (i *info) AlbumArtist() string {
	return queryFirstFoundTags(i.tags, []string{"album_artist"})
}

// AlbumArtists implements tagcommon.Info.
func (i *info) AlbumArtists() []string {
	return strings.Split(
		queryFirstFoundTags(i.tags, []string{"ALBUMARTISTS"}), ";")
}

// AlbumBrainzID implements tagcommon.Info.
func (i *info) AlbumBrainzID() string {
	return queryFirstFoundTags(i.tags, []string{"MUSICBRAINZ_ALBUMID", "MUSICBRAINZ ALBUM ID"})
}

// Artist implements tagcommon.Info.
func (i *info) Artist() string {
	return queryFirstFoundTags(i.tags, []string{"ARTIST", "artist"})
}

// Artists implements tagcommon.Info.
func (i *info) Artists() []string {
	return strings.Split(
		queryFirstFoundTags(i.tags, []string{"ARTISTS"}), ";")
}

// Bitrate implements tagcommon.Info.
func (i *info) Bitrate() int {
	ffbr, _ := i.ffprobe["bit_rate"].(string)
	br, _ := strconv.Atoi(ffbr)
	return br
}

// BrainzID implements tagcommon.Info.
func (i *info) BrainzID() string {
	return queryFirstFoundTags(i.tags, []string{"MUSICBRAINZ_TRACKID"})
}

// Compilation implements tagcommon.Info.
func (i *info) Compilation() bool {
	return queryFirstFoundTags(i.tags, []string{"COMPILATION"}) == "1"
}

// DiscNumber implements tagcommon.Info.
func (i *info) DiscNumber() int {
	ffd, _ := i.ffprobe["disc"].(string)
	d, _ := strconv.Atoi(ffd)
	return d
}

// Genre implements tagcommon.Info.
func (i *info) Genre() string {
	return queryFirstFoundTags(i.tags, []string{"genre"})
}

// Genres implements tagcommon.Info.
func (i *info) Genres() []string {
	return strings.Split(queryFirstFoundTags(i.tags, []string{"GENRES"}), ";")
}

// Length implements tagcommon.Info.
func (i *info) Length() int {
	duration, _ := i.ffprobe["duration"].(string)
	d, _ := strconv.ParseFloat(duration, 64)
	return int(d)
}

// Lyrics implements tagcommon.Info.
func (i *info) Lyrics() string {
	return "" // todo i have no tracks with lyric tags
}

// ReleaseType implements tagcommon.Info.
func (i *info) ReleaseType() string {
	return queryFirstFoundTags(i.tags, []string{"RELEASETYPE"})
}

// ReplayGainAlbumGain implements tagcommon.Info.
func (i *info) ReplayGainAlbumGain() float32 {
	return 0
}

// ReplayGainAlbumPeak implements tagcommon.Info.
func (i *info) ReplayGainAlbumPeak() float32 {
	return 0
}

// ReplayGainTrackGain implements tagcommon.Info.
func (i *info) ReplayGainTrackGain() float32 {
	return 0
}

// ReplayGainTrackPeak implements tagcommon.Info.
func (i *info) ReplayGainTrackPeak() float32 {
	return 0
}

// Title implements tagcommon.Info.
func (i *info) Title() string {
	return queryFirstFoundTags(i.tags, []string{"title", "TITLE"})
}

// TrackNumber implements tagcommon.Info.
func (i *info) TrackNumber() int {
	track := queryFirstFoundTags(i.tags, []string{"track"})
	t, _ := strconv.Atoi(track)
	return t
}

// Year implements tagcommon.Info.
func (i *info) Year() int {
	date := queryFirstFoundTags(i.tags, []string{"DATE", "date", "TDOR"})
	year, _, _ := strings.Cut(date, "-")
	res, _ := strconv.Atoi(year)
	return res
}
