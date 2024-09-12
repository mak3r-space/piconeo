package main

func newMaliLetters() letterConfig {
	// M: top start: 17,end: 36; bottom start: 36, end: 37
	// a: top start: 12,end: 17; bottom start: 37, end: 40
	// l: top start: 6, end: 12; bottom start: 40, end: 42
	// i: top start: 0, end: 6; bottom start: 42, end: 44

	letters := make([]letter, 44)
	// M
	letters[23] = letter{0, 0}
	letters[24] = letter{0, 1}
	letters[25] = letter{0, 1}
	letters[22] = letter{0, 1}
	letters[26] = letter{0, 2}
	letters[27] = letter{0, 2}
	letters[21] = letter{0, 2}
	letters[28] = letter{0, 3}
	letters[29] = letter{0, 3}
	letters[20] = letter{0, 3}
	letters[30] = letter{0, 4}
	letters[31] = letter{0, 4}
	letters[19] = letter{0, 4}
	letters[32] = letter{0, 5}
	letters[33] = letter{0, 5}
	letters[18] = letter{0, 5}
	letters[34] = letter{0, 6}
	letters[35] = letter{0, 6}
	letters[36] = letter{0, 6}
	letters[17] = letter{0, 6}
	// a
	letters[37] = letter{1, 7}
	letters[16] = letter{1, 7}
	letters[38] = letter{1, 8}
	letters[15] = letter{1, 8}
	letters[14] = letter{1, 8}
	letters[39] = letter{1, 9}
	letters[13] = letter{1, 9}
	letters[12] = letter{1, 9}
	// l
	letters[40] = letter{2, 10}
	letters[11] = letter{2, 10}
	letters[10] = letter{2, 11}
	letters[9] = letter{2, 12}
	letters[8] = letter{2, 13}
	letters[7] = letter{2, 14}
	letters[6] = letter{2, 14}
	// i
	letters[41] = letter{3, 16}
	letters[5] = letter{3, 17}
	letters[4] = letter{3, 18}
	letters[3] = letter{3, 19}
	letters[2] = letter{3, 20}
	letters[1] = letter{3, 21}
	letters[42] = letter{3, 21}
	letters[0] = letter{3, 22}
	letters[43] = letter{3, 22}

	return letterConfig{
		numLetters:     4,
		numDrawIndices: 22,
		letters:        letters,
		startIdx:       27,
	}
}
