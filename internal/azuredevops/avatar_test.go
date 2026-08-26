package azuredevops

import (
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"azuredevops-notify/internal/config"
	"azuredevops-notify/internal/notify"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

const azureOrganization = "https://devops.example.com/SampleCollection"

// a 1x1 transparent GIF, small enough to inline as the stub's image body
var fakeImage = []byte("GIF89a\x01\x00\x01\x00\x80\x00\x00\x00\x00\x00\x00\x00\x00!\xf9\x04\x01\x00\x00\x00\x00,\x00\x00\x00\x00\x01\x00\x01\x00\x00\x02\x02D\x01\x00;")

func TestAvatarProxyURLFor(t *testing.T) {
	proxy := &AvatarProxy{
		BaseURL: "https://notify.example.com",
		Azure:   config.AzureConfig{Organization: azureOrganization},
	}

	azureImage := azureOrganization + "/_apis/GraphProfile/MemberAvatars/aad.abc"

	tests := []struct {
		name     string
		proxy    *AvatarProxy
		imageURL string
		want     string
	}{
		{
			name:     "rewrites an azure image url",
			proxy:    proxy,
			imageURL: azureImage,
			want:     "https://notify.example.com/avatar/" + base64.RawURLEncoding.EncodeToString([]byte(azureImage)),
		},
		{
			name:     "rewrites a url elsewhere on the azure host",
			proxy:    proxy,
			imageURL: "https://devops.example.com/_api/_common/identityImage?id=a4dd11b6",
			want: "https://notify.example.com/avatar/" +
				base64.RawURLEncoding.EncodeToString([]byte("https://devops.example.com/_api/_common/identityImage?id=a4dd11b6")),
		},
		{
			name: "trims a trailing slash on the base url",
			proxy: &AvatarProxy{
				BaseURL: "https://notify.example.com/",
				Azure:   config.AzureConfig{Organization: azureOrganization},
			},
			imageURL: azureImage,
			want:     "https://notify.example.com/avatar/" + base64.RawURLEncoding.EncodeToString([]byte(azureImage)),
		},
		{
			name:     "nil proxy leaves the url alone",
			proxy:    nil,
			imageURL: azureImage,
			want:     azureImage,
		},
		{
			name:     "no public base url leaves the url alone",
			proxy:    &AvatarProxy{Azure: config.AzureConfig{Organization: azureOrganization}},
			imageURL: azureImage,
			want:     azureImage,
		},
		{
			name:     "a url on another host is left alone",
			proxy:    proxy,
			imageURL: "https://evil.example.com/avatar.png",
			want:     "https://evil.example.com/avatar.png",
		},
		{
			name:     "an empty image url stays empty",
			proxy:    proxy,
			imageURL: "",
			want:     "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.proxy.URLFor(tt.imageURL))
		})
	}
}

func TestAvatarProxyMessagesCarryTheProxiedIcon(t *testing.T) {
	proxy := &AvatarProxy{
		BaseURL: "https://notify.example.com",
		Azure:   config.AzureConfig{Organization: azureOrganization},
	}

	azureImage := azureOrganization + "/_apis/GraphProfile/MemberAvatars/aad.abc"
	want := proxy.URLFor(azureImage)

	req := &AzureRequest{}
	req.Resource.CreatedBy.ImageUrl = azureImage
	req.Resource.RequestedFor.ImageUrl = azureImage
	req.Resource.Deployment.RequestedFor.ImageUrl = azureImage

	assert.Equal(t, want, req.toPRMessage("t", notify.LevelPending, proxy).Author.IconURL)
	assert.Equal(t, want, req.toPipelineMessage("t", notify.LevelSuccess, AzureRepository{}, proxy).Author.IconURL)
	assert.Equal(t, want, req.toReleaseMessage("t", notify.LevelSuccess, proxy).Author.IconURL)

	// With no proxy configured the raw Azure URL is used, as before.
	assert.Equal(t, azureImage, req.toPRMessage("t", notify.LevelPending, nil).Author.IconURL)
}

// serveAvatar runs one GET against the avatar route, pointed at an upstream
// stub standing in for Azure DevOps.
func serveAvatar(t *testing.T, upstream *httptest.Server, imageURL string) *httptest.ResponseRecorder {
	t.Helper()

	gin.SetMode(gin.TestMode)

	proxy := &AvatarProxy{
		BaseURL: "https://notify.example.com",
		Azure: config.AzureConfig{
			Organization: upstream.URL,
			PATToken:     "fake-token",
		},
		HTTPClient: upstream.Client(),
	}

	r := gin.New()
	r.GET(RouteAvatar, proxy.Avatar)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, proxy.URLFor(imageURL), nil))

	return w
}

func TestAvatarProxyServesTheImage(t *testing.T) {
	var gotAuth, gotPath string

	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		gotPath = r.URL.Path
		w.Header().Set("Content-Type", "image/gif")
		_, _ = w.Write(fakeImage)
	}))
	defer upstream.Close()

	w := serveAvatar(t, upstream, upstream.URL+"/_apis/GraphProfile/MemberAvatars/aad.abc")

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "image/gif", w.Header().Get("Content-Type"))
	assert.Equal(t, avatarCacheControl, w.Header().Get("Cache-Control"))
	assert.Equal(t, fakeImage, w.Body.Bytes())
	assert.Equal(t, "/_apis/GraphProfile/MemberAvatars/aad.abc", gotPath)
	// The PAT is what the chat platforms cannot send themselves.
	assert.Equal(t, "Basic "+base64.StdEncoding.EncodeToString([]byte(":fake-token")), gotAuth)
}

func TestAvatarProxyRejectsAnotherHost(t *testing.T) {
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("upstream should not be called for a foreign host")
	}))
	defer upstream.Close()

	gin.SetMode(gin.TestMode)

	proxy := &AvatarProxy{
		BaseURL: "https://notify.example.com",
		Azure:   config.AzureConfig{Organization: upstream.URL, PATToken: "fake-token"},
	}

	r := gin.New()
	r.GET(RouteAvatar, proxy.Avatar)

	ref := base64.RawURLEncoding.EncodeToString([]byte("https://evil.example.com/steal"))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/avatar/"+ref, nil))

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAvatarProxyRejectsAnUndecodableReference(t *testing.T) {
	gin.SetMode(gin.TestMode)

	proxy := &AvatarProxy{
		BaseURL: "https://notify.example.com",
		Azure:   config.AzureConfig{Organization: azureOrganization},
	}

	r := gin.New()
	r.GET(RouteAvatar, proxy.Avatar)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/avatar/not-base64!!", nil))

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAvatarProxyReportsUpstreamFailures(t *testing.T) {
	// respondError marshals the error as gin.H{"err": err}, which renders as
	// an empty object, so the status code is what the test can assert on;
	// the detail goes to stdout like every other error in this package.
	tests := []struct {
		name    string
		handler http.HandlerFunc
	}{
		{
			name: "non-200 from azure",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotFound)
				_, _ = w.Write([]byte("no such identity"))
			},
		},
		{
			name: "an image bigger than the limit",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "image/gif")
				w.Header().Set("Content-Length", strconv.Itoa(maxAvatarBytes+1))
				_, _ = w.Write(make([]byte, maxAvatarBytes+1))
			},
		},
		{
			// The symptom this whole proxy exists for: an unauthenticated
			// fetch is answered with a sign-in page, not an error status.
			name: "a sign-in page instead of an image",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "text/html; charset=utf-8")
				_, _ = w.Write([]byte("<html>sign in</html>"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			upstream := httptest.NewServer(tt.handler)
			defer upstream.Close()

			w := serveAvatar(t, upstream, upstream.URL+"/_apis/GraphProfile/MemberAvatars/aad.abc")

			assert.Equal(t, http.StatusBadGateway, w.Code)
		})
	}
}
