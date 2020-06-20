package checksum

var strongestHashOrder = []string{"sha512", "sha256", "sha1", "md5"}

func StrongestChecksum(in []Checksum) Checksum {
	for _, strong := range strongestHashOrder {
		for _, checksum := range in {
			if checksum.Algorithm() == strong {
				return checksum
			}
		}
	}

	return nil
}
