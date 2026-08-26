package azuredevops

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/mail"
	"strings"
)

// gravatarSize is the pixel size requested from Gravatar. Both Discord and
// Google Chat render an author icon small, so Gravatar's own default of 80
// is plenty.
const gravatarSize = 80

// gravatarDefault is what Gravatar serves for an address with no account:
// "identicon" is a geometric pattern derived from the same hash, so a user
// without a Gravatar still gets a stable, distinct image rather than a
// blank circle. Gravatar's other options include "mp" (a generic
// silhouette), "retro" and "404".
const gravatarDefault = "identicon"

// avatarURL maps an Azure DevOps identity to an author icon the chat
// platforms can actually fetch.
//
// The imageUrl Azure DevOps puts in a webhook payload points at an endpoint
// that requires authentication, and on Azure DevOps Server that host is
// usually internal-only. Discord and Google Chat fetch an author icon
// anonymously from their own servers, so they get a sign-in page instead of
// an image and render a blank avatar. Gravatar serves the same person's
// picture from a public URL keyed by their email address, which keeps image
// delivery out of this service altogether — nothing here fetches, proxies
// or caches an image.
//
// uniqueName is the identity's email address (the UPN, for Azure AD
// accounts). When it is not an address — an on-premises DOMAIN\user, a
// service account, an empty field — there is nothing to key on, so the
// result is empty and the notification carries no author icon at all.
func avatarURL(uniqueName string) string {
	address, err := mail.ParseAddress(strings.TrimSpace(uniqueName))
	if err != nil {
		return ""
	}

	// Gravatar keys on a hash of the trimmed, lowercased address. MD5 is
	// the form their API has always accepted; it is an identity hash
	// defined by that API, not a security primitive.
	sum := md5.Sum([]byte(strings.ToLower(address.Address)))

	return fmt.Sprintf(
		"https://www.gravatar.com/avatar/%s?s=%d&d=%s",
		hex.EncodeToString(sum[:]), gravatarSize, gravatarDefault,
	)
}
