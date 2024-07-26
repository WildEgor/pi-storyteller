package tg

import tele "gopkg.in/telebot.v3"

// JobLimit ...
func JobLimit(next tele.HandlerFunc) tele.HandlerFunc {
	return func(c tele.Context) error {
		// TODO
		return next(c) // continue execution chain
	}
}
