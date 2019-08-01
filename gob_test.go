// Encoding and decoding on gob objects.

package main

import "testing"

func TestGobString(t *testing.T) {
	const test = "Hello World ❤"

	b, err := encodeString(test)
	if err != nil {
		t.Errorf("Couldn't encode a string to gob: %s.", err.Error())
	}

	decoded, err := decodeString(b)
	if err != nil {
		t.Errorf("Couldn't decode a string from gob: %s.", err.Error())
	}

	if decoded != test {
		t.Errorf("Gob+String was incorrect, got: %s, want: %s.", decoded, test)
	}
}

func TestGobRepo(t *testing.T) {
	test := Repo{
		UUID: "feaaade6-43bb-4031-b388-731ef80f63d7",
		URL:  "git@github.com:Permaweb/Host.git",
		IPFS: "bafybeieceh6aoqp35neajnm23z5kjctgtaujmiprlbpxmoytngd3jys5am",
		Key:  "12D3KooWGTgRDFis5f1yByfwvU6iUyndVHMWQKXmdU3kjaDKHkwt",
		IPNS: "12D3KooWGTgRDFis5f1yByfwvU6iUyndVHMWQKXmdU3kjaDKHkwt",
	}

	b, err := test.encode()
	if err != nil {
		t.Errorf("Couldn't encode a repo to gob: %s.", err.Error())
	}

	decoded, err := decodeRepo(b)
	if err != nil {
		t.Errorf("Couldn't decode a repo from gob: %s.", err.Error())
	}

	// UUID
	if decoded.UUID != test.UUID {
		t.Errorf("Gob+Repo was incorrect, got: %s, want: %s.", decoded.UUID, test.UUID)
	}

	// URL
	if decoded.URL != test.URL {
		t.Errorf("Gob+Repo was incorrect, got: %s, want: %s.", decoded.URL, test.URL)
	}

	// IPFS
	if decoded.IPFS != test.IPFS {
		t.Errorf("Gob+Repo was incorrect, got: %s, want: %s.", decoded.IPFS, test.IPFS)
	}

	// Key
	if decoded.Key != test.Key {
		t.Errorf("Gob+Repo was incorrect, got: %s, want: %s.", decoded.Key, test.Key)
	}

	// IPNS
	if decoded.IPNS != test.IPNS {
		t.Errorf("Gob+Repo was incorrect, got: %s, want: %s.", decoded.IPNS, test.IPNS)
	}
}
