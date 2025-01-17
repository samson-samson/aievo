package feedback

import (
	"context"

	"github.com/antgroup/aievo/schema"
)

// feedbackChain Define a struct contains slice of Feedback
type feedbackChain struct {
	chains []Feedback
}

// Feedback Implement the Feedback method for the feedbackChain
func (fc *feedbackChain) Feedback(ctx context.Context, agent schema.Agent, messages []schema.Message, actions []schema.StepAction,
	steps []schema.StepAction, prompt string) *FeedbackInfo {

	info := &FeedbackInfo{
		Type: Approved,
	}

	for _, feedback := range fc.chains {
		if feedback == nil {
			continue
		}
		info = feedback.Feedback(ctx, agent, messages, actions, steps, prompt)
		if info.Type == NotApproved {
			return info
		}
	}

	return info
}

// Chain function to create a new Feedback that chains multiple Feedback
func Chain(chains ...Feedback) Feedback {
	return &feedbackChain{chains: chains}
}
