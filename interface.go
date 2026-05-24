package response_analyser

import (
	"context"
	"io"
	"time"
)

type DataExporter interface {
	// Este é o seu ponto de entrada do hotpath. O hotpath chamará esta função e passará o conteúdo da resposta antes de escrevê-lo no corpo do http.Response.
	Analyse(context context.Context, content []byte)
	// Esta função deve retornar os valores atuais de cada código único. Ela não deve redefinir os contadores.
	GetCurrentCounts() map[string]uint64

	// Subscribe adiciona um novo assinante de relatório com o intervalo especificado
	Subscribe(writer io.WriteCloser, interval time.Duration) error
	// Espera-se que Shutdown seja chamado pelo aplicativo principal quando o servidor receber uma chamada de término com contexto para desligamento correto.
	Shutdown(ctx context.Context) error
}
