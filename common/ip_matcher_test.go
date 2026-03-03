package common

import (
	"net"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// ---------------------------------------------------------------------------
// TestValidateCIDRList
// ---------------------------------------------------------------------------

func TestValidateCIDRList_EmptyList(t *testing.T) {
	err := ValidateCIDRList([]string{})
	assert.NoError(t, err)
}

func TestValidateCIDRList_ValidIPv4(t *testing.T) {
	err := ValidateCIDRList([]string{"192.168.1.1"})
	assert.NoError(t, err)
}

func TestValidateCIDRList_ValidIPv6(t *testing.T) {
	err := ValidateCIDRList([]string{"::1"})
	assert.NoError(t, err)
}

func TestValidateCIDRList_ValidCIDRv4(t *testing.T) {
	err := ValidateCIDRList([]string{"10.0.0.0/8", "172.16.0.0/12"})
	assert.NoError(t, err)
}

func TestValidateCIDRList_ValidCIDRv6(t *testing.T) {
	err := ValidateCIDRList([]string{"2001:db8::/32"})
	assert.NoError(t, err)
}

func TestValidateCIDRList_InvalidIPSegment(t *testing.T) {
	// IP segment out of range
	err := ValidateCIDRList([]string{"256.0.0.1"})
	assert.Error(t, err)
}

func TestValidateCIDRList_InvalidPrefixRange(t *testing.T) {
	// Prefix length out of range for IPv4
	err := ValidateCIDRList([]string{"10.0.0.0/33"})
	assert.Error(t, err)
}

func TestValidateCIDRList_NonIPString(t *testing.T) {
	err := ValidateCIDRList([]string{"not-an-ip"})
	assert.Error(t, err)
}

func TestValidateCIDRList_EmptyString(t *testing.T) {
	err := ValidateCIDRList([]string{""})
	assert.Error(t, err)
}

func TestValidateCIDRList_MixedValidInvalid_ReturnsFirstInvalid(t *testing.T) {
	err := ValidateCIDRList([]string{"10.0.0.0/8", "bad-entry", "192.168.0.0/16"})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "bad-entry")
}

// ---------------------------------------------------------------------------
// TestParseCIDRList
// ---------------------------------------------------------------------------

func TestParseCIDRList_EmptyList(t *testing.T) {
	result, err := ParseCIDRList([]string{})
	assert.NoError(t, err)
	assert.Empty(t, result)
}

func TestParseCIDRList_SingleIPv4AutoComplete(t *testing.T) {
	result, err := ParseCIDRList([]string{"192.168.1.1"})
	assert.NoError(t, err)
	assert.Len(t, result, 1)
	// The network mask must be /32
	ones, bits := result[0].Mask.Size()
	assert.Equal(t, 32, ones)
	assert.Equal(t, 32, bits)
}

func TestParseCIDRList_SingleIPv6AutoComplete(t *testing.T) {
	result, err := ParseCIDRList([]string{"::1"})
	assert.NoError(t, err)
	assert.Len(t, result, 1)
	ones, bits := result[0].Mask.Size()
	assert.Equal(t, 128, ones)
	assert.Equal(t, 128, bits)
}

func TestParseCIDRList_ValidCIDR(t *testing.T) {
	result, err := ParseCIDRList([]string{"10.0.0.0/8"})
	assert.NoError(t, err)
	assert.Len(t, result, 1)
	ones, _ := result[0].Mask.Size()
	assert.Equal(t, 8, ones)
}

func TestParseCIDRList_InvalidEntry_ReturnsError(t *testing.T) {
	result, err := ParseCIDRList([]string{"10.0.0.0/8", "not-valid"})
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestParseCIDRList_IndependentSlices(t *testing.T) {
	// Two calls must return independent slices; modifying one must not affect the other.
	slice1, err := ParseCIDRList([]string{"10.0.0.0/8"})
	assert.NoError(t, err)
	slice2, err := ParseCIDRList([]string{"10.0.0.0/8"})
	assert.NoError(t, err)

	// Overwrite the first element of slice1.
	slice1[0] = nil

	assert.NotNil(t, slice2[0], "modifying slice1 must not affect slice2")
}

func TestParseCIDRList_ConcurrentSafety(t *testing.T) {
	// 100 goroutines calling ParseCIDRList concurrently — must pass -race.
	const goroutines = 100
	var wg sync.WaitGroup
	wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			_, _ = ParseCIDRList([]string{"10.0.0.0/8", "192.168.1.1"})
		}()
	}
	wg.Wait()
}

// ---------------------------------------------------------------------------
// TestIPMatchesCIDRList
// ---------------------------------------------------------------------------

func mustParseCIDRList(t *testing.T, ips []string) []*net.IPNet {
	t.Helper()
	result, err := ParseCIDRList(ips)
	assert.NoError(t, err)
	return result
}

func TestIPMatchesCIDRList_Hit(t *testing.T) {
	cidrs := mustParseCIDRList(t, []string{"10.0.0.0/8"})
	assert.True(t, IPMatchesCIDRList("10.1.2.3", cidrs))
}

func TestIPMatchesCIDRList_Miss(t *testing.T) {
	cidrs := mustParseCIDRList(t, []string{"10.0.0.0/8"})
	assert.False(t, IPMatchesCIDRList("192.168.1.1", cidrs))
}

func TestIPMatchesCIDRList_EmptyCIDRs(t *testing.T) {
	assert.False(t, IPMatchesCIDRList("10.1.2.3", nil))
}

func TestIPMatchesCIDRList_NetworkBoundaryFirst(t *testing.T) {
	cidrs := mustParseCIDRList(t, []string{"192.168.0.0/24"})
	assert.True(t, IPMatchesCIDRList("192.168.0.0", cidrs))
}

func TestIPMatchesCIDRList_NetworkBoundaryLast(t *testing.T) {
	cidrs := mustParseCIDRList(t, []string{"192.168.0.0/24"})
	assert.True(t, IPMatchesCIDRList("192.168.0.255", cidrs))
}

func TestIPMatchesCIDRList_IPv6Hit(t *testing.T) {
	cidrs := mustParseCIDRList(t, []string{"2001:db8::/32"})
	assert.True(t, IPMatchesCIDRList("2001:db8::1", cidrs))
}

func TestIPMatchesCIDRList_InvalidIP(t *testing.T) {
	cidrs := mustParseCIDRList(t, []string{"10.0.0.0/8"})
	assert.False(t, IPMatchesCIDRList("not-an-ip", cidrs))
}

// ---------------------------------------------------------------------------
// TestGetClientIP
// ---------------------------------------------------------------------------

func newGinContext(method, path string) (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(method, path, nil)
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	return c, w
}

func TestGetClientIP_NoTrustedProxy_IgnoresXFF(t *testing.T) {
	// Reset trusted proxies.
	InitTrustedProxies("")

	c, _ := newGinContext("GET", "/")
	c.Request.RemoteAddr = "1.2.3.4:5678"
	c.Request.Header.Set("X-Forwarded-For", "9.9.9.9")

	ip := GetClientIP(c)
	assert.Equal(t, "1.2.3.4", ip)
}

func TestGetClientIP_TrustedProxy_UsesXFFLeftmost(t *testing.T) {
	InitTrustedProxies("127.0.0.1")

	c, _ := newGinContext("GET", "/")
	c.Request.RemoteAddr = "127.0.0.1:9000"
	c.Request.Header.Set("X-Forwarded-For", "5.6.7.8, 10.0.0.1")

	ip := GetClientIP(c)
	assert.Equal(t, "5.6.7.8", ip)
}

func TestGetClientIP_TrustedProxy_XFFEmpty_FallsBackToRemoteAddr(t *testing.T) {
	InitTrustedProxies("127.0.0.1")

	c, _ := newGinContext("GET", "/")
	c.Request.RemoteAddr = "127.0.0.1:9000"
	// No XFF header.

	ip := GetClientIP(c)
	assert.Equal(t, "127.0.0.1", ip)
}

func TestGetClientIP_NonTrustedProxy_IgnoresXFF(t *testing.T) {
	InitTrustedProxies("10.0.0.0/8")

	c, _ := newGinContext("GET", "/")
	c.Request.RemoteAddr = "1.2.3.4:5678" // not in trusted range
	c.Request.Header.Set("X-Forwarded-For", "9.9.9.9")

	ip := GetClientIP(c)
	assert.Equal(t, "1.2.3.4", ip)
}

// ---------------------------------------------------------------------------
// TestInitTrustedProxies
// ---------------------------------------------------------------------------

func TestInitTrustedProxies_EmptyString_SetsNil(t *testing.T) {
	InitTrustedProxies("")
	assert.Nil(t, _trustedProxyCidrs)
}

func TestInitTrustedProxies_SingleIP(t *testing.T) {
	InitTrustedProxies("127.0.0.1")
	assert.Len(t, _trustedProxyCidrs, 1)
}

func TestInitTrustedProxies_SingleCIDR(t *testing.T) {
	InitTrustedProxies("10.0.0.0/8")
	assert.Len(t, _trustedProxyCidrs, 1)
}

func TestInitTrustedProxies_MultipleCommaSeparated(t *testing.T) {
	InitTrustedProxies("127.0.0.1, 10.0.0.0/8, ::1")
	assert.Len(t, _trustedProxyCidrs, 3)
}

func TestInitTrustedProxies_InvalidEntrySkipped(t *testing.T) {
	InitTrustedProxies("127.0.0.1, bad-entry, 10.0.0.0/8")
	// Only two valid entries; bad-entry must be skipped.
	assert.Len(t, _trustedProxyCidrs, 2)
}
