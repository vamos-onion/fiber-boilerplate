package v1

import (
	"bufio"
	"errors"
	"fmt"
	"net/http"
	"time"

	"fiber-boilerplate/internal/defs"
	"fiber-boilerplate/internal/pkg/logging"
	"fiber-boilerplate/internal/pkg/session"

	"github.com/gofiber/fiber/v2"
)

// getRedisStore Redis Store 조회
func getRedisStore(ctx *fiber.Ctx) (*session.StoreBlock, error) {
	if store := ctx.Locals(defs.KeyStore); store != nil {
		if redisCtx, ok := store.(*session.StoreBlock); ok {
			return redisCtx, nil
		}
	}

	redisCtx := session.GetStore(defs.KeyStore)
	if redisCtx == nil {
		return nil, defs.ErrInvalid
	}

	return redisCtx, nil
}

func SseOpen(ctx *fiber.Ctx) error {
	// HTTP 헤더 설정
	ctx.Set("Content-Type", "text/event-stream")
	ctx.Set("Cache-Control", "no-cache")
	ctx.Set("Connection", "keep-alive")
	ctx.Set("Transfer-Encoding", "chunked")

	// Redis Store 초기화
	redisCtx, err := getRedisStore(ctx)
	if err != nil {
		return SendError(ctx, http.StatusInternalServerError, err)
	}

	// 메시지 채널과 context 미리 저장
	messageChan := make(chan string, 1)
	done := make(chan struct{})

	requestCtx := ctx.Context()

	// Redis Subscribe를 별도 goroutine에서 실행
	go func() {
		defer close(done)

		goerr := redisCtx.SubscribeChannel(requestCtx, defs.KeyStore)

		// context가 이미 cancel되었으면 보내지 않음
		select {
		case <-requestCtx.Done():
			return
		default:
		}

		var msg string
		switch {
		case goerr == nil:
			msg = "End"
		case errors.Is(goerr, defs.ErrTimeout):
			msg = "Timeout"
		default:
			logging.Error(goerr, "Redis subscribe error")
			msg = fmt.Sprintf("Error: %s", goerr.Error())
		}

		// 무한 블로킹을 피하기 위한 논블로킹 전송
		select {
		case messageChan <- msg:
			logging.Trace("Sent message to channel: %s", msg)
		case <-requestCtx.Done():
			logging.Trace("Request context already done, skipping message send")
		default:
			logging.Trace("Message channel full or receiver closed, skipping send")
		}
	}()

	// SSE 스트리밍
	ctx.Status(http.StatusOK).Context().SetBodyStreamWriter(func(w *bufio.Writer) {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		var i int = 1
		for {
			select {
			case <-requestCtx.Done():
				logging.Trace("Client connection closed")
				return
			case msg := <-messageChan:
				// Redis 메시지 수신
				logging.Trace("Received Redis message: %s", msg)
				message := fmt.Sprintf("data: %s\n\n", msg)
				if _, err := w.Write([]byte(message)); err != nil {
					logging.Trace("Write error: %v", err)
					return
				}
				if err := w.Flush(); err != nil {
					logging.Trace("Flush error: %v", err)
					return
				}
				return // Redis 메시지 수신 후 종료
			case <-done:
				// Goroutine이 종료됨 (메시지를 보낼 수 없는 경우)
				logging.Trace("Subscribe goroutine completed")
				return
			case <-ticker.C:
				// 주기적으로 heartbeat 전송
				message := fmt.Sprintf("data: %d sec\n\n", i)
				if _, err := w.Write([]byte(message)); err != nil {
					logging.Trace("Write error: %v", err)
					return
				}
				if err := w.Flush(); err != nil {
					logging.Trace("Flush error: %v", err)
					return
				}
				i++
			}
		}
	})

	return nil
}

func SseClose(ctx *fiber.Ctx) error {
	// Redis Store 초기화
	redisCtx, err := getRedisStore(ctx)
	if err != nil {
		return SendError(ctx, http.StatusInternalServerError, err)
	}

	// Publish message to channel
	requestCtx := ctx.Context()
	err = redisCtx.PublishChannel(requestCtx, defs.KeyStore, "Hello World")
	if err != nil {
		logging.Error(err, "Failed to publish SSE close message")
		return SendError(ctx, http.StatusInternalServerError, err)
	}

	return SendGeneric(ctx, http.StatusOK)
}
