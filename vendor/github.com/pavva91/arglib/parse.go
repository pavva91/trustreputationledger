package arglib

import (
	"strings"
)

func ParseStringSliceToByteSlice(args []string) (argsAsBytes [][]byte){
	for _,arg := range args {
		argAsBytes := []byte(arg)
		argsAsBytes = append(argsAsBytes,argAsBytes)
	}
	return argsAsBytes
}

func ParseStringToStringSlice(stringToDecompose string) (stringSlice []string){
	if stringToDecompose=="" {
		return stringSlice
	}
	stringSlice = strings.Split(stringToDecompose, ",")
	return stringSlice
}