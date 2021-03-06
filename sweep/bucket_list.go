package sweep

// bucket contains a set of inputs that are not mutually exclusive.
type bucket pendingInputs

// tryAdd tries to add a new input to this bucket.
func (b bucket) tryAdd(input *pendingInput) bool {
	exclusiveGroup := input.params.ExclusiveGroup
	if exclusiveGroup != nil {
		for _, input := range b {
			existingGroup := input.params.ExclusiveGroup
			if existingGroup != nil &&
				*existingGroup == *exclusiveGroup {

				return false
			}
		}
	}

	b[*input.OutPoint()] = input

	return true
}

// bucketList is a list of buckets that contain non-mutually exclusive inputs.
type bucketList struct {
	buckets []bucket
}

// add adds a new input. If the input is not accepted by any of the existing
// buckets, a new bucket will be created.
func (b *bucketList) add(input *pendingInput) {
	for _, existingBucket := range b.buckets {
		if existingBucket.tryAdd(input) {
			return
		}
	}

	// Create a new bucket and add the input. It is not necessary to check
	// the return value of tryAdd because it will always succeed on an empty
	// bucket.
	newBucket := make(bucket)
	newBucket.tryAdd(input)
	b.buckets = append(b.buckets, newBucket)
}
