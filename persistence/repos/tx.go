package repos

type Tx interface {
	Commit() error
	Rollback() error
}
