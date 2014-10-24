package main

import (
	"fmt"
	"time"

	. "github.com/glycerine/HyperDex/bindings/go"
	. "github.com/glycerine/HyperDex/bindings/go/client"
)

func putStuff(client *Client, key string) {
	attrs := Attributes{
		"name":          "john",
		"height":        float64(241.12421),
		"profile_views": int64(6075551024),
		"pending_requests": List{
			"haha",
			"hehe",
		},
		"ratings": List{
			1.22141,
			-5235.241,
			92804.14,
		},
		"hobbies": Set{
			"qowue",
			"waoihdwao",
		},
		"ages": Set{
			-41,
			284,
			2304,
		},
		"unread_messages": Map{
			"oahd":      "waohdaw",
			"wapodajwp": "waohdwoqpd",
		},
		"upvotes": Map{
			"wadwa": 10294,
			"aowid": 98571,
		},
	}

	err := client.Put("profiles", key, attrs)

	if err != nil {
		panic(err)
	}
}

func setupStuff() {
	admin, err := NewAdmin("127.0.0.1", 1982)
	if err != nil {
		panic(err)
	}

	err = admin.RemoveSpace(`profiles`)
	if err != nil {
		if err.Error() != "Error 8777: cannot rm space: does not exist" &&
			err.Error() != "unknown hyperdex_client_returncode: cannot rm space: does not exist" {
			panic(err)
		}
	}

	err = admin.AddSpace(`space profiles
key username
attributes
    string name,
    float height,
    int profile_views,
    list(string) pending_requests,
    list(float) ratings,
    set(string) hobbies,
    set(int) ages,
    map(string, string) unread_messages,
    map(string, int) upvotes
subspace name
subspace height
subspace profile_views
`)
	if err != nil {
		panic(err)
	}
}

func write() {
	client, err, _ := NewClient("127.0.0.1", 1982)
	if err != nil {
		panic(err)
	}
	defer client.Destroy()
	t0 := time.Now()
	fmt.Printf("starting writes at %v\n", t0)
	for i := 0; i < 10000; i++ {
		putStuff(client, fmt.Sprintf("%d", i))
	}
	fmt.Printf("elapsed time: %v\n", time.Since(t0))
}

func main() {
	for i := 0; i < 100; i++ {
		setupStuff()
		for k := 0; k < 10; k++ {
			write()
		}
	}

	// with 1 daemon,
	// 6.25469113s for 10k writes => 1600 writes/sec.

	// with 3 daemons: 3x slower
	// results so far, trying to write or update 10k times.
	// 1st run: 10k writes in 15 seconds => 666.67 writes/second
	// 2nd run interrupted with:
	// panic: HYPERDEX_CLIENT_RECONFIGURE: reconfiguration affecting virtual_server(13995)/server(6235075844444110843)
	// 3rd run: appears to hang forever
	// 4th run: hangs forever
	// 5th run: reduce write count to 10: still hangs forever.
}
