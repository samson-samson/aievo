package feedback

type LLMFeedbackOption func(*LLMFeedback)

func WithPromptTemplate(template string) LLMFeedbackOption {
	return func(fd *LLMFeedback) {
		fd.prompt = template
	}
}

func WithMaxConBlock(block int) LLMFeedbackOption {
	return func(fd *LLMFeedback) {
		fd.maxcb = block
	}
}

func WithMiddlewares(middlewares ...Middleware) LLMFeedbackOption {
	return func(fd *LLMFeedback) {
		fd.middlewaresChain = middlewareChain(middlewares...)
	}
}

func WithExpertNum(num int) LLMFeedbackOption {
	return func(fd *LLMFeedback) {
		fd.expert = num
	}
}
