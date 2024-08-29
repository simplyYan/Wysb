package memorymanager

import (
	"fmt"
	"sync"
)

type MemoryManager struct {
	memory []byte
	free   []bool
	mu     sync.Mutex
}

func New(size int) *MemoryManager {
	return &MemoryManager{
		memory: make([]byte, size),
		free:   make([]bool, size),
	}
}

func (mm *MemoryManager) Allocate(size int) (int, error) {
	mm.mu.Lock()
	defer mm.mu.Unlock()

	start := -1
	count := 0
	for i := 0; i < len(mm.memory); i++ {
		if !mm.free[i] {
			if start == -1 {
				start = i
			}
			count++
			if count == size {
				for j := start; j < start+size; j++ {
					mm.free[j] = true
				}
				return start, nil
			}
		} else {
			start = -1
			count = 0
		}
	}
	return -1, fmt.Errorf("memória insuficiente")
}

func (mm *MemoryManager) Free(address int, size int) error {
	mm.mu.Lock()
	defer mm.mu.Unlock()

	if address < 0 || address+size > len(mm.memory) {
		return fmt.Errorf("endereço fora dos limites")
	}

	for i := address; i < address+size; i++ {
		if !mm.free[i] {
			return fmt.Errorf("memória já liberada ou inválida")
		}
		mm.free[i] = false
	}
	return nil
}

func (mm *MemoryManager) Read(address int) (byte, error) {
	mm.mu.Lock()
	defer mm.mu.Unlock()

	if address < 0 || address >= len(mm.memory) || !mm.free[address] {
		return 0, fmt.Errorf("endereço inválido ou memória não alocada")
	}
	return mm.memory[address], nil
}

func (mm *MemoryManager) Write(address int, value byte) error {
	mm.mu.Lock()
	defer mm.mu.Unlock()

	if address < 0 || address >= len(mm.memory) || !mm.free[address] {
		return fmt.Errorf("endereço inválido ou memória não alocada")
	}
	mm.memory[address] = value
	return nil
}

func (mm *MemoryManager) Dump() {
	mm.mu.Lock()
	defer mm.mu.Unlock()

	for i := 0; i < len(mm.memory); i++ {
		if mm.free[i] {
			fmt.Printf("X")
		} else {
			fmt.Printf(".")
		}
		if (i+1)%10 == 0 {
			fmt.Println()
		}
	}
}