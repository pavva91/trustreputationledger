package arglib

func ParseStringSliceToByteSlice(args []string) (argsAsBytes [][]byte){
	for _,arg := range args {
		argAsBytes := []byte(arg)
		argsAsBytes = append(argsAsBytes,argAsBytes)
	}
	return argsAsBytes
}
