//go:build long

package gositemapfetcher

import (
	"bufio"
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"strconv"
	"testing"
	"time"
)

func TestSitemapFetcher_LargeRandom(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping long test in short mode")
	}
	if os.Getenv("GO_SITEMAP_FETCHER_LONG") == "" {
		t.Skip("set GO_SITEMAP_FETCHER_LONG=1 to run")
	}

	const totalURLs = 10_000_000

	rng := rand.New(rand.NewSource(42))

	server := newTestServer(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/sitemap.xml" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/xml")
		writer := bufio.NewWriterSize(w, 1<<20)
		_, _ = writer.WriteString("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n")
		_, _ = writer.WriteString("<urlset xmlns=\"http://www.sitemaps.org/schemas/sitemap/0.9\">\n")

		base := r.Host
		var numBuf [32]byte
		for i := 0; i < totalURLs; i++ {
			randomID := rng.Uint64()
			writer.WriteString("  <url><loc>http://")
			writer.WriteString(base)
			writer.WriteString("/")
			writer.Write(strconv.AppendUint(numBuf[:0], randomID, 36))
			writer.WriteString("/page-")
			writer.Write(strconv.AppendInt(numBuf[:0], int64(i), 10))
			writer.WriteString("</loc></url>\n")
		}

		_, _ = writer.WriteString("</urlset>")
		_ = writer.Flush()
	}))
	defer server.Close()

	sitemapURL, err := url.Parse(server.URL + "/sitemap.xml")
	if err != nil {
		t.Fatalf("failed to parse sitemap URL: %v", err)
	}

	fetcher := New(Options{})
	var count int
	lastReport := time.Now()
	reportEvery := 1_000_000
	if err := fetcher.Walk(context.Background(), sitemapURL, func(Item) error {
		count++
		if count%reportEvery == 0 {
			reportMem(t, count, time.Since(lastReport))
			lastReport = time.Now()
		}
		return nil
	}); err != nil {
		t.Fatalf("walk failed: %v", err)
	}
	reportMem(t, count, time.Since(lastReport))
	if count != totalURLs {
		t.Fatalf("expected %d URLs, got %d", totalURLs, count)
	}
	fmt.Fprintf(os.Stdout, "parsed %d URLs\n", count)
}

func reportMem(t *testing.T, count int, elapsed time.Duration) {
	t.Helper()
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)
	fmt.Printf("urls=%d elapsed=%s alloc_mb=%d heap_inuse_mb=%d sys_mb=%d\n",
		count,
		elapsed.Truncate(time.Millisecond),
		ms.Alloc/1024/1024,
		ms.HeapInuse/1024/1024,
		ms.Sys/1024/1024,
	)
}
