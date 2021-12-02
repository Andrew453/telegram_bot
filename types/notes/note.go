package notes

import "time"

type Note struct {
	Hash string
	Class string `json:"class"`
	Name string `json:"name"`
	Text string `json:"text"`
	Deadline time.Time `json:"deadline"`
}

func NewNote(class string, name string, text string) *Note {
	return &Note{Class: class, Name: name, Text: text}
}
func (n *Note) GetClass(){

}


