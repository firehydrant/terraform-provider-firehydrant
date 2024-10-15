package firehydrant

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"golang.org/x/time/rate"
)

var (
	ErrClientRateLimitedTimeout = errors.New("client rate limit timeout")
)

// RateLimitedHTTPDoer is a rate-limited HTTP client which wraps http.Client and rate.Limiter together.
// It exposes Doer interface accepted by github.com/dghubble/sling.
type RateLimitedHTTPDoer struct {
	client  *http.Client
	limiter *rate.Limiter
	mutex   sync.Mutex
	backoff time.Duration
}

const (
	RateLimitSeconds  = 5 * time.Second
	RateLimitRequests = 10
)

func DefaultHTTPDoerRateLimit() *rate.Limiter {
	return rate.NewLimiter(rate.Every(RateLimitSeconds), RateLimitRequests)
}

func NewRateLimitedHTTPDoer() *RateLimitedHTTPDoer {
	return &RateLimitedHTTPDoer{
		client:  &http.Client{},
		limiter: DefaultHTTPDoerRateLimit(),
		mutex:   sync.Mutex{},
		backoff: 5 * time.Second,
	}
}

func (c *RateLimitedHTTPDoer) WithLimiter(l *rate.Limiter) *RateLimitedHTTPDoer {
	c.limiter = l
	return c
}

func (c *RateLimitedHTTPDoer) WithBackoff(d time.Duration) *RateLimitedHTTPDoer {
	c.backoff = d
	return c
}

func (c *RateLimitedHTTPDoer) Do(req *http.Request) (*http.Response, error) {
	start := time.Now()

	// In the event of a rate limit, make sure we don't end up making requests that we know will fail.
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if err := c.limiter.Wait(req.Context()); err != nil {
		return nil, fmt.Errorf("%w after %s: %s", ErrClientRateLimitedTimeout, time.Since(start).String(), err)
	}

	var resp *http.Response
	var err error
	// At this point, only retry on 429. Any other scenario should pass.
	for i := range 5 {
		resp, err = c.client.Do(req)
		if err != nil {
			return resp, err
		}
		// If response is nil, there is no way to actually inspect the status code.
		if resp == nil {
			return resp, errors.New("response is nil")
		}
		if resp.StatusCode == http.StatusTooManyRequests && (i+1) < 5 {
			resp.Body.Close()
			tflog.Warn(req.Context(), "rate limited, queueing for retry", map[string]any{
				"status_code": resp.StatusCode,
				"duration":    time.Since(start).String(),
				"attempt":     i + 1,
				"backoff":     c.backoff.String(),
			})

			backoff := c.backoff
			// If `Retry-After` header is present, _and_ it's less than 30s, use that as the backoff.
			// 30s is set as ceiling to prevent any unknown timeouts from Terraform.
			if resp.Header.Get("Retry-After") != "" {
				if n, err := strconv.Atoi(resp.Header.Get("Retry-After")); err == nil && n < 30 {
					backoff = time.Duration(n) * time.Second
				}
			}
			time.Sleep(backoff)
			continue
		}
		break
	}

	return resp, err
}
