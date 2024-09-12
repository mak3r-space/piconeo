package main

func newAnikaLetters() letterConfig {
	letters := make([]letter, 83)
	// A
	letters[52] = letter{0, 0}
	letters[53] = letter{0, 0}
	letters[51] = letter{0, 1}
	letters[54] = letter{0, 1}
	letters[50] = letter{0, 2}
	letters[55] = letter{0, 2}
	letters[49] = letter{0, 3}
	letters[56] = letter{0, 3}
	letters[48] = letter{0, 4}
	letters[57] = letter{0, 4}
	letters[47] = letter{0, 5}
	letters[58] = letter{0, 5}
	letters[46] = letter{0, 6}
	letters[59] = letter{0, 6}
	letters[45] = letter{0, 7}
	letters[60] = letter{0, 7}
	letters[44] = letter{0, 8}
	letters[61] = letter{0, 8}
	letters[43] = letter{0, 9}
	letters[42] = letter{0, 10}
	letters[41] = letter{0, 11}
	letters[40] = letter{0, 12}
	letters[39] = letter{0, 13}
	letters[62] = letter{0, 13}
	letters[38] = letter{0, 14}
	letters[63] = letter{0, 14}

	// n
	letters[37] = letter{1, 15}
	letters[64] = letter{1, 15}
	letters[36] = letter{1, 16}
	letters[35] = letter{1, 17}
	letters[34] = letter{1, 18}
	letters[65] = letter{1, 19}
	letters[66] = letter{1, 20}
	letters[33] = letter{1, 21}
	letters[32] = letter{1, 22}
	letters[67] = letter{1, 22}
	letters[31] = letter{1, 23}
	letters[30] = letter{1, 24}
	letters[29] = letter{1, 25}
	letters[68] = letter{1, 25}
	letters[69] = letter{1, 26}
	letters[28] = letter{1, 25}
	letters[70] = letter{1, 27}

	// i
	letters[27] = letter{2, 28}
	letters[26] = letter{2, 29}
	letters[23] = letter{2, 30}
	letters[22] = letter{2, 31}
	letters[71] = letter{2, 32}
	letters[21] = letter{2, 33}
	letters[25] = letter{2, 34}
	letters[24] = letter{2, 34}

	//k
	letters[20] = letter{3, 35}
	letters[19] = letter{3, 36}
	letters[18] = letter{3, 37}
	letters[17] = letter{3, 38}
	letters[16] = letter{3, 39}
	letters[15] = letter{3, 40}
	letters[14] = letter{3, 41}
	letters[13] = letter{3, 42}
	letters[12] = letter{3, 43}
	letters[72] = letter{3, 44}
	letters[73] = letter{3, 45}
	letters[11] = letter{3, 45}
	letters[74] = letter{3, 46}
	letters[10] = letter{3, 46}
	letters[75] = letter{3, 47}
	letters[9] = letter{3, 47}
	letters[76] = letter{3, 48}
	letters[8] = letter{3, 48}
	letters[7] = letter{3, 49}
	letters[77] = letter{3, 50}
	letters[6] = letter{3, 50}
	letters[78] = letter{3, 51}
	letters[5] = letter{3, 51}
	letters[79] = letter{3, 52}

	// a
	letters[2] = letter{4, 53}
	letters[3] = letter{4, 54}
	letters[4] = letter{4, 55}
	letters[80] = letter{4, 56}
	letters[81] = letter{4, 57}
	letters[82] = letter{4, 58}
	letters[1] = letter{4, 59}
	letters[0] = letter{4, 60}

	return letterConfig{
		numLetters:     5, // anika
		numDrawIndices: 61,
		letters:        letters,
		startIdx:       27,
	}
}
