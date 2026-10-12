package cmd

import (
	"fmt"
	"os"

	sudoku "github.com/sagernet/sing-box/transport/sudoku"
)

// generateSudokuKeyPair prints a Sudoku key pair using the core's own
// generator, so the key format always matches what the inbound/outbound
// expects. Mirrors `sing-box generate sudoku-keypair`.
func generateSudokuKeyPair() {
	privateKey, publicKey, err := sudoku.GenKeyPair()
	if err != nil {
		fmt.Println("Failed to generate Sudoku key pair:", err)
		os.Exit(1)
	}
	fmt.Println("PrivateKey:", privateKey)
	fmt.Println("PublicKey: ", publicKey)
}
