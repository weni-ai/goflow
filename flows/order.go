package flows

import (
	"encoding/json"

	"github.com/nyaruka/gocommon/jsonx"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/utils"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

type Order struct {
	CatalogID    string        `json:"catalog_id,omitempty"`
	ProductItems []ProductItem `json:"product_items,omitempty"`
	Text         string        `json:"text,omitempty"`
}

type orderEnvelope struct {
	CatalogID    string        `json:"catalog_id,omitempty"`
	ProductItems []ProductItem `json:"product_items,omitempty"`
	Text         string        `json:"text,omitempty"`
}

type ProductItem struct {
	Currency          string          `json:"currency,omitempty"`
	ItemPrice         decimal.Decimal `json:"item_price,omitempty"`
	ProductRetailerID string          `json:"product_retailer_id"`
	Quantity          int64           `json:"quantity,omitempty"`
}

func (o *Order) Context(env envs.Environment) map[string]types.XValue {
	array := make([]types.XValue, len(o.ProductItems))
	for i, p := range o.ProductItems {
		array[i] = types.NewXObject(map[string]types.XValue{
			"currency":            types.NewXText(p.Currency),
			"item_price":          types.NewXNumber(p.ItemPrice),
			"product_retailer_id": types.NewXText(p.ProductRetailerID),
			"quantity":            types.NewXNumberFromInt64(p.Quantity),
		})
	}

	return map[string]types.XValue{
		"catalog_id":    types.NewXText(o.CatalogID),
		"product_items": types.NewXArray(array...),
		"text":          types.NewXText(o.Text),
	}
}

func ReadOrder(sa SessionAssets, data json.RawMessage, missing assets.MissingCallback) (*Order, error) {
	var envelope orderEnvelope
	var err error

	if err = utils.UnmarshalAndValidate(data, &envelope); err != nil {
		return nil, errors.Wrap(err, "unable to read order")
	}

	o := &Order{
		CatalogID:    envelope.CatalogID,
		ProductItems: envelope.ProductItems,
		Text:         envelope.Text,
	}

	return o, nil
}

func (o *Order) MarshalJSON() ([]byte, error) {
	oe := &orderEnvelope{
		CatalogID:    o.CatalogID,
		ProductItems: o.ProductItems,
		Text:         o.Text,
	}

	return jsonx.Marshal(oe)
}
