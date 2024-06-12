package flows

import (
	"encoding/json"

	"github.com/nyaruka/gocommon/jsonx"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/utils"
	"github.com/pkg/errors"
)

type NFMReply struct {
	Name         string                 `json:"name,omitempty"`
	ResponseJSON map[string]interface{} `json:"response_json,omitempty"`
}

type nfmReplyEnvelope struct {
	Name         string                 `json:"name,omitempty"`
	ResponseJSON map[string]interface{} `json:"response_json,omitempty"`
}

func (n *NFMReply) Context(env envs.Environment) map[string]types.XValue {
	jsonData, _ := json.Marshal(n.ResponseJSON)
	responseXValue := types.JSONToXValue(jsonData)

	return map[string]types.XValue{
		"name":          types.NewXText(n.Name),
		"response_json": responseXValue,
	}
}

func ReadNFMReply(sa SessionAssets, data json.RawMessage, missing assets.MissingCallback) (*NFMReply, error) {
	var envelope nfmReplyEnvelope
	var err error

	if err = utils.UnmarshalAndValidate(data, &envelope); err != nil {
		return nil, errors.Wrap(err, "unable to read nfm reply")
	}

	n := &NFMReply{
		Name:         envelope.Name,
		ResponseJSON: envelope.ResponseJSON,
	}

	return n, nil
}

func (n *NFMReply) MarshalJSON() ([]byte, error) {
	ne := &nfmReplyEnvelope{
		Name:         n.Name,
		ResponseJSON: n.ResponseJSON,
	}

	return jsonx.Marshal(ne)
}
