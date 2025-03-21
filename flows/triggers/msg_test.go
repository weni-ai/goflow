package triggers

import (
	"encoding/json"
	"log"
	"testing"

	"github.com/nyaruka/gocommon/jsonx"
)

func TestMsgTrigger(t *testing.T) {

	var msgTrigger MsgTrigger

	err := json.Unmarshal([]byte(triggerJSON), &msgTrigger)
	if err != nil {
		t.Fatal(err)
	}
	tm, err := msgTrigger.MarshalJSON()
	if err != nil {
		t.Fatal(err)
	}

	log.Println(string(tm))
}

func TestMsgTriggerMarshallJSON(t *testing.T) {
	var mtEnvelop msgTriggerEnvelope

	err := json.Unmarshal([]byte(triggerJSON), &mtEnvelop)
	if err != nil {
		t.Fatal(err)
	}

	res, err := jsonx.Marshal(mtEnvelop)
	if err != nil {
		t.Fatal(err)
	}
	log.Println(string(res))
}

const triggerJSON = `{
	"type": "msg",
	"environment": {
			"date_format": "DD-MM-YYYY",
			"time_format": "tt:mm",
			"timezone": "America/Argentina/Buenos_Aires",
			"number_format": {
					"decimal_symbol": ".",
					"digit_grouping_symbol": ","
			},
			"default_country": "BR",
			"redaction_policy": "none",
			"max_value_length": 3840
	},
	"flow": {
			"uuid": "7a2236bc-4c92-4ad6-b81a-0d86d2fdedc0",
			"name": "catalogo"
	},
	"contact": {
			"uuid": "7eed1ae9-4f7d-4211-a788-8d73ebdd518c",
			"id": 389,
			"name": "Roberta Moreira",
			"status": "active",
			"timezone": "America/Argentina/Buenos_Aires",
			"created_on": "2023-10-05T19:35:25.351816Z",
			"urns": [
					"whatsapp:5582999489287?channel=3065fa26-593b-4517-9318-6050165c78d7\\u0026id=387\\u0026priority=1000"
			]
	},
	"params": {
			"order": {
				"catalog_id": "1729754620797120",
				"product_items": [
						{
								"currency": "BRL",
								"item_price": 2.46,
								"product_retailer_id": "1",
								"quantity": 1
						}
				],
				"text": ""
			},
			"nfm_reply": {
				"name": "Flow Wpp",
				"response_json": "{\"flow_token\": \"<FLOW_TOKEN>\", \"optional_param1\": \"<value1>\", \"optional_param2\": \"<value2>\"}"
			},
			"ig_comment": {
				"text": "hello",
				"from": {
					"id": "1234567890",
					"username": "bob"
			}
	},
	"triggered_on": "2023-10-05T19:35:25.576204487Z",
	"msg": {
			"uuid": "22454ba0-b8db-440d-a1a6-1a6bb56e7add",
			"id": 162145,
			"urn": "whatsapp:5582999489287",
			"channel": {
					"uuid": "3065fa26-593b-4517-9318-6050165c78d7",
					"name": "Teste Weni Cloud 5"
			},
			"text": "",
			"external_id": "wamid.HB1gMNTU4ODkzaaN2asaTasasYda1fa1asaOsaTA21FsQaa3I1AaaqEa2hssgWa23M0V2aCMEafY4NDYaswNDaRGOaUFBNzcwNjQysRgA=",
			"order": {
					"catalog_id": "1729754620797120",
					"product_items": [
							{
									"currency": "BRL",
									"item_price": 2.46,
									"product_retailer_id": "1",
									"quantity": 1
							}
					],
					"text": ""
			}
	}
}`
