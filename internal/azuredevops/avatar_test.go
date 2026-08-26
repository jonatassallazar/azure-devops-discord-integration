package azuredevops

import (
	"testing"

	"azuredevops-notify/internal/notify"

	"github.com/stretchr/testify/assert"
)

func TestAvatarURL(t *testing.T) {
	// The digest Gravatar expects for david.martinez@example.com, spelled
	// out rather than recomputed, so a change of hash or of normalization
	// fails the test instead of following it.
	const davidAvatar = "https://www.gravatar.com/avatar/47746b155b64f3ae4cffeef88a400b2d?s=80&d=identicon"

	tests := []struct {
		name       string
		uniqueName string
		want       string
	}{
		{
			name:       "an email address maps to its gravatar",
			uniqueName: "david.martinez@example.com",
			want:       davidAvatar,
		},
		{
			name:       "case and surrounding space are normalized away",
			uniqueName: "  David.Martinez@Example.COM  ",
			want:       davidAvatar,
		},
		{
			// An on-premises Active Directory identity, which has no
			// address to key on: better no icon than a wrong one.
			name:       "a domain account has no gravatar",
			uniqueName: `CONTOSO\david.martinez`,
			want:       "",
		},
		{
			name:       "an empty unique name has no gravatar",
			uniqueName: "",
			want:       "",
		},
		{
			name:       "a display name is not an address",
			uniqueName: "David Martinez",
			want:       "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, avatarURL(tt.uniqueName))
		})
	}
}

func TestMessagesCarryTheGravatarIcon(t *testing.T) {
	const want = "https://www.gravatar.com/avatar/47746b155b64f3ae4cffeef88a400b2d?s=80&d=identicon"

	req := &AzureRequest{}
	req.Resource.CreatedBy.UniqueName = "david.martinez@example.com"
	req.Resource.RequestedFor.UniqueName = "david.martinez@example.com"
	req.Resource.Deployment.RequestedFor.UniqueName = "david.martinez@example.com"

	assert.Equal(t, want, req.toPRMessage("t", notify.LevelPending).Author.IconURL)
	assert.Equal(t, want, req.toPipelineMessage("t", notify.LevelSuccess, AzureRepository{}).Author.IconURL)
	assert.Equal(t, want, req.toReleaseMessage("t", notify.LevelSuccess).Author.IconURL)
}

// An identity with no usable address leaves the icon off the message
// entirely rather than pointing the chat platforms at an image they cannot
// load.
func TestMessagesDropTheIconWithoutAnAddress(t *testing.T) {
	req := &AzureRequest{}
	req.Resource.CreatedBy.UniqueName = `CONTOSO\david.martinez`
	req.Resource.CreatedBy.ImageUrl = "https://devops.example.com/_apis/GraphProfile/MemberAvatars/aad.abc"

	assert.Equal(t, "", req.toPRMessage("t", notify.LevelPending).Author.IconURL)
}
