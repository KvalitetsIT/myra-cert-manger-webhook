package client

import (
	"context"
	"log/slog"
	"runtime"
	"time"
)

// ClientLogger is a generic decorator for any Client[T] implementation that adds structured logging.
//
// It wraps a Client[T] and logs calls to its OnAdd and OnDelete methods,
// including both the input record and any errors returned by the client.
//
// This pattern allows separation of concerns: the underlying client focuses
// purely on business logic, while logging is handled by this wrapper.
type ClientLogger[T any] struct {
	// client is the underlying Client[T] that performs the actual operations.
	client Client[T]

	// logger is the slog.Logger used for structured logging.
	logger *slog.Logger
}

// NewClientLogger creates a new ClientLogger[T] that wraps the given client
// and logs to the provided slog.Logger.
//
// Parameters:
//   - client: The Client[T] to wrap
//   - logger: The slog.Logger used for structured logging
//
// Returns:
//   - A new *ClientLogger[T] instance that implements Client[T].
func NewClientLogger[T any](client Client[T], logger *slog.Logger) *ClientLogger[T] {
	return &ClientLogger[T]{
		client: client,
		logger: logger,
	}
}

// OnDelete wraps the underlying client's OnDelete method with logging.
//
// Logs the input record before calling the underlying client,
// and logs any error returned as an error-level message.
//
// Parameters:
//   - record: The record to delete.
//
// Returns:
//   - The result returned by the underlying client's OnDelete method.
//   - Any error returned by the underlying client's OnDelete method.
func (p *ClientLogger[T]) OnDelete(record T) (T, error) {
	response, err := p.client.OnDelete(record)

	ctx := context.Background()

	if err != nil {
		p.logWithCaller(ctx, slog.LevelError, "Failed to delete dns record",
			slog.Any("record", record),
			slog.Any("error", err),
		)
	} else {
		p.logWithCaller(ctx, slog.LevelInfo, "Deleted record",
			slog.Any("record", record),
		)
	}

	return response, err

}

// OnAdd wraps the underlying client's OnAdd method with logging.
//
// Logs the input record before calling the underlying client,
// and logs any error returned as an error-level message.
//
// Parameters:
//   - record: The record to add.
//
// Returns:
//   - The result returned by the underlying client's OnAdd method.
//   - Any error returned by the underlying client's OnAdd method.
// func (p *ClientLogger[T]) OnAdd(record T) (T, error) {
// 	response, err := p.client.OnAdd(record)
// 	if err != nil {
// 		p.logger.Error("Faied to create dns record", slog.Any("record", record), slog.Any("error", err))
// 	}
// 	p.logger.Info("Added record", slog.Any("record", record))
// 	return response, err
// }

func (p *ClientLogger[T]) OnAdd(record T) (T, error) {
	response, err := p.client.OnAdd(record)

	ctx := context.Background()

	if err != nil {
		p.logWithCaller(ctx, slog.LevelError, "Failed to create dns record",
			slog.Any("record", record),
			slog.Any("error", err),
		)
	} else {
		p.logWithCaller(ctx, slog.LevelInfo, "Added record",
			slog.Any("record", record),
		)
	}

	return response, err
}

// logWithCaller is a helper that captures the PC of the one calling the decorator
func (p *ClientLogger[T]) logWithCaller(ctx context.Context, level slog.Level, msg string, attrs ...slog.Attr) {
	if !p.logger.Enabled(ctx, level) {
		return
	}
	var pcs [1]uintptr
	runtime.Callers(3, pcs[:])
	pc := pcs[0]

	r := slog.NewRecord(time.Now(), level, msg, pc)
	r.AddAttrs(attrs...)
	_ = p.logger.Handler().Handle(ctx, r)
}
