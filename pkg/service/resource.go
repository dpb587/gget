package service

type ResourceType string

// ArchiveResourceType is a tar/zip export of the repository from the ref.
const ArchiveResourceType ResourceType = "archive"

// AssetResourceType is a user-provided file associated with the ref.
const AssetResourceType ResourceType = "asset"

// BlobResourceType is a blob of the repository at the ref.
const BlobResourceType ResourceType = "blob"

type ResourceName string
