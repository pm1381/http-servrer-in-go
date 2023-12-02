package pkg

type chunk struct {
	bufSize int
	offset int64
}

func fillChunkSizes(concurrencySteps int, fileSize int) ([]chunk) {
	chunkSizes := make([]chunk, concurrencySteps)
	for i := 0; i < concurrencySteps; i++ {
		chunkSizes[i].bufSize = bufferSize
		chunkSizes[i].offset = int64(i * bufferSize)
		if (i == concurrencySteps - 1) {
			chunkSizes[i].bufSize = fileSize % bufferSize
		}
	}
	return chunkSizes	
}
