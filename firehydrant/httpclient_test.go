package firehydrant_test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/firehydrant/terraform-provider-firehydrant/firehydrant"
	"golang.org/x/time/rate"
)

func TestRateLimitedHTTPDoer(t *testing.T) {
	t.Run("PassWithinLimit", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, "Hello, doer")
		}))
		defer ts.Close()

		doer := firehydrant.NewRateLimitedHTTPDoer()
		doer.WithLimiter(rate.NewLimiter(rate.Every(1*time.Minute), 1))

		req, err := http.NewRequest(http.MethodGet, ts.URL, nil)
		if err != nil {
			t.Fatal(err)
		}
		resp, err := doer.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected status code %d, got %d", http.StatusOK, resp.StatusCode)
		}
	})

	t.Run("RateLimitedOnClient", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, "Hello, doer")
		}))
		defer ts.Close()

		doer := firehydrant.NewRateLimitedHTTPDoer()
		doer.WithLimiter(rate.NewLimiter(rate.Every(1*time.Minute), 1))

		req, err := http.NewRequest(http.MethodGet, ts.URL, nil)
		if err != nil {
			t.Fatal(err)
		}

		// First one should pass
		resp, err := doer.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected status code %d, got %d", http.StatusOK, resp.StatusCode)
		}

		// Second one should fail
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		t.Cleanup(cancel)
		req2 := req.WithContext(ctx)
		if _, err = doer.Do(req2); err == nil {
			t.Fatal("expected error, got nil")
		}
		if !errors.Is(err, firehydrant.ErrClientRateLimitedTimeout) {
			t.Errorf("expected error %v, got %v", firehydrant.ErrClientRateLimitedTimeout, err)
		}
	})

	t.Run("RateLimitedOnServer", func(t *testing.T) {
		calls := 0
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			calls++
			if calls > 1 {
				w.WriteHeader(http.StatusTooManyRequests)
				fmt.Fprintln(w, "Rate limited")
			}
			fmt.Fprintln(w, "Hello, doer")
		}))
		defer ts.Close()

		doer := firehydrant.NewRateLimitedHTTPDoer()
		doer.WithLimiter(rate.NewLimiter(rate.Every(1*time.Minute), 1))
		doer.WithBackoff(1 * time.Second)

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		t.Cleanup(cancel)
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, ts.URL, nil)
		if err != nil {
			t.Fatal(err)
		}

		// This should succeed because internally retry has been handled.
		resp, err := doer.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected status code %d, got %d", http.StatusOK, resp.StatusCode)
		}
	})
}
