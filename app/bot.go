package app

import (
	"TelegramBot/config"
	"TelegramBot/router"
	"github.com/pkg/errors"
	"github.com/yanzay/tbot/v2"
)

type Bot struct {
	config config.Configuration
	router *router.Router
	//client *tbot.Client
	server *tbot.Server
}

func NewBot() *Bot {
	return &Bot{}
}

func (b *Bot) Start () (err error) {
	defer func() {
		err = errors.Wrap(err, "Bot Start")
	}()
	err = b.manage()
	if err != nil {
		return err
	}
	b.server = tbot.New(b.config.Token)
	b.router = router.NewRouter(b.server)
	b.server.HandleMessage("/start", b.router.StartMsg)
	b.server.HandleCallback(b.router.CallbackHandler)
	b.server.HandleMessage("/help", b.router.Help)
	b.server.HandleMessage("/create *", b.router.CreateNote)
	b.server.HandleMessage("/readall", b.router.ReadAllNotes)
	b.server.HandleMessage("/read", b.router.ReadNoteByName)
	return b.server.Start()
}


func (b *Bot) Stop () {
	b.router.WriteToJson()
}

func(b *Bot) manage () (err error) {
	defer func() {
		err = errors.Wrap(err,"bot manage")
	}()
	err = config.LoadConfigHCL("./bot.hcl",&b.config)
	if err != nil {
		err =errors.Wrap(err,"load config")
		return err
	}
	return nil
}


// json.Marshall -> []byte -> string -> int ....




