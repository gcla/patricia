package int32_tree

import (
	"math/bits"

	"github.com/gcla/patricia"
)

const _leftmost64Bit = uint64(1 << 63)

type treeNodeV6 struct {
	treeNode
	Left         uint // left node index: 0 for not set
	Right        uint // right node index: 0 for not set
	prefixLeft   uint64
	prefixRight  uint64
	prefixLength uint
	TagCount     int
}

func (n *treeNodeV6) MatchCount(address patricia.IPv6Address) uint {
	length := address.Length
	if length > n.prefixLength {
		length = n.prefixLength
	}

	matches := uint(bits.LeadingZeros64(n.prefixLeft ^ address.Left))
	if matches == 64 && length > 64 {
		matches += uint(bits.LeadingZeros64(n.prefixRight ^ address.Right))
	}
	if matches > length {
		return length
	}
	return matches
}

// ShiftPrefix shifts the prefix by the input shiftCount
func (n *treeNodeV6) ShiftPrefix(shiftCount uint) {
	n.prefixLeft, n.prefixRight, n.prefixLength = patricia.ShiftLeftIPv6(n.prefixLeft, n.prefixRight, n.prefixLength, shiftCount)
}

// IsLeftBitSet returns whether the leftmost bit is set
func (n *treeNodeV6) IsLeftBitSet() bool {
	return n.prefixLeft >= _leftmost64Bit
}

// MergeFromNodes updates the prefix and prefix length from the two input nodes
func (n *treeNodeV6) MergeFromNodes(left *treeNodeV6, right *treeNodeV6) {
	n.prefixLeft, n.prefixRight, n.prefixLength = patricia.MergePrefixes64(left.prefixLeft, left.prefixRight, left.prefixLength, right.prefixLeft, right.prefixRight, right.prefixLength)
}
