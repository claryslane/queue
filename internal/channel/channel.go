package channel

type Message struct {
	Type    string `json:"type,omiempty"`
	From    string `json:"from,omiempty"`
	To      string `json:"to,omiempty"`
	Message string `json:"msg,omiempty"`
	Key     string `json:"Key,omiempty"`
}

type Channel struct {
	Key string
	In  map[string]*chan Message
	Out map[string]*chan Message
}

func New(key string) Channel {
	return Channel{
		Key: key,
	}
}

func (c Channel) Verify(key string) bool {
	return key == c.Key
}
