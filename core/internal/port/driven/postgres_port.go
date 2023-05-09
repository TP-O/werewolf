package driven

type PostgrePort interface {
	// Querier

	Close() error
	// StoreGame(ctx context.Context, params *StoreGameParams) error
}
