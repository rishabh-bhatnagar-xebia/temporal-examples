package shared

type WorkflowIn struct {
	Data string
}

type WorkflowBasicOut struct {
	DBOut  *DBOut
	GitOut *GitOut
}

type WorkflowAsyncV1Out struct {
	DBOut *DBOut
}

type WorkflowAsyncV2Out struct {
	DBOut  *DBOut
	GitOut *GitOut
}

type DBOut struct {
	ID string
}

type GitOut struct {
	ID string
}
