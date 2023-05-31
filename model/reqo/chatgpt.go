package reqo

// Chat POST "/chat" request object
type PostQuestion struct {
	Content string `json:"content"`
}
