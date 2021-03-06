// Copyright 2020 The Reed Developers
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

package discover

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/reed/log"
	"net"
	"testing"
)

func init() {
	log.Init()
}

func TestLogarithmDist(t *testing.T) {
	id1 := NodeID{31: byte(1)}
	id2 := NodeID{31: byte(2)}

	id3 := NodeID{31: byte(1)}

	dist := logarithmDist(id1, id2)
	dist2 := logarithmDist(id1, id3)

	if dist != 1 || dist2 != 0 {
		t.Error("logarithmDist error")
	}
}

func TestContains(t *testing.T) {
	id1, _ := hex.DecodeString("7b52009b64fd0a2a49e6d8a939753077792b0554")
	id2, _ := hex.DecodeString("40bd001563085fc35165329ea1ff5c5ecbdbbeef")
	id3, _ := hex.DecodeString("7110eda4d09e062aa5e4a390b0a572ac0d2c0220")

	var ns []*bNode
	ns = append(ns, &bNode{node: &Node{ID: BytesToHash(id1)}})
	ns = append(ns, &bNode{node: &Node{ID: BytesToHash(id2)}})

	if !contains(ns, BytesToHash(id1)) {
		t.Error("contains error (expect id in ns)")
	}

	if contains(ns, BytesToHash(id3)) {
		t.Error("contains error (expect ns does not exist request id)")
	}
}

func TestComputeDist(t *testing.T) {
	ta := NodeID{31: byte(3)}
	id1 := NodeID{31: byte(3)}
	id2 := NodeID{0, 0, 0, 0, 0, byte(8), 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, byte(2)}

	dist := computeDist(ta, id1, id2)
	fmt.Println(dist)
}

func TestNodesByDistance(t *testing.T) {
	nd := nodesByDistance{
		target: NodeID{31: byte(3)},
		entries: []*Node{
			{ID: NodeID{31: byte(1)}},
			{ID: NodeID{30: byte(2), 31: 0}},
			{ID: NodeID{9: byte(12), 29: byte(2)}},
		},
	}
	node := &Node{
		ID: NodeID{29: byte(2), 31: byte(1)},
	}
	nd.push(node)

	if len(nd.entries) != 4 {
		t.Error("failed to push node")
	}

	if !bytes.Equal(nd.entries[2].ID.Bytes(), node.ID.Bytes()) {
		t.Error("nodesByDistance push error")
	}
}

type tn struct {
	name string
	node *Node
}

func TestGetWithExclude(t *testing.T) {

	tb := newTable()

	minDist := []byte{30: byte(2), 31: byte(1)}
	secondDist := []byte{29: byte(3), 30: byte(2), 31: byte(1)}

	ns := tb.GetRandNodes(1, nil)
	if !bytes.Equal(ns[0].ID.Bytes(), minDist) {
		t.Fatal("the first(minimum distance) not right")
	}

	ns2 := tb.GetRandNodes(4, []NodeID{{29: byte(3), 30: byte(2), 31: byte(1)}})
	if len(ns2) != 4 {
		t.Fatal("wrong count")
	}
	for _, n := range ns2 {
		if bytes.Equal(n.ID.Bytes(), secondDist) {
			t.Fatal("does not exclude the given node")
		}
	}

}

func newTable() *Table {
	our := &Node{
		ID: NodeID{31: byte(1)},
	}
	tb, _ := NewTable(our)

	n2 := &Node{
		ID: NodeID{30: byte(2), 31: byte(1)},
	}
	n3 := &Node{
		ID:      NodeID{29: byte(3), 30: byte(2), 31: byte(1)},
		TCPPort: 8002,
		UDPPort: 8001,
		IP:      net.IP{123, 123, 123, 13},
	}
	n4 := &Node{
		ID: NodeID{28: byte(4), 29: byte(3), 30: byte(2), 31: byte(1)},
	}
	n5 := &Node{
		ID: NodeID{27: byte(5), 28: byte(4), 29: byte(3), 30: byte(2), 31: byte(1)},
	}
	n6 := &Node{
		ID: NodeID{26: byte(6), 27: byte(5), 28: byte(4), 29: byte(3), 30: byte(2), 31: byte(1)},
	}
	n7 := &Node{
		ID: NodeID{25: byte(7), 26: byte(6), 27: byte(5), 28: byte(4), 29: byte(3), 30: byte(2), 31: byte(1)},
	}
	var tns []*tn
	tns = append(tns, &tn{
		name: "n2",
		node: n2,
	})
	tns = append(tns, &tn{
		name: "n3",
		node: n3,
	})
	tns = append(tns, &tn{
		name: "n4",
		node: n4,
	})
	tns = append(tns, &tn{
		name: "n5",
		node: n5,
	})
	tns = append(tns, &tn{
		name: "n6",
		node: n6,
	})
	tns = append(tns, &tn{
		name: "n7",
		node: n7,
	})

	for _, t := range tns {
		tb.Add(t.node)
	}

	for i, b := range tb.Bucket {
		for _, n := range b {
			for _, t := range tns {
				if n.node == t.node {
					fmt.Printf("name:%s k-bucket:%d\n", t.name, i)
				}
			}
		}
	}
	return tb
}
