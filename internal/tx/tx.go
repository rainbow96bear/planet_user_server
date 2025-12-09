package tx

import (
	"context"

	"gorm.io/gorm"
)

// txKey는 context.Context에서 트랜잭션 객체를 저장하고 추출하는 데 사용되는 타입입니다.
// 충돌을 방지하기 위해 빈 구조체를 사용합니다.
type txKey struct{}

// WithTx는 주어진 Context에 GORM 트랜잭션(tx)을 담아 새 Context를 반환합니다.
func WithTx(ctx context.Context, tx *gorm.DB) context.Context {
	return context.WithValue(ctx, txKey{}, tx)
}

// GetTx는 Context에서 GORM 트랜잭션을 추출합니다.
// 트랜잭션이 없으면 nil을 반환합니다.
func GetTx(ctx context.Context) *gorm.DB {
	if tx, ok := ctx.Value(txKey{}).(*gorm.DB); ok {
		return tx
	}
	return nil // 트랜잭션이 Context에 없을 경우
}

// BeginTx는 DB에서 새 트랜잭션을 시작하고 Context에 담아 반환합니다.
func BeginTx(ctx context.Context, db *gorm.DB) (*gorm.DB, context.Context, error) {
	tx := db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return nil, nil, tx.Error
	}
	newCtx := WithTx(ctx, tx)
	return tx, newCtx, nil
}
