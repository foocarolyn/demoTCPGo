package main

func Byte32bArrToInt(byteArr []byte) int {
	var result int

	// 4 3 2 1 => 1 2 3 4
	result += int(byteArr[0])
	result += int(byteArr[1]) << 8
	result += int(byteArr[2]) << 16
	result += int(byteArr[3]) << 24

	return int(result)
}

func IntTo32bByteArr(targetInt int) []byte {
	byteArr := make([]byte, 4)

	// 1 2 3 4 => 4 3 2 1
	byteArr[0] = byte(targetInt & 0xff);
	byteArr[1] = byte((targetInt >> 8) & 0xff);
	byteArr[2] = byte((targetInt >> 16) & 0xff);
	byteArr[3] = byte((targetInt >> 32) & 0xff);

	return byteArr
}

