package writer

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/hoshsadiq/m3ufilter/m3u"
	"io"
	"strconv"
	"strings"
)

func writeM3U(w io.Writer, streams []*m3u.Stream) {
	_, err := w.Write([]byte("#EXTM3U"))
	if err != nil {
		log.Fatalf("unable to write extm3u, err = %v", err)
	}

	for i, stream := range streams {
		extinf := getStreamExtinf(stream)
		_, err := w.Write(extinf)
		if err != nil {
			log.Fatalf("unable to write new streams, i = %d, err = %v, extinf = %v", i, err, extinf)
		}
	}
}

func GetMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

func getStreamExtinf(stream *m3u.Stream) []byte {
	b := &strings.Builder{}
	b.WriteString("\n")
	b.WriteString("#EXTINF:")
	b.WriteString(stream.Duration)

	if stream.ChNo != "" {
		writeKV(b, "tvg-chno", stream.ChNo)
	}

	writeKV(b, "CUID", GetMD5Hash(stream.Uri))
	writeKV(b, "tvg-id", stream.Id)
	writeKV(b, "tvg-name", stream.Name)
	writeKV(b, "group-title", stream.Group)
	writeKV(b, "tvg-logo", stream.Logo)

	if stream.Shift != "" {
		writeKV(b, "tvg-shift", stream.Shift)
	}

	b.WriteRune(',')
	b.WriteString(stream.Name)
	b.WriteString("\n")
	b.WriteString(stream.Uri)

	return []byte(b.String())
}

func writeKV(b *strings.Builder, key string, value string) {
	b.WriteRune(' ')
	b.WriteString(key)
	b.WriteRune('=')
	b.WriteString(strconv.Quote(strings.TrimSpace(value)))
}
