# Bitmap

One of the more perplexing things from the paper ["Ideal Hash Trees" by Phil Bagwell](https://lampwww.epfl.ch/papers/idealhashtrees.pdf) is the handling of the bitmap.

It describes a way to map a group of bits onto a sparse array without actually allocating a sparse array. The bitmap will be sparse but the array will be packed and we can grow it as needed. I find this fascinating. This works like this:

```go
func set[S ~[]E, E any](bitmap uint32, slice S, hash uint32, value E) (uint32, []E) {
	// Map the hash to a bit position
	bit := uint32(1) << hash

	// Use it to mask the bitmap and compute the population count
	index := bits.OnesCount32(bitmap & (bit - 1))

	// If the bit position is vacant. We insert at this index
	if bitmap&bit == 0 {
		bitmap, slice = bitmap|bit, slices.Insert(slice, index, value)
	} else {
		slice[index] = value
	}

	return bitmap, slice
}
```

The function above is part of a small [Go program](https://play.golang.com/p/18QCdZ2i9DA) that illustrates how we build a sorted array using this method. What we're doing is that we're making a linear search in "chunks" of this word size. A word size of 32-bits will use a chunk size of 5 bits and a word size of 64-bits will use a chunk size of 6 bits.
