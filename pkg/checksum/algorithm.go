package checksum

import (
	"strings"
)

type Algorithm string

const MD5 Algorithm = "md5"
const SHA1 Algorithm = "sha1"
const SHA256 Algorithm = "sha256"
const SHA384 Algorithm = "sha384"
const SHA512 Algorithm = "sha512"

type AlgorithmList []Algorithm

func (l AlgorithmList) FilterMin(min Algorithm) AlgorithmList {
	var found bool
	var res AlgorithmList

	for _, a := range l {
		res = append(res, a)

		if a == min {
			found = true

			break
		}
	}

	if !found {
		return AlgorithmList{}
	}

	return res
}

func (l AlgorithmList) Contains(in Algorithm) bool {
	for _, v := range l {
		if v == in {
			return true
		}
	}

	return false
}

func (l AlgorithmList) Intersection(in AlgorithmList) AlgorithmList {
	var res AlgorithmList

	for _, exp := range l {
		for _, des := range in {
			if exp == des {
				res = append(res, exp)

				break
			}
		}
	}

	return res
}

func (l AlgorithmList) Join(sep string) string {
	var res []string

	for _, v := range l {
		res = append(res, string(v))
	}

	return strings.Join(res, sep)
}

var AlgorithmsByStrength = AlgorithmList{SHA512, SHA384, SHA256, SHA1, MD5}
