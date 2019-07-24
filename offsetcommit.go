package kafka

type offsetCommitRequestV2Partition struct {
	// Partition ID
	Partition int32

	// Offset to be committed
	Offset int64

	// Metadata holds any associated metadata the client wants to keep
	Metadata string
}

func (t offsetCommitRequestV2Partition) size() int32 {
	return sizeofInt32(t.Partition) +
		sizeofInt64(t.Offset) +
		sizeofString(t.Metadata)
}

func (t offsetCommitRequestV2Partition) writeTo(wb *writeBuffer) {
	wb.writeInt32(t.Partition)
	wb.writeInt64(t.Offset)
	wb.writeString(t.Metadata)
}

type offsetCommitRequestV2Topic struct {
	// Topic name
	Topic string

	// Partitions to commit offsets
	Partitions []offsetCommitRequestV2Partition
}

func (t offsetCommitRequestV2Topic) size() int32 {
	return sizeofString(t.Topic) +
		sizeofArray(len(t.Partitions), func(i int) int32 { return t.Partitions[i].size() })
}

func (t offsetCommitRequestV2Topic) writeTo(wb *writeBuffer) {
	wb.writeString(t.Topic)
	wb.writeArray(len(t.Partitions), func(i int) { t.Partitions[i].writeTo(wb) })
}

type offsetCommitRequestV2 struct {
	// GroupID holds the unique group identifier
	GroupID string

	// GenerationID holds the generation of the group.
	GenerationID int32

	// MemberID assigned by the group coordinator
	MemberID string

	// RetentionTime holds the time period in ms to retain the offset.
	RetentionTime int64

	// Topics to commit offsets
	Topics []offsetCommitRequestV2Topic
}

func (t offsetCommitRequestV2) size() int32 {
	return sizeofString(t.GroupID) +
		sizeofInt32(t.GenerationID) +
		sizeofString(t.MemberID) +
		sizeofInt64(t.RetentionTime) +
		sizeofArray(len(t.Topics), func(i int) int32 { return t.Topics[i].size() })
}

func (t offsetCommitRequestV2) writeTo(wb *writeBuffer) {
	wb.writeString(t.GroupID)
	wb.writeInt32(t.GenerationID)
	wb.writeString(t.MemberID)
	wb.writeInt64(t.RetentionTime)
	wb.writeArray(len(t.Topics), func(i int) { t.Topics[i].writeTo(wb) })
}

type offsetCommitResponseV2PartitionResponse struct {
	Partition int32

	// ErrorCode holds response error code
	ErrorCode int16
}

func (t offsetCommitResponseV2PartitionResponse) size() int32 {
	return sizeofInt32(t.Partition) +
		sizeofInt16(t.ErrorCode)
}

func (t offsetCommitResponseV2PartitionResponse) writeTo(wb *writeBuffer) {
	wb.writeInt32(t.Partition)
	wb.writeInt16(t.ErrorCode)
}

func (t *offsetCommitResponseV2PartitionResponse) readFrom(rb *readBuffer) {
	t.Partition = rb.readInt32()
	t.ErrorCode = rb.readInt16()
}

type offsetCommitResponseV2Response struct {
	Topic              string
	PartitionResponses []offsetCommitResponseV2PartitionResponse
}

func (t offsetCommitResponseV2Response) size() int32 {
	return sizeofString(t.Topic) +
		sizeofArray(len(t.PartitionResponses), func(i int) int32 { return t.PartitionResponses[i].size() })
}

func (t offsetCommitResponseV2Response) writeTo(wb *writeBuffer) {
	wb.writeString(t.Topic)
	wb.writeArray(len(t.PartitionResponses), func(i int) { t.PartitionResponses[i].writeTo(wb) })
}

func (t *offsetCommitResponseV2Response) readFrom(rb *readBuffer) {
	t.Topic = rb.readString()

	rb.readArray(func() {
		res := offsetCommitResponseV2PartitionResponse{}
		res.readFrom(rb)
		t.PartitionResponses = append(t.PartitionResponses, res)
	})
}

type offsetCommitResponseV2 struct {
	Responses []offsetCommitResponseV2Response
}

func (t offsetCommitResponseV2) size() int32 {
	return sizeofArray(len(t.Responses), func(i int) int32 { return t.Responses[i].size() })
}

func (t offsetCommitResponseV2) writeTo(wb *writeBuffer) {
	wb.writeArray(len(t.Responses), func(i int) { t.Responses[i].writeTo(wb) })
}

func (t *offsetCommitResponseV2) readFrom(rb *readBuffer) {
	rb.readArray(func() {
		res := offsetCommitResponseV2Response{}
		res.readFrom(rb)
		t.Responses = append(t.Responses, res)
	})
}
