package main

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/bobg/scp"
	"github.com/chain/txvm/crypto/ed25519"
	"github.com/chain/txvm/protocol/bc"
)

type (
	marshaled struct {
		M json.RawMessage
		S string // hex-encoded signature over marshaledPayload
	}

	marshaledPayload struct {
		C int32
		V string
		I int
		Q scp.QSet
		T marshaledTopic
	}

	marshaledTopic struct {
		Type        int // scp.Phase values
		X, Y        []bc.Hash
		B, C, P, PP marshaledBallot
		PN, HN, CN  int
	}

	marshaledBallot struct {
		N int
		X bc.Hash
	}
)

func marshal(msg *scp.Msg) ([]byte, error) {
	var mt marshaledTopic
	switch topic := msg.T.(type) {
	case *scp.NomTopic:
		mt.Type = int(scp.PhNom)

		var x, y []bc.Hash
		for _, val := range topic.X {
			if val != nil {
				x = append(x, valToHash(val))
			}
		}
		for _, val := range topic.Y {
			if val != nil {
				y = append(y, valToHash(val))
			}
		}
		mt.X = x
		mt.Y = y

	case *scp.PrepTopic:
		mt.Type = int(scp.PhPrep)

		mt.B = marshaledBallot{N: topic.B.N, X: valToHash(topic.B.X)}
		mt.P = marshaledBallot{N: topic.P.N, X: valToHash(topic.P.X)}
		mt.PP = marshaledBallot{N: topic.PP.N, X: valToHash(topic.PP.X)}
		mt.HN = topic.HN
		mt.CN = topic.CN

	case *scp.CommitTopic:
		mt.Type = int(scp.PhCommit)

		mt.B = marshaledBallot{N: topic.B.N, X: valToHash(topic.B.X)}
		mt.PN = topic.PN
		mt.HN = topic.HN
		mt.CN = topic.CN

	case *scp.ExtTopic:
		mt.Type = int(scp.PhExt)

		mt.C = marshaledBallot{N: topic.C.N, X: valToHash(topic.C.X)}
		mt.HN = topic.HN
	}
	mp := marshaledPayload{
		C: msg.C,
		V: string(msg.V),
		I: int(msg.I),
		Q: msg.Q,
		T: mt,
	}
	mpbytes, err := json.Marshal(mp) // xxx json is subject to mutation in transit!
	if err != nil {
		return nil, err
	}
	sig := ed25519.Sign(prv, mpbytes)
	m := marshaled{
		M: mpbytes,
		S: hex.EncodeToString(sig),
	}
	return json.Marshal(m)
}

func valToHash(v scp.Value) (result bc.Hash) {
	if v != nil {
		result = bc.Hash(v.(valtype))
	}
	return result
}

func unmarshalBallot(mb marshaledBallot) scp.Ballot {
	return scp.Ballot{
		N: mb.N,
		X: valtype(mb.X),
	}
}

func unmarshal(b []byte) (*scp.Msg, error) {
	var m marshaled
	err := json.Unmarshal(b, &m)
	if err != nil {
		return nil, err
	}

	var mp marshaledPayload
	err = json.Unmarshal(m.M, &mp)
	if err != nil {
		return nil, err
	}

	sig, err := hex.DecodeString(m.S)
	if err != nil {
		return nil, err
	}

	u, err := url.Parse(mp.V)
	if err != nil {
		return nil, err
	}
	pubkeyHex := u.Path
	pubkeyHex = strings.Trim(pubkeyHex, "/")
	pubkey, err := hex.DecodeString(pubkeyHex)
	if err != nil {
		return nil, err
	}
	if !ed25519.Verify(pubkey, m.M, sig) {
		return nil, errors.New("bad signature")
	}

	var topic scp.Topic
	switch scp.Phase(mp.T.Type) {
	case scp.PhNom:
		var x, y scp.ValueSet
		for _, v := range mp.T.X {
			x = append(x, valtype(v))
		}
		for _, v := range mp.T.Y {
			y = append(y, valtype(v))
		}
		topic = &scp.NomTopic{
			X: x,
			Y: y,
		}

	case scp.PhPrep:
		topic = &scp.PrepTopic{
			B:  unmarshalBallot(mp.T.B),
			P:  unmarshalBallot(mp.T.P),
			PP: unmarshalBallot(mp.T.PP),
			HN: mp.T.HN,
			CN: mp.T.CN,
		}

	case scp.PhCommit:
		topic = &scp.CommitTopic{
			B:  unmarshalBallot(mp.T.B),
			PN: mp.T.PN,
			HN: mp.T.HN,
			CN: mp.T.CN,
		}

	case scp.PhExt:
		topic = &scp.ExtTopic{
			C:  unmarshalBallot(mp.T.C),
			HN: mp.T.HN,
		}

	default:
		return nil, fmt.Errorf("unknown topic type %d", mp.T.Type)
	}

	msg := &scp.Msg{
		C: mp.C,
		V: scp.NodeID(mp.V),
		I: scp.SlotID(mp.I),
		Q: mp.Q,
		T: topic,
	}
	return msg, nil
}
