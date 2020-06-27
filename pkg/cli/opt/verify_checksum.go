package opt

import (
	"fmt"
	"strings"

	"github.com/dpb587/gget/pkg/checksum"
)

type VerifyChecksum []string

func (vR VerifyChecksum) Profile() (checksum.VerificationProfile, error) {
	var standaloneValues []string
	var usedCustomAlgos bool

	res := checksum.VerificationProfile{
		Selector: checksum.StrongestChecksumSelector{},
	}

	v := vR

	for _, data := range v {
		if data == "required" { // modifier
			res.Required = true
		} else if data == "all" { // modifier
			res.Selector = checksum.AllChecksumSelector{}
		} else if data == "auto" { // standalone
			// explicit default
			standaloneValues = append(standaloneValues, data)
		} else if data == "none" {
			res.Acceptable = checksum.AlgorithmList{}
			standaloneValues = append(standaloneValues, data)
			usedCustomAlgos = true
		} else if strings.HasSuffix(data, "-min") {
			algo := checksum.Algorithm(strings.TrimSuffix(data, "-min"))

			add := checksum.AlgorithmsByStrength.FilterMin(algo)
			if len(add) == 0 {
				return checksum.VerificationProfile{}, fmt.Errorf("unsupported algorithm: %s", algo)
			}

			res.Acceptable = append(res.Acceptable, add...)
			res.Required = true

			usedCustomAlgos = true
		} else {
			algo := checksum.Algorithm(data)

			if !checksum.AlgorithmsByStrength.Contains(algo) {
				return checksum.VerificationProfile{}, fmt.Errorf("unsupported algorithm: %s", data)
			}

			res.Acceptable = append(res.Acceptable, algo)
			res.Required = true

			usedCustomAlgos = true
		}
	}

	if len(standaloneValues) > 1 {
		return checksum.VerificationProfile{}, fmt.Errorf("standalone value combined with others: %s", strings.Join(standaloneValues, ", "))
	}

	if !usedCustomAlgos {
		res.Acceptable = checksum.AlgorithmsByStrength
	}

	return res, nil
}
