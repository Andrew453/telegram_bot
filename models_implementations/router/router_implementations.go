package router

import (
	"TelegramBot/types/notes"
	"crypto/sha1"
	"fmt"
	"github.com/pkg/errors"
	"github.com/yanzay/tbot/v2"
	"regexp"
	"strings"
	"time"
)

type Storage map[string][]notes.Note //map[идентификатор клиента]map[хэш]заметка
type Router struct {
	storage Storage //
	Client  *tbot.Client
	NowNote notes.Note
}

func NewRouter(server *tbot.Server) *Router {
	var err error
	defer func() {
		err = errors.Wrap(err, "NewRouter")
	}()
	router := &Router{}
	router.storage = make(Storage)
	router.Client = server.Client()
	return router
}

func (r Router) CreateNoteText(message *tbot.Message) {
	err := r.Client.SendChatAction(message.Chat.ID, tbot.ActionTyping)
	if err != nil {
		return
	}
	time.Sleep(time.Millisecond)
	r.Client.SendMessage(message.Chat.ID, fmt.Sprintf("Просто напиши свою заметку в формате: \nТип:<твой_текст> Имя:<твой_текст> Текст:<твой_текст>"))
}

func (r *Router) ReadNoteByName(message *tbot.Message) {
	err := r.Client.SendChatAction(message.Chat.ID, tbot.ActionTyping)
	if err != nil {
		return
	}
	time.Sleep(time.Millisecond)
	text := message.Text
	reg := regexp.MustCompile(`/note (.+)`)
	arr := reg.FindStringSubmatch(text)
	if len(arr) != 2 {
		r.Client.SendMessage(message.Chat.ID, fmt.Sprintf("К сожалению, ты говоришь на незнакомом мне языке:("))
		return
	}
	arrNote := strings.Split(arr[1], "/////")
	for _, note := range r.storage[message.Chat.ID] {
		if note.Name == strings.TrimPrefix(arrNote[0], " ") && note.Class == strings.TrimPrefix(arrNote[1], " ") {
			r.Client.SendMessage(message.Chat.ID, fmt.Sprintf("Тип:"+
				"	%s\n\nНазвание:"+
				"	%s\n\nТекст:\n"+
				"%s", note.Class, note.Name, note.Text))
			return
		}
	}
	msg := "Прости, но кажется, что ты не создавал такую заметку"
	r.Client.SendMessage(message.Chat.ID, msg)
}

// ReadAllNotes функция получения всех записей
func (r *Router) ReadAllNotes(message *tbot.Message) {
	err := r.Client.SendChatAction(message.Chat.ID, tbot.ActionTyping)
	if err != nil {
		return
	}
	time.Sleep(time.Millisecond)
	//var list string
	if len(r.storage[message.Chat.ID]) == 0 {
		r.Client.SendMessage(message.Chat.ID, "Ты меня за кого держишь? Нет заметок, что ищешь?:)")
		return
	}
	buttons := r.makeButtonsOfNotes(message.Chat.ID)
	r.Client.SendMessage(message.Chat.ID, "Вот список твоих записей:", tbot.OptInlineKeyboardMarkup(buttons))
}

func (r *Router) Help(message *tbot.Message) {
	err := r.Client.SendChatAction(message.Chat.ID, tbot.ActionTyping)
	if err != nil {
		r.Client.SendMessage(message.Chat.ID, "Some error, sorry:(")
		return
	}
	time.Sleep(time.Millisecond)
	commands, err := r.Client.GetMyCommands()
	if err != nil {
		r.Client.SendMessage(message.Chat.ID, "Some error, sorry:(")
		return
	}
	msg := ""
	for _, command := range *commands {
		msg = fmt.Sprintf("%sКоманда \"%s\" \nОписание: %s\n\n", msg, command.Command, command.Description)

	}
	buttons := makeButtons()
	r.Client.SendMessage(message.Chat.ID, msg, tbot.OptInlineKeyboardMarkup(buttons))
	return
}

func (r *Router) StartMsg(m *tbot.Message) {
	buttons := makeButtons()
	_, err := r.Client.SendPhotoFile(m.Chat.ID, "./imgs/okeyletsgo.jpg", tbot.OptCaption("Привет, это бот-заметочник. Я могу быть твоим помощником в быстром запоминании какой-либо информации. \n У меня будут храниться твои записи, которые ты в любой момент может просмотреть. \nПросто создай свою заметку и мы начнем! "), tbot.OptInlineKeyboardMarkup(buttons))
	if err != nil {
		fmt.Println(err)
	}
	}

func (r *Router) CallbackHandler(cq *tbot.CallbackQuery) {
	r.Client.AnswerCallbackQuery(cq.ID, tbot.OptText("Хороший выбор, мой дружок"))
	fmt.Println(cq.Data)
	msg := new(tbot.Message)
	msg = cq.Message
	if strings.Contains(cq.Data, "/note") {
		//msg1 := new(tbot.Message)
		msg.Text = cq.Data
		r.ReadNoteByName(msg)
	}
	switch cq.Data {
	case "/create":
		r.CreateNoteText(msg)
	case "/readall":
		r.ReadAllNotes(msg)
	}
}

func makeButtons() *tbot.InlineKeyboardMarkup {
	button1 := tbot.InlineKeyboardButton{
		Text:         fmt.Sprintf("Создать"),
		CallbackData: "/create",
	}
	button3 := tbot.InlineKeyboardButton{
		Text:         fmt.Sprintf("Прочесть"),
		CallbackData: "/readall",
	}
	return &tbot.InlineKeyboardMarkup{
		InlineKeyboard: [][]tbot.InlineKeyboardButton{
			[]tbot.InlineKeyboardButton{button1, button3},
		},
	}
}
func (r *Router) makeButtonsOfNotes(chatID string) *tbot.InlineKeyboardMarkup {
	var buttons = new([]tbot.InlineKeyboardButton)
	for _, note := range r.storage[chatID] {
		button := tbot.InlineKeyboardButton{
			Text:         fmt.Sprintf(note.Name),
			CallbackData: fmt.Sprintf("/note %s/////%s", note.Name, note.Class),
		}
		*buttons = append(*buttons, button)
	}
	return &tbot.InlineKeyboardMarkup{InlineKeyboard: [][]tbot.InlineKeyboardButton{
		*buttons,
	},
	}

}

func (r *Router) СreateNote(message *tbot.Message) {
	err := r.Client.SendChatAction(message.Chat.ID, tbot.ActionTyping)
	if err != nil {
		return
	}
	time.Sleep(time.Millisecond)
	text := message.Text
	reg := regexp.MustCompile(`Тип:(.+) Имя:(.+) Текст:(.+)`)
	arr := reg.FindStringSubmatch(text)
	if len(arr) != 4 {
		r.Client.SendMessage(message.Chat.ID, fmt.Sprintf("Не понимаю, о чем ты мне пытаешься сказать)"))
		return
	}
	note := notes.NewNote(strings.TrimPrefix(arr[1], " "), strings.TrimPrefix(arr[2], " "), strings.TrimPrefix(arr[3], " "))
	hb := sha1.Sum([]byte(note.Text))
	var sb []byte
	for _, i2 := range hb {
		sb = append(sb, i2)

	}
	note.Hash = string(sb)
	r.storage[message.Chat.ID] = append(r.storage[message.Chat.ID], *note)
	msg := fmt.Sprintf("Всё получилось, теперь у тебя есть заметка %s", note.Name)
	r.Client.SendMessage(message.Chat.ID, msg)
}
