package feedback

import (
	"context"

	"github.com/antgroup/aievo/schema"
)

type FeedbackType = string

const (
	Approved    FeedbackType = "Approved"
	NotApproved FeedbackType = "NotApproved"
)

type Feedback interface {
	Feedback(ctx context.Context, agent schema.Agent, messages []schema.Message, actions []schema.StepAction,
		steps []schema.StepAction, prompt string) *FeedbackInfo
}
