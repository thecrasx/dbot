package errors

import "fmt"

type FailedToCreateChannelID struct{}

func (e *FailedToCreateChannelID) Error() string {
	return "Failed to create ChannelID\n"
}

type UserMentionNotFound struct{}

func (e *UserMentionNotFound) Error() string {
	return "Failed to find user mention\n"
}

type InvalidRange struct {
	Start int
	End   int
}

func (e *InvalidRange) Error() string {
	return fmt.Sprintf("Range must be from %d to %d", e.Start, e.End)
}

type TTCNoAvailablePosition struct{}

func (e *TTCNoAvailablePosition) Error() string {
	return "No available position to place sign\n"
}

type TTCPositionTaken struct {
	Position int
}

func (e *TTCPositionTaken) Error() string {
	return fmt.Sprintf("Sign at [%d] position already exist\n", e.Position)
}
