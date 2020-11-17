package main

import (
	"context"
	"testing"
	"time"
)

// Deadline を超えた時に、どういう値になる？
func TestContext_Deadline(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	const key = "hello"
	const value = "world"
	ctx = context.WithValue(ctx, key, value)
	_, ok := ctx.Deadline()
	if !ok {
		t.Error("deadline が存在していると ok になる")
	}
	if err := ctx.Err(); err != nil {
		t.Error("なんかエラーがあるぞ", err)
	}

	time.Sleep(250 * time.Millisecond) // deadline 超え

	if ctx.Value(key) != value {
		t.Error("deadline 超えると、Value が取れなくなる？")
	}
	if err := ctx.Err(); err == nil {
		t.Error("deadline exceeded がないだと！？")
	}
	if err := ctx.Err(); err != context.DeadlineExceeded {
		t.Error("DeadlineExceeded じゃない！？", err)
	}

	// deadline 超えている 親で context を新たに作ると、子どもはどうなってる？ -> deadline 超え状態のまま
	ctx, cancel = context.WithCancel(ctx)
	if ctx.Value(key) != value {
		t.Error("deadline 超えると、Value が取れなくなる？")
	}
	if err := ctx.Err(); err == nil {
		t.Error("deadline exceeded がないだと！？")
	}
	if err := ctx.Err(); err != context.DeadlineExceeded {
		t.Error("DeadlineExceeded じゃない！？", err)
	}
}
