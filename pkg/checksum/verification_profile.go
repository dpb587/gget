package checksum

type VerificationProfile struct {
	Required   bool
	Acceptable AlgorithmList
	Selector   ChecksumSelector
}
