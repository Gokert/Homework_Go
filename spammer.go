package main

import (
	"fmt"
	"sort"
	"sync"
)

func RunPipeline(cmds ...cmd) {
	var out, in chan interface{}
	in = make(chan interface{})
	wg := &sync.WaitGroup{}

	for _, c := range cmds {
		out = make(chan interface{})
		wg.Add(1)
		go func(c cmd, in, out chan interface{}) {
			defer wg.Done()
			c(in, out)
			close(out)
		}(c, in, out)
		in = out
	}
	wg.Wait()
}

func SelectUsers(in, out chan interface{}) {
	var seen sync.Map
	var wg sync.WaitGroup

	for email := range in {
		wg.Add(1)
		go func(email interface{}) {
			defer wg.Done()

			emailStr := email.(string)

			alias := usersAliases[emailStr]

			if alias != "" {
				emailStr = alias
			}

			_, loaded1 := seen.LoadOrStore(emailStr, true)
			if !loaded1 {
				user := GetUser(emailStr)
				out <- user
				return
			} else {

				GetUser("")
			}
		}(email)
	}
	wg.Wait()
}

func SelectMessages(in, out chan interface{}) {
	var wg sync.WaitGroup

	batchSize := 2

	users := make([]User, 0, batchSize)

	for user := range in {
		users = append(users, user.(User))

		if len(users) == batchSize {
			wg.Add(1)
			go func(users []User) {
				defer wg.Done()

				bufs := make([]User, 0, batchSize)

				for _, u := range users {
					bufs = append(bufs, u)
				}

				messages, _ := GetMessages(bufs...)

				for _, msg := range messages {
					out <- msg
				}
			}(users)

			users = make([]User, 0, batchSize)
		}
	}

	if len(users) == 1 {
		messages, _ := GetMessages(users[0])

		for _, msg := range messages {
			out <- msg
		}
	}

	wg.Wait()
}

func CheckSpam(in, out chan interface{}) {
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, HasSpamMaxAsyncRequests)

	for msg := range in {
		wg.Add(1)
		semaphore <- struct{}{}

		go func(msg interface{}) {
			defer wg.Done()

			msgID, ok := msg.(MsgID)
			if !ok {
				return
			}

			hasSpam, _ := HasSpam(msgID)
			out <- MsgData{
				ID:      msgID,
				HasSpam: hasSpam,
			}
			<-semaphore
		}(msg)
	}

	wg.Wait()
}

func CombineResults(in, out chan interface{}) {
	var spamResults []MsgData
	var nonSpamResults []MsgData

	for data := range in {
		msgData, ok := data.(MsgData)
		if !ok {
			continue
		}
		if msgData.HasSpam {
			spamResults = append(spamResults, msgData)
		} else {
			nonSpamResults = append(nonSpamResults, msgData)
		}
	}

	sort.Slice(spamResults, func(i, j int) bool {
		return spamResults[i].ID < spamResults[j].ID
	})

	sort.Slice(nonSpamResults, func(i, j int) bool {
		return nonSpamResults[i].ID < nonSpamResults[j].ID
	})

	for _, result := range spamResults {
		out <- fmt.Sprintf("%v %v", result.HasSpam, result.ID)
	}

	for _, result := range nonSpamResults {
		out <- fmt.Sprintf("%v %v", result.HasSpam, result.ID)
	}

}
