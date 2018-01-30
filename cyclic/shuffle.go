package cyclic

// Shuffles a passed interface array using a Fisher-Yates shuffle
func Shuffle(shufflee *[]uint64) {
	// Skip empty lists or lists of only 1 element, they cannot be shuffled
	if len(*shufflee) <= 1 {
		return
	}

	g := NewRandom(NewInt(0), NewInt(int64(len(*shufflee))-1))

	x := NewInt(1<<63 - 1)

	for curPos := int64(0); curPos < int64(len(*shufflee)); curPos++ {
		randPos := g.Rand(x).Int64()
		(*shufflee)[randPos], (*shufflee)[curPos] = (*shufflee)[curPos], (*shufflee)[randPos]
	}

}
