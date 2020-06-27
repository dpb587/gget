package checksum

type ChecksumSelector interface {
	SelectChecksums(ChecksumList) ChecksumList
}
