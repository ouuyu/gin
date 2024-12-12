package common

import (
	"sync"
	"time"
)

type VerificationCode struct {
    Code      string
    CreatedAt time.Time
}

var (
    verificationCodes = make(map[string]VerificationCode)
    codeMutex        sync.RWMutex
)

// VerifyCodeWithKey 验证邮箱验证码
func VerifyCodeWithKey(email, code, purpose string) bool {
    key := email + ":" + purpose
    codeMutex.RLock()
    defer codeMutex.RUnlock()
    
    if storedCode, exists := verificationCodes[key]; exists {
        // 验证码5分钟内有效
        if time.Since(storedCode.CreatedAt) <= 5*time.Minute && storedCode.Code == code {
            delete(verificationCodes, key)
            return true
        }
    }
    return false
}

// SaveVerificationCode 保存验证码
func SaveVerificationCode(email, code, purpose string) {
    key := email + ":" + purpose
    codeMutex.Lock()
    defer codeMutex.Unlock()
    
    verificationCodes[key] = VerificationCode{
        Code:      code,
        CreatedAt: time.Now(),
    }
}