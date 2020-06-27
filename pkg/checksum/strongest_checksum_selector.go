package checksum

type StrongestChecksumSelector struct{}

var _ ChecksumSelector = StrongestChecksumSelector{}

func (StrongestChecksumSelector) SelectChecksums(in ChecksumList) ChecksumList {
	for _, algo := range AlgorithmsByStrength {
		for _, cs := range in {
			if cs.Algorithm() == algo {
				return ChecksumList{cs}
			}
		}
	}

	return nil
}
