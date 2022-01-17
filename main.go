package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Claryslane/queue/internal/channel"
	"github.com/Claryslane/queue/pkg/util"
	Sex "github.com/Plankiton/SexPistol"
)

var channels = MakeChannels()

func main() {
	sex := Sex.NewPistol().
		Add("/channel/{channel}", func(r Sex.Request) (Sex.Json, int) {
			chanName := r.PathVars["channel"]

			chanKey := util.Sha1(chanName + time.Now().Format(time.RFC3339))

			if _, ok := channels[chanName]; !ok {
				channels[chanName] = channel.New(chanKey)
			} else {
				chanName, chanKey, ok = r.BasicAuth()
				if !ok {
					return Sex.Dict{
						"error": "You need be logged to set this channel key",
					}, Sex.StatusUnauthorized
				}

				b := []byte{}
				r.RawBody(&b)

				if len(b) < 8 {
					return Sex.Dict{
						"error": "The channel key need more than 8 characters",
					}, Sex.StatusBadRequest
				}

				channels[chanName] = channel.New(string(b))
			}

			return Sex.Dict{
				"chan_name": chanName,
				"cha_key":   chanKey,
			}, Sex.StatusOK
		}, "post").
		Add("/{channel}", func(r Sex.Request) (Sex.Json, int) {
			ctx, _ := context.WithTimeout(context.Background(), 2*time.Hour)

			nick := r.URL.Query().Get("nick")
			chanName, chanKey, ok := r.BasicAuth()
			if !ok {
				return Sex.Dict{
					"error": "You need a basic auth with 'channel' as user and 'channel key' as password",
				}, Sex.StatusUnauthorized
			}

			ch := channels[chanName]
			if ok := ch.Verify(chanKey); !ok {
				return Sex.Dict{
					"error": "Channel name or channel key is wrong",
				}, Sex.StatusUnauthorized
			}

			if ch.In == nil {
				ch.In = map[string]*chan channel.Message{}
			}

			if _, ok := ch.In[nick]; !ok {
				ch.In[nick] = new(chan channel.Message)
				go SendAll(ctx, ch.In[nick], chanName)
			}

			var message channel.Message
			if err := json.NewDecoder(r.Body).Decode(&message); err != nil {
				return Sex.Dict{
					"error": fmt.Sprintf("Error parsing json: %v", err),
				}, 400
			}

			*ch.In[nick] <- message

			return nil, 200
		}, "post")

	sex.Run()
}
