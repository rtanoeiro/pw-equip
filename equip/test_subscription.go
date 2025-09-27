package equip

import (
	"fmt"
)

func TestSubscription() {
	fmt.Println("=== Teste do Sistema de Assinatura ===")

	// Display HWID
	fmt.Println("\n1. Obtendo HWID:")
	DisplayHWID()

	// Test subscription check
	fmt.Println("\n2. Testando verificação de assinatura:")
	isActive, err := CheckSubscription()
	if err != nil {
		fmt.Printf("Erro esperado (servidor não existe): %v\n", err)
	} else {
		fmt.Printf("Status da assinatura: %v\n", isActive)
	}

	fmt.Println("\n3. Testando com retry:")
	isActive, err = ValidateSubscriptionWithRetry(2)
	if err != nil {
		fmt.Printf("Erro esperado (servidor não existe): %v\n", err)
	} else {
		fmt.Printf("Status da assinatura com retry: %v\n", isActive)
	}
}
