package middleware

import "gopkg.in/telebot.v3"

// Middleware проверяет chat_id перед каждым хендлером
func AuthMiddleware(allowedChatID int64) telebot.MiddlewareFunc {
	return func(next telebot.HandlerFunc) telebot.HandlerFunc {
		return func(c telebot.Context) error {
			if c.Sender().ID != allowedChatID {
				return nil 
			}
			return next(c) 
		}
	}
}
