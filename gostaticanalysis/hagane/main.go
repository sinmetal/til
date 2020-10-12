package hagane

//go:generate hagane -template template.go.tmpl -o getter_setter.go -data '{"type":"Hoge"}' main.go
type Hoge struct {
	Mome string
	meme int
}

func (h *Hoge) Mogeta() {

}
