package middleware

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"sort"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"

	"github.com/lingvo-ai/lingvo/internal/db"
)

type telegramUser struct {
	ID           int64  `json:"id"`
	Username     string `json:"username"`
	LanguageCode string `json:"language_code"`
}

func AuthMiddleware(botToken string, database *sqlx.DB, devMode bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		if devMode {
			log.Printf("[auth] DEV_MODE active, bypassing auth for %s", c.ClientIP())
			mockTGID := int64(12345)
			if err := db.UpsertUser(c.Request.Context(), database, mockTGID, "Dev", "uz"); err != nil {
				c.AbortWithStatusJSON(500, gin.H{"error": "internal_error"})
				return
			}
			user, err := db.GetUserByTelegramID(c.Request.Context(), database, mockTGID)
			if err != nil {
				c.AbortWithStatusJSON(500, gin.H{"error": "internal_error"})
				return
			}
			c.Set("telegram_id", mockTGID)
			c.Set("user_id", user.ID)
			c.Set("lang", "uz")
			c.Set("level", user.Level)
			c.Next()
			return
		}

		initData := c.GetHeader("X-Telegram-Init-Data")

		var tgUser *telegramUser
		var err error

		if initData == "" {
			log.Printf("[auth] missing init data header from %s", c.ClientIP())
			c.AbortWithStatusJSON(401, gin.H{"error": "unauthorized", "message": "Missing init data"})
			return
		}

		tgUser, err = verifyInitData(initData, botToken)
		if err != nil {
			log.Printf("[auth] verifyInitData failed: %v (initData length=%d, ip=%s)", err, len(initData), c.ClientIP())
			c.AbortWithStatusJSON(401, gin.H{"error": "unauthorized", "message": "Invalid init data"})
			return
		}

		lang := tgUser.LanguageCode
		if lang != "uz" && lang != "ru" {
			lang = "uz"
		}

		if err := db.UpsertUser(c.Request.Context(), database, tgUser.ID, tgUser.Username, lang); err != nil {
			c.AbortWithStatusJSON(500, gin.H{"error": "internal_error"})
			return
		}

		user, err := db.GetUserByTelegramID(c.Request.Context(), database, tgUser.ID)
		if err != nil {
			c.AbortWithStatusJSON(500, gin.H{"error": "internal_error"})
			return
		}

		c.Set("telegram_id", tgUser.ID)
		c.Set("user_id", user.ID)
		c.Set("lang", user.Lang)
		c.Set("level", user.Level)
		c.Next()
	}
}

func verifyInitData(initData, botToken string) (*telegramUser, error) {
	values, err := url.ParseQuery(initData)
	if err != nil {
		return nil, err
	}

	hash := values.Get("hash")
	if hash == "" {
		return nil, fmt.Errorf("no hash")
	}
	values.Del("hash")

	keys := make([]string, 0, len(values))
	for k := range values {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var pairs []string
	for _, k := range keys {
		for _, v := range values[k] {
			pairs = append(pairs, fmt.Sprintf("%s=%s", k, v))
		}
	}
	dataCheck := strings.Join(pairs, "\n")

	secret := hmac.New(sha256.New, []byte("WebAppData"))
	secret.Write([]byte(botToken))
	secretKey := secret.Sum(nil)

	mac := hmac.New(sha256.New, secretKey)
	mac.Write([]byte(dataCheck))
	expectedHash := hex.EncodeToString(mac.Sum(nil))

	if !hmac.Equal([]byte(hash), []byte(expectedHash)) {
		return nil, fmt.Errorf("invalid hash: expected=%s got=%s keys=%v initData=%q", expectedHash, hash, keys, initData)
	}

	userStr := values.Get("user")
	if userStr == "" {
		return nil, fmt.Errorf("no user data")
	}

	var tgUser telegramUser
	if err := json.Unmarshal([]byte(userStr), &tgUser); err != nil {
		return nil, err
	}

	return &tgUser, nil
}
