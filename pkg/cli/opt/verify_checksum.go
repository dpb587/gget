package opt

import (
	"fmt"
	"strings"

	"github.com/dpb587/gget/pkg/checksum"
)

type VerifyChecksum struct {
	algorithms checksum.AlgorithmList
	mode       string
}

func (vc *VerifyChecksum) UnmarshalFlag(data string) error {
	if data == "required" {
		vc.algorithms = checksum.AlgorithmsByStrength
		vc.mode = "required"

		return nil
	} else if data == "auto" {
		// default; empty
		vc.algorithms = checksum.AlgorithmsByStrength

		return nil
	} else if data == "none" {
		vc.mode = "none"

		return nil
	} else if strings.HasSuffix(data, "-min") {
		algo := checksum.Algorithm(strings.TrimSuffix(data, "-min"))

		vc.algorithms = checksum.AlgorithmsByStrength.FilterMin(algo)
		if len(vc.algorithms) == 0 {
			return fmt.Errorf("unsupported algorithm: %s", algo)
		}

		vc.mode = "required"

		return nil
	}

	a := checksum.Algorithm(data)

	if !checksum.AlgorithmsByStrength.Contains(a) {
		return fmt.Errorf("unsupported algorithm: %s", data)
	}

	vc.algorithms = checksum.AlgorithmList{a}
	vc.mode = "required"

	return nil
}

func (vc VerifyChecksum) Mode() string {
	if vc.mode == "" {
		return "auto"
	}

	return vc.mode
}

func (vc VerifyChecksum) AcceptableAlgorithms() checksum.AlgorithmList {
	if vc.algorithms != nil {
		return vc.algorithms
	}

	return checksum.AlgorithmsByStrength
}
