package router

import (
	"TelegramBot/types/notes"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/yanzay/tbot/v2"
	"os"
	"regexp"
	"strings"
	"time"
)

type Storage map[string][]notes.Note //map[идентификатор клиента]map[хэш]заметка
type Router struct {
	storage Storage //
	Client *tbot.Client
	NowNote notes.Note
	votings map[string]*voting
}
type voting struct {
	ups   int
	downs int
}


func NewRouter(server *tbot.Server) *Router {
	var err error
	defer func() {
		err = errors.Wrap(err,"NewRouter")
	}()
	router := &Router{}
	router.storage = make(Storage)
	router.Client = server.Client()
	router.votings = make(map[string]*voting)
	return router
}

func (r Router)CreateNote(message *tbot.Message) {
	err :=r.Client.SendChatAction(message.Chat.ID,tbot.ActionTyping)
	if err != nil {
		return
	}
	time.Sleep(time.Millisecond)
	text := message.Text
	reg:= regexp.MustCompile(`/create-note Class:(.+) Name:(.+) Text:(.+)`)
	arr := reg.FindStringSubmatch(text)
	if len(arr) != 4 {
		r.Client.SendMessage(message.Chat.ID,fmt.Sprintf("Failed creating note \n" +
			" Pattern for creation: /create-note " +
			"Class: <your_name_of_class> " +
			"Name:<your_name_of_note> " +
			"Text:<yout_text_of_note>"))
		return
	}
	r.Client.SendMessage(message.Chat.ID,"Сработало, супер")
	note := notes.NewNote(strings.TrimPrefix(arr[1]," "),strings.TrimPrefix(arr[2], " "),strings.TrimPrefix(arr[3]," "))
	hb := sha1.Sum([]byte(note.Text))
	var sb []byte
	for _, i2 := range hb {
		sb = append(sb, i2)

	}
	note.Hash = string(sb)
	r.storage[message.Chat.ID] = append(r.storage[message.Chat.ID],*note)
	msg := fmt.Sprintf("Successfully Creation of a note with name: %s", note.Name)
	fmt.Println(message.Chat.ID)
	r.Client.SendMessage(message.Chat.ID,msg)
}


func (r *Router) ReadNoteByName(message *tbot.Message) {
	err :=r.Client.SendChatAction(message.Chat.ID,tbot.ActionTyping)
	if err != nil {
		return
	}
	time.Sleep(time.Millisecond)
	text := message.Text
	reg:= regexp.MustCompile(`/read-note Name:(.+)`)
	arr := reg.FindStringSubmatch(text)
	if len(arr) != 2 {
		r.Client.SendMessage(message.Chat.ID,fmt.Sprintf("Failed reading Note: \n" +
			"patter for reading: " +
			"/read-note Name:<name_of_your_note> \n" +
			"if you forget name, you can get a list of all your notes by /read-all"))
		return
	}
	for _, note := range r.storage[message.Chat.ID] {
		if note.Name == strings.TrimPrefix(arr[1]," ") {
			r.Client.SendMessage(message.Chat.ID,fmt.Sprintf("Class: %s \n Name: %s \n Text: %s", note.Class,note.Name, note.Text))
			return
		}
	}
	msg := "Sorry, I can not find your note:("
	r.Client.SendMessage(message.Chat.ID,msg)
}

// ReadAllNotes функция получения всех записей
func (r *Router)ReadAllNotes(message *tbot.Message) {
	err :=r.Client.SendChatAction(message.Chat.ID,tbot.ActionTyping)
	if err != nil {
		return
	}
	time.Sleep(time.Millisecond)
	var list string
	if len(r.storage[message.Chat.ID]) == 0 {
		r.Client.SendMessage(message.Chat.ID,"Ты меня за кого держишь? Нет заметок, что ищешь?:)")

	}
	for i, note := range r.storage[message.Chat.ID] {
		list = fmt.Sprintf("%s%d. Class:%s, Name:%s \n",list,i+1,note.Class,note.Name)
	}
	r.Client.SendMessage(message.Chat.ID,list)
}

// WriteToJson функция записи в json
func (r *Router)WriteToJson(){
	var err error
	defer func() {
		err = errors.Wrap(err,"WriteToJson")
	}()
	file , err := os.Open("./storage.json")
	if err != nil {
		fmt.Println(err)
		return
	}
	b, err := json.Marshal(r.storage)
	_,err = file.Write(b)
	if err != nil {
		return
	}
	file.Close()
	return
}
func (r *Router)Help(message *tbot.Message) {
	err :=r.Client.SendChatAction(message.Chat.ID,tbot.ActionTyping)
	if err != nil {
		r.Client.SendMessage(message.Chat.ID,"Some error, sorry:(")
		return
	}
	time.Sleep(time.Millisecond)
	commands,err := r.Client.GetMyCommands()
	if err != nil {
		r.Client.SendMessage(message.Chat.ID,"Some error, sorry:(")
		return
	}
	msg := ""
	for _, command := range *commands {
		msg = fmt.Sprintf("%sКоманда \"%s\" \nОписание: %s\n\n",msg,command.Command,command.Description)

	}
	buttons := makeButtons()
	r.Client.SendMessage(message.Chat.ID,msg,tbot.OptInlineKeyboardMarkup(buttons))
	return
}

func (r *Router) StartMsg(m *tbot.Message) {
	buttons := makeButtons()
	r.Client.SendMessage(m.Chat.ID, "Привет, это бот-заметочник-дедлайнчик", tbot.OptInlineKeyboardMarkup(buttons))
}

func (r *Router) CallbackHandler(cq *tbot.CallbackQuery) {
	//votingID := fmt.Sprintf("%s:%d", cq.Message.Chat.ID, cq.Message.MessageID)
	//vtng := new(voting)
	r.Client.AnswerCallbackQuery(cq.ID, tbot.OptText("Хороший выбор, мой дружок"))
	fmt.Println(cq.Data)
	msg := new(tbot.Message)
	msg = cq.Message
	switch cq.Data {
	case "/create":
		r.CreateNote(msg)
	case "/read":
		r.ReadNoteByName(msg)
	case "/readall":
		r.ReadAllNotes(msg)
	}
}

func makeButtons() *tbot.InlineKeyboardMarkup {
	button1 := tbot.InlineKeyboardButton{
		Text:         fmt.Sprintf("Создать запись"),
		CallbackData: "/create",
	}
	button2 := tbot.InlineKeyboardButton{
		Text:         fmt.Sprintf("Прочесть запись"),
		CallbackData: "/read",
	}
	button3 := tbot.InlineKeyboardButton{
		Text:         fmt.Sprintf("Вывести список всех записей"),
		CallbackData: "/readall",
	}
	return &tbot.InlineKeyboardMarkup{
		InlineKeyboard: [][]tbot.InlineKeyboardButton{
			[]tbot.InlineKeyboardButton{button1, button2, button3},
		},
	}
}


