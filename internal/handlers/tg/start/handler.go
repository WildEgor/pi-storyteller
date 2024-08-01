// Package tg_start_handler responsible to show instructions
package tg_start_handler

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/WildEgor/pi-storyteller/internal/adapters/bot"
)

type Layout struct {
	Lang     string
	Guide    string
	Commands []string
}

type Lang struct {
	Ru string
	En string
}

type WelcomeData struct {
	Messages         Lang
	GenerateExamples []Lang
	RandomExamples   []Lang
}

var data = &WelcomeData{
	Messages: Lang{
		En: `
Welcome, %s! 🎉

Hello and thank you for starting me!

🔍 Key Features:
- /generate [text] - use this command to create content based on the text you provide. 
Simply type /generate followed by your text, and I’ll generate images related to your input.
Please note that sentences should be separated by periods. The text should not contain any offensive language and should not be too lengthy!

- /random - this command provides you with a random story with images. It’s perfect for when you need a quick bit of inspiration or just want to learn something new!

Example:

`,
		Ru: `
Добро пожаловать, %s! 🎉

Привет и спасибо, что запустили меня!

🔍 Основные функции:
- /generate [текст] - используйте эту команду для создания контента на основе предоставленного вами текста. 
Просто напишите /generate, за которым следует ваш текст, и я сгенерирую изображения, связанные с вашим вводом.
Учтите, что предложения должны быть разделены точкой, текст не должен содержать нецензурные выражения и быть слишком большим!

- /random - эта команда предоставляет вам случайную историю с картинками. Это идеально подходит, когда вам нужно быстрое вдохновение или просто хотите узнать что-то новое!

Пример:

`,
	},
	GenerateExamples: []Lang{
		{
			En: "/generate Geralt of Rivia accidentally getting stuck in a magical sauna, where all his attempts to escape are thwarted by enchanted towels and talking soap bars.",
			Ru: "/generate Геральт из Ривии, который случайно застрял в волшебной сауне, где все его попытки выбраться мешают заколдованные полотенца и говорящие мыльные пузыри.",
		},
		{
			En: "/generate Frodo Baggins trying to use a modern smartphone but keeps accidentally sending selfies to the Dark Lord. His quest turns into a comedic race to stop Sauron from discovering his silly photos.",
			Ru: "/generate Фродо Бэггинс, который пытается использовать современный смартфон, но всё время случайно отправляет селфи Темному Лорду. Его путешествие превращается в комедийную гонку, чтобы помешать Саурону увидеть его глупые фотографии.",
		},
		{
			En: "/generate Joker trying to host a cooking show. His attempts at making extravagant dishes lead to chaos in the kitchen, with ingredients exploding and a pie fight that ends with him covered in flour and cream.",
			Ru: "/generate Джокер, который пытается вести кулинарное шоу. Его попытки приготовить экстравагантные блюда приводят к хаосу на кухне, с взрывами ингредиентов и пирогами, в итоге он оказывается покрытым мукой и кремом.",
		},
	},
}

// StartHandler ...
type StartHandler struct {
	tgBot bot.Bot
}

// NewStartHandler ...
func NewStartHandler(tgBot bot.Bot) *StartHandler {
	return &StartHandler{
		tgBot,
	}
}

// Handle ...
func (h *StartHandler) Handle(ctx context.Context, payload *StartDTO) error {
	layout := &Layout{
		Lang:     payload.Lang,
		Commands: make([]string, 0),
		Guide:    data.Messages.En,
	}

	if payload.Lang == "ru" {
		layout.Guide = data.Messages.Ru
	}

	for _, item := range data.GenerateExamples {
		if payload.Lang == "ru" {
			layout.Commands = append(layout.Commands, item.Ru)
			continue
		}

		layout.Commands = append(layout.Commands, item.En)
	}

	//nolint
	_, err := h.tgBot.SendMsg(ctx, &bot.MessageRecipient{
		ID: payload.ChatID,
	}, fmt.Sprintf(layout.Guide, payload.Nickname))

	time.Sleep(5 * time.Second)

	//nolint
	_, err = h.tgBot.SendMsg(ctx, &bot.MessageRecipient{
		ID: payload.ChatID,
	}, layout.Commands[rand.Intn(len(layout.Commands))])

	return err
}
