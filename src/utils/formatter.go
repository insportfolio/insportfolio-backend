package utils

import (
	"errors"
	"strconv"
	"strings"
)

func StringToArray(data string) ([]int, error) {
	data = strings.Replace(data, "[", "", 1)
	data = strings.Replace(data, "]", "", 1)

	splitted := strings.Split(data, ",")
	var toInt []int

	for _, c := range splitted {
		i, err := strconv.Atoi(c)
		if err != nil {
			return nil, errors.New("Invalid integer: " + c)
		}
		toInt = append(toInt, i)
	}

	return toInt, nil
}
