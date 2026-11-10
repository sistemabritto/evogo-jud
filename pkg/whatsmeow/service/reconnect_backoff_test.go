package whatsmeow_service

import (
	"testing"
	"time"
)

func TestReconnectBackoff(t *testing.T) {
	id := "test-instance"
	want := []time.Duration{0, 15 * time.Second, 30 * time.Second, 60 * time.Second, 120 * time.Second, 240 * time.Second, 5 * time.Minute, 5 * time.Minute}
	for i, w := range want {
		if got := reconnectBackoff(id); got != w {
			t.Fatalf("reconexão #%d: esperado %s, veio %s", i+1, w, got)
		}
	}

	// Fora da janela de 10 min volta a reconectar sem espera.
	reconnectMu.Lock()
	old := time.Now().Add(-reconnectWindow - time.Second)
	for i := range reconnectHistory[id] {
		reconnectHistory[id][i] = old
	}
	reconnectMu.Unlock()
	if got := reconnectBackoff(id); got != 0 {
		t.Fatalf("depois da janela: esperado 0, veio %s", got)
	}

	// Instâncias são independentes.
	if got := reconnectBackoff("outra"); got != 0 {
		t.Fatalf("outra instância: esperado 0, veio %s", got)
	}
}
