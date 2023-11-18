package future

func NewUnprocessResult() *UnProcessResult {
	return &UnProcessResult{
		resultFuture: make(map[string]*CompleteFuture),
	}
}

type UnProcessResult struct {
	resultFuture map[string]*CompleteFuture
}

func (unProcessResult *UnProcessResult) Complete(id string, result interface{}) {
	completeFuture := unProcessResult.resultFuture[id]
	complete := completeFuture.getComplete()
	complete <- result

}

func (unProcessResult *UnProcessResult) Put(id string, completeFuture *CompleteFuture) {
	unProcessResult.resultFuture[id] = completeFuture
}
