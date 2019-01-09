package dktest

import (
	"testing"
	"time"
)

import (
	"github.com/stretchr/testify/assert"
)

func TestOptionsInit(t *testing.T) {
	timeout := 9 * time.Second

	testCases := []struct {
		name     string
		opts     Options
		expected Options
	}{
		{name: "default timeouts used", opts: Options{},
			expected: Options{
				PullTimeout:    DefaultPullTimeout,
				Timeout:        DefaultTimeout,
				ReadyTimeout:   DefaultReadyTimeout,
				CleanupTimeout: DefaultCleanupTimeout,
			},
		},
		{name: "default timeouts not used",
			opts: Options{
				PullTimeout:    timeout,
				Timeout:        timeout,
				ReadyTimeout:   timeout,
				CleanupTimeout: timeout,
			},
			expected: Options{
				PullTimeout:    timeout,
				Timeout:        timeout,
				ReadyTimeout:   timeout,
				CleanupTimeout: timeout,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.opts.init()
			assert.Equal(t, tc.expected, tc.opts, "Expected post-init Options to match expected")
		})
	}
}

func TestOptionsEnv(t *testing.T) {
	testCases := []struct {
		name        string
		env         map[string]string
		expectedEnv []string
	}{
		{name: "nil", env: nil, expectedEnv: nil},
		{name: "empty", env: nil, expectedEnv: nil},
		{name: "1 var", env: map[string]string{"foo": "bar"}, expectedEnv: []string{"foo=bar"}},
		{name: "1 var - empty value", env: map[string]string{"foo": ""}, expectedEnv: []string{"foo="}},
		{name: "1 var - empty key", env: map[string]string{"": "bar"}, expectedEnv: []string{"=bar"}},
		{name: "3 vars", env: map[string]string{"foo": "bar", "hello": "world", "dead": "beef"},
			expectedEnv: []string{"foo=bar", "hello=world", "dead=beef"}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			opts := Options{Env: tc.env}
			assert.ElementsMatch(t, tc.expectedEnv, opts.env(), "Options environment to match expected")
		})
	}
}
