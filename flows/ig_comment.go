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

type IGComment struct {
	Text string `json:"text,omitempty"`
	From struct {
		ID       string `json:"id,omitempty"`
		Username string `json:"username,omitempty"`
	} `json:"from,omitempty"`
	Media struct {
		AdID             string `json:"ad_id,omitempty"`
		ID               string `json:"id,omitempty"`
		MediaProductType string `json:"media_product_type,omitempty"`
		OriginalMediaID  string `json:"original_media_id,omitempty"`
	} `json:"media,omitempty"`
	Time int64  `json:"time,omitempty"`
	ID   string `json:"id,omitempty"`
}

type igCommentEnvelope struct {
	Text string `json:"text,omitempty"`
	From struct {
		ID       string `json:"id,omitempty"`
		Username string `json:"username,omitempty"`
	} `json:"from,omitempty"`
	Media struct {
		AdID             string `json:"ad_id,omitempty"`
		ID               string `json:"id,omitempty"`
		MediaProductType string `json:"media_product_type,omitempty"`
		OriginalMediaID  string `json:"original_media_id,omitempty"`
	} `json:"media,omitempty"`
	Time int64  `json:"time,omitempty"`
	ID   string `json:"id,omitempty"`
}

func (i *IGComment) Context(env envs.Environment) map[string]types.XValue {
	fromMap := map[string]types.XValue{
		"id":       types.NewXText(i.From.ID),
		"username": types.NewXText(i.From.Username),
	}

	mediaMap := map[string]types.XValue{
		"ad_id":              types.NewXText(i.Media.AdID),
		"id":                 types.NewXText(i.Media.ID),
		"media_product_type": types.NewXText(i.Media.MediaProductType),
		"original_media_id":  types.NewXText(i.Media.OriginalMediaID),
	}

	return map[string]types.XValue{
		"text":  types.NewXText(i.Text),
		"from":  types.NewXObject(fromMap),
		"media": types.NewXObject(mediaMap),
		"time":  types.NewXNumberFromInt64(i.Time),
		"id":    types.NewXText(i.ID),
	}
}

func ReadIGComment(sa SessionAssets, data json.RawMessage, missing assets.MissingCallback) (*IGComment, error) {
	var comment IGComment
	var err error

	if err = utils.UnmarshalAndValidate(data, &comment); err != nil {
		return nil, errors.Wrap(err, "unable to read ig comment")
	}

	return &comment, nil
}

func (i *IGComment) MarshalJSON() ([]byte, error) {
	ie := &igCommentEnvelope{
		Text:  i.Text,
		From:  i.From,
		Media: i.Media,
		Time:  i.Time,
		ID:    i.ID,
	}

	return jsonx.Marshal(ie)
}
