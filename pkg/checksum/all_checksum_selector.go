package checksum

type AllChecksumSelector struct{}

var _ ChecksumSelector = AllChecksumSelector{}

func (AllChecksumSelector) SelectChecksums(in ChecksumList) ChecksumList {
	return in
}
