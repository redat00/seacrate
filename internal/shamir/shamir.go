package shamir

import (
	"fmt"

	"github.com/redat00/seacrate/internal/shamir/sss"
)

func Split(toSplit []byte, partCount int, thresholdCount int) ([][]byte, error) {
	d, err := sss.Split(toSplit, partCount, thresholdCount)
	if err != nil {
		return [][]byte{}, fmt.Errorf("could not split secret using the Shamir algorithm due to the following error : %v", err.Error())
	}
	return d, nil
}

func Combine(toCombine [][]byte) ([]byte, error) {
	d, err := sss.Combine(toCombine)
	if err != nil {
		return []byte{}, fmt.Errorf("could not combine using the Shamir algorithm due to the following error : %v", err)
	}
	return d, nil
}
