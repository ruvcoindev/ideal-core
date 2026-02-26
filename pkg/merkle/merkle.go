package merkle

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
)

// MerkleNode представляет узел дерева Меркла
type MerkleNode struct {
	Hash  string
	Left  *MerkleNode
	Right *MerkleNode
	Data  string // ID связи или человека
}

// BuildTree строит дерево Меркла из списка хешей связей
func BuildTree(hashes []string) *MerkleNode {
	if len(hashes) == 0 {
		return nil
	}
	
	// Сортируем для детерминизма
	sort.Strings(hashes)
	
	// Создаём листовые узлы
	var nodes []*MerkleNode
	for _, h := range hashes {
		nodes = append(nodes, &MerkleNode{Hash: h, Data: h})
	}
	
	// Строим дерево по уровням
	for len(nodes) > 1 {
		var level []*MerkleNode
		for i := 0; i < len(nodes); i += 2 {
			left := nodes[i]
			var right *MerkleNode
			if i+1 < len(nodes) {
				right = nodes[i+1]
			} else {
				// Нечётное количество: дублируем последний
				right = left
			}
			combined := left.Hash + right.Hash
			hash := sha256.Sum256([]byte(combined))
			level = append(level, &MerkleNode{
				Hash:  hex.EncodeToString(hash[:]),
				Left:  left,
				Right: right,
			})
		}
		nodes = level
	}
	
	return nodes[0]
}

// Verify проверяет, входит ли хеш в дерево с данным root
func Verify(rootHash, targetHash, proof []string) bool {
	current := targetHash
	for _, sibling := range proof {
		combined := current + sibling
		hash := sha256.Sum256([]byte(combined))
		current = hex.EncodeToString(hash[:])
	}
	return current == rootHash
}

// HashRelation создаёт хеш для связи между двумя людьми
func HashRelation(fromID, toID, relationType string, timestamp int64) string {
	data := fmt.Sprintf("%s:%s:%s:%d", fromID, toID, relationType, timestamp)
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}
