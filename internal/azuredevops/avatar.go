package azuredevops

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"azuredevops-notify/internal/config"

	"github.com/gin-gonic/gin"
)

// maxAvatarBytes caps how much of an upstream response the proxy will copy
// back, so a wrong URL cannot stream something huge into a chat client.
const maxAvatarBytes = 2 << 20 // 2 MiB

// avatarCacheControl lets the chat platforms' image caches hold an avatar
// for a while instead of re-fetching it for every notification.
const avatarCacheControl = "public, max-age=3600"

// AvatarProxy serves Azure DevOps identity avatars on this service's behalf.
//
// The `imageUrl` Azure DevOps puts in a webhook payload points at an
// endpoint that requires authentication (and on Azure DevOps Server is
// usually only reachable from inside the network). Discord and Google Chat
// fetch that URL anonymously from their own servers, so they get a sign-in
// redirect or an HTML page instead of an image and render a blank avatar.
// Pointing the URL back at this service — which already holds the PAT —
// makes the image load.
//
// Rewriting is opt-in: with BaseURL empty the raw Azure URL is used
// unchanged, which is the right behaviour when Azure is publicly readable.
type AvatarProxy struct {
	// BaseURL is this service's own externally reachable base URL, e.g.
	// https://notify.example.com (PUBLIC_BASE_URL). It has to be reachable
	// by the chat platform, not just by Azure DevOps.
	BaseURL string
	// Azure supplies the organization URL the proxy is allowed to fetch
	// from, and the PAT used to authenticate against it.
	Azure config.AzureConfig
	// HTTPClient is optional; a client with a sane timeout is used when nil.
	HTTPClient *http.Client
}

// URLFor rewrites an Azure DevOps identity image URL into one served by
// this service. It returns imageURL unchanged when rewriting is disabled,
// when there is no image, or when the URL is not on the configured Azure
// DevOps host — the proxy only ever fetches from that one host, so a URL it
// could not serve is better left as-is than pointed at a route that would
// refuse it.
func (p *AvatarProxy) URLFor(imageURL string) string {
	if p == nil || p.BaseURL == "" || imageURL == "" || !p.isAzureURL(imageURL) {
		return imageURL
	}

	ref := base64.RawURLEncoding.EncodeToString([]byte(imageURL))

	return strings.TrimSuffix(p.BaseURL, "/") + strings.Replace(RouteAvatar, ":ref", ref, 1)
}

// Avatar fetches the referenced Azure DevOps avatar with the PAT and
// streams it back to the caller.
func (p *AvatarProxy) Avatar(c *gin.Context) {
	ref, err := base64.RawURLEncoding.DecodeString(c.Param("ref"))
	if err != nil {
		respondError(c, fmt.Errorf("avatar: invalid reference: %w", err))
		return
	}

	imageURL := string(ref)
	if !p.isAzureURL(imageURL) {
		respondError(c, fmt.Errorf("avatar: refusing to fetch %q: not on the configured Azure DevOps host", imageURL))
		return
	}

	req, err := http.NewRequestWithContext(c.Request.Context(), http.MethodGet, imageURL, nil)
	if err != nil {
		respondError(c, fmt.Errorf("avatar: build request: %w", err))
		return
	}

	authToken := base64.StdEncoding.EncodeToString([]byte(":" + p.Azure.PATToken))
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", authToken))

	resp, err := p.client().Do(req)
	if err != nil {
		respondErrorStatus(c, http.StatusBadGateway, fmt.Errorf("avatar: fetch from azure: %w", err))
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		respondErrorStatus(c, http.StatusBadGateway, fmt.Errorf("avatar: azure returned status %d: %s", resp.StatusCode, strings.TrimSpace(string(body))))
		return
	}

	// An unauthenticated or unauthorized request is answered with a sign-in
	// page rather than an error status, so the content type is what tells
	// a working PAT from a broken one.
	contentType := resp.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "image/") {
		respondErrorStatus(c, http.StatusBadGateway, fmt.Errorf("avatar: azure returned %q, not an image (check AZURE_PAT_TOKEN)", contentType))
		return
	}

	// Reject an oversized image rather than truncating it into a corrupt
	// one; the LimitReader then only guards a response that declares no
	// length at all.
	if resp.ContentLength > maxAvatarBytes {
		respondErrorStatus(c, http.StatusBadGateway, fmt.Errorf("avatar: azure returned %d bytes, over the %d byte limit", resp.ContentLength, maxAvatarBytes))
		return
	}

	c.Header("Cache-Control", avatarCacheControl)
	c.DataFromReader(http.StatusOK, resp.ContentLength, contentType, io.LimitReader(resp.Body, maxAvatarBytes), nil)
}

// isAzureURL reports whether rawURL lives on the same scheme and host as
// the configured Azure DevOps organization. It is what keeps the route from
// being an open proxy that would attach the PAT to any URL a caller names.
func (p *AvatarProxy) isAzureURL(rawURL string) bool {
	if p == nil || p.Azure.Organization == "" {
		return false
	}

	organization, err := url.Parse(p.Azure.Organization)
	if err != nil {
		return false
	}

	target, err := url.Parse(rawURL)
	if err != nil {
		return false
	}

	return target.Scheme == organization.Scheme &&
		strings.EqualFold(target.Host, organization.Host) &&
		target.Host != ""
}

func (p *AvatarProxy) client() *http.Client {
	if p.HTTPClient != nil {
		return p.HTTPClient
	}

	return &http.Client{Timeout: 10 * time.Second}
}
