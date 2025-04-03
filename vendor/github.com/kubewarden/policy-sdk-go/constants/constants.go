package constants

const (
	// ProtocolVersion is the version of the protocol used by the
	// Kubewarden waPC host and guest to exchange information.
	ProtocolVersion = "v1"
	// These two media types are manifests media types that the
	// oci-distribution accepts when fetching images manifests. But they
	// are not present in the open containers go lib used in the SDK.
	// Therefore, we need to check for them as well in in the code used to
	// fetch and parse the OCI image manifests.
	ImageManifestListMediaType = "application/vnd.docker.distribution.manifest.list.v2+json"
	ImageManifestMediaType     = "application/vnd.docker.distribution.manifest.v2+json"
)
