package actions_test

import (
	"encoding/json"
	"io/ioutil"
	"testing"
	"time"

	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/actions"
	"github.com/nyaruka/goflow/flows/triggers"
	"github.com/nyaruka/goflow/test"
	"github.com/nyaruka/goflow/utils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var actionUUID = flows.ActionUUID("ad154980-7bf7-4ab8-8728-545fd6378912")

var actionTests = []struct {
	action flows.Action
	json   string
}{
	{
		actions.NewAddContactGroupsAction(
			actionUUID,
			[]*assets.GroupReference{
				assets.NewGroupReference(assets.GroupUUID("b7cf0d83-f1c9-411c-96fd-c511a4cfa86d"), "Testers"),
				assets.NewVariableGroupReference("@(format_location(contact.fields.state)) Members"),
			},
		),
		`{
			"type": "add_contact_groups",
			"uuid": "ad154980-7bf7-4ab8-8728-545fd6378912",
			"groups": [
				{
					"uuid": "b7cf0d83-f1c9-411c-96fd-c511a4cfa86d",
					"name": "Testers"
				},
				{
					"name_match": "@(format_location(contact.fields.state)) Members"
				}
			]
		}`,
	},
	{
		actions.NewAddContactURNAction(
			actionUUID,
			"tel",
			"+234532626677",
		),
		`{
			"type": "add_contact_urn",
			"uuid": "ad154980-7bf7-4ab8-8728-545fd6378912",
			"scheme": "tel",
			"path": "+234532626677"
		}`,
	},
	{
		actions.NewAddInputLabelsAction(
			actionUUID,
			[]*assets.LabelReference{
				assets.NewLabelReference(assets.LabelUUID("3f65d88a-95dc-4140-9451-943e94e06fea"), "Spam"),
				assets.NewVariableLabelReference("@(format_location(contact.fields.state)) Messages"),
			},
		),
		`{
			"type": "add_input_labels",
			"uuid": "ad154980-7bf7-4ab8-8728-545fd6378912",
			"labels": [
				{
					"uuid": "3f65d88a-95dc-4140-9451-943e94e06fea",
					"name": "Spam"
				},
				{
					"name_match": "@(format_location(contact.fields.state)) Messages"
				}
			]
		}`,
	},
	{
		actions.NewCallResthookAction(
			actionUUID,
			"new-registration",
			"My Result",
		),
		`{
			"type": "call_resthook",
			"uuid": "ad154980-7bf7-4ab8-8728-545fd6378912",
			"resthook": "new-registration",
			"result_name": "My Result"
		}`,
	},
	{
		actions.NewCallWebhookAction(
			actionUUID,
			"POST",
			"http://example.com/ping",
			map[string]string{
				"Authentication": "Token @contact.fields.token",
			},
			`{"contact_id": 234}`, // body
			"Webhook Response",
		),
		`{
			"type": "call_webhook",
			"uuid": "ad154980-7bf7-4ab8-8728-545fd6378912",
			"method": "POST",
			"url": "http://example.com/ping",
			"headers": {
				"Authentication": "Token @contact.fields.token"
			},
			"body": "{\"contact_id\": 234}",
			"result_name": "Webhook Response"
		}`,
	},
	{
		actions.NewRemoveContactGroupsAction(
			actionUUID,
			[]*assets.GroupReference{
				assets.NewGroupReference(assets.GroupUUID("b7cf0d83-f1c9-411c-96fd-c511a4cfa86d"), "Testers"),
				assets.NewVariableGroupReference("@(format_location(contact.fields.state)) Members"),
			},
			false,
		),
		`{
			"type": "remove_contact_groups",
			"uuid": "ad154980-7bf7-4ab8-8728-545fd6378912",
			"groups": [
				{
					"uuid": "b7cf0d83-f1c9-411c-96fd-c511a4cfa86d",
					"name": "Testers"
				},
				{
					"name_match": "@(format_location(contact.fields.state)) Members"
				}
			],
			"all_groups": false
		}`,
	},
	{
		actions.NewSendBroadcastAction(
			actionUUID,
			"Hi there",
			[]string{"http://example.com/red.jpg"},
			[]string{"Red", "Blue"},
			[]urns.URN{"twitter:nyaruka"},
			[]*flows.ContactReference{
				flows.NewContactReference(flows.ContactUUID("cbe87f5c-cda2-4f90-b5dd-0ac93a884950"), "Bob Smith"),
			},
			[]*assets.GroupReference{
				assets.NewGroupReference(assets.GroupUUID("b7cf0d83-f1c9-411c-96fd-c511a4cfa86d"), "Testers"),
			},
			nil,
		),
		`{
			"type": "send_broadcast",
			"uuid": "ad154980-7bf7-4ab8-8728-545fd6378912",
			"text": "Hi there",
			"attachments": ["http://example.com/red.jpg"],
			"quick_replies": ["Red", "Blue"],
			"urns": ["twitter:nyaruka"],
            "contacts": [
				{
					"uuid": "cbe87f5c-cda2-4f90-b5dd-0ac93a884950",
					"name": "Bob Smith"
				}
			],
			"groups": [
				{
					"uuid": "b7cf0d83-f1c9-411c-96fd-c511a4cfa86d",
					"name": "Testers"
				}
			]
		}`,
	},
	{
		actions.NewSendEmailAction(
			actionUUID,
			[]string{"bob@example.com"},
			"Hi there",
			"So I was thinking...",
		),
		`{
			"type": "send_email",
			"uuid": "ad154980-7bf7-4ab8-8728-545fd6378912",
			"addresses": ["bob@example.com"],
			"subject": "Hi there",
			"body": "So I was thinking..."
		}`,
	},
	{
		actions.NewSendMsgAction(
			actionUUID,
			"Hi there",
			[]string{"http://example.com/red.jpg"},
			[]string{"Red", "Blue"},
			true,
		),
		`{
			"type": "send_msg",
			"uuid": "ad154980-7bf7-4ab8-8728-545fd6378912",
			"text": "Hi there",
			"attachments": ["http://example.com/red.jpg"],
			"quick_replies": ["Red", "Blue"],
			"all_urns": true
		}`,
	},
	{
		actions.NewSetContactChannelAction(
			actionUUID,
			assets.NewChannelReference(assets.ChannelUUID("57f1078f-88aa-46f4-a59a-948a5739c03d"), "My Android Phone"),
		),
		`{
			"type": "set_contact_channel",
			"uuid": "ad154980-7bf7-4ab8-8728-545fd6378912",
			"channel": {
				"uuid": "57f1078f-88aa-46f4-a59a-948a5739c03d",
				"name": "My Android Phone"
			}
		}`,
	},
	{
		actions.NewSetContactFieldAction(
			actionUUID,
			assets.NewFieldReference("gender", "Gender"),
			"Male",
		),
		`{
			"type": "set_contact_field",
			"uuid": "ad154980-7bf7-4ab8-8728-545fd6378912",
			"field": {
				"key": "gender",
				"name": "Gender"
			},
			"value": "Male"
		}`,
	},
	{
		actions.NewSetContactLanguageAction(
			actionUUID,
			"eng",
		),
		`{
			"type": "set_contact_language",
			"uuid": "ad154980-7bf7-4ab8-8728-545fd6378912",
			"language": "eng"
		}`,
	},
	{
		actions.NewSetContactNameAction(
			actionUUID,
			"Bob",
		),
		`{
			"type": "set_contact_name",
			"uuid": "ad154980-7bf7-4ab8-8728-545fd6378912",
			"name": "Bob"
		}`,
	},
	{
		actions.NewSetContactTimezoneAction(
			actionUUID,
			"Africa/Kigali",
		),
		`{
			"type": "set_contact_timezone",
			"uuid": "ad154980-7bf7-4ab8-8728-545fd6378912",
			"timezone": "Africa/Kigali"
		}`,
	},
	{
		actions.NewSetRunResultAction(
			actionUUID,
			"Response 1",
			"yes",
			"Yes",
		),
		`{
			"type": "set_run_result",
			"uuid": "ad154980-7bf7-4ab8-8728-545fd6378912",
			"name": "Response 1",
			"value": "yes",
			"category": "Yes"
		}`,
	},
	{
		actions.NewStartFlowAction(
			actionUUID,
			assets.NewFlowReference(assets.FlowUUID("fece6eac-9127-4343-9269-56e88f391562"), "Parent"),
			true, // terminal
		),
		`{
			"type": "start_flow",
			"uuid": "ad154980-7bf7-4ab8-8728-545fd6378912",
			"flow": {
				"uuid": "fece6eac-9127-4343-9269-56e88f391562",
				"name": "Parent"
			},
			"terminal": true
		}`,
	},
	{
		actions.NewStartSessionAction(
			actionUUID,
			assets.NewFlowReference(assets.FlowUUID("fece6eac-9127-4343-9269-56e88f391562"), "Parent"),
			[]urns.URN{"twitter:nyaruka"},
			[]*flows.ContactReference{
				flows.NewContactReference(flows.ContactUUID("cbe87f5c-cda2-4f90-b5dd-0ac93a884950"), "Bob Smith"),
			},
			[]*assets.GroupReference{
				assets.NewGroupReference(assets.GroupUUID("b7cf0d83-f1c9-411c-96fd-c511a4cfa86d"), "Testers"),
			},
			nil,  // legacy vars
			true, // create new contact
		),
		`{
			"type": "start_session",
			"uuid": "ad154980-7bf7-4ab8-8728-545fd6378912",
			"flow": {
				"uuid": "fece6eac-9127-4343-9269-56e88f391562",
				"name": "Parent"
			},
			"urns": ["twitter:nyaruka"],
            "contacts": [
				{
					"uuid": "cbe87f5c-cda2-4f90-b5dd-0ac93a884950",
					"name": "Bob Smith"
				}
			],
			"groups": [
				{
					"uuid": "b7cf0d83-f1c9-411c-96fd-c511a4cfa86d",
					"name": "Testers"
				}
			],
			"create_contact": true
		}`,
	},
}

func TestMarshaling(t *testing.T) {
	session, err := test.CreateTestSession("", nil)
	require.NoError(t, err)

	for _, tc := range actionTests {
		// test validating the action
		err := tc.action.Validate(session.Assets(), flows.NewValidationContext())
		assert.NoError(t, err)

		// test marshaling the action
		actualJSON, err := json.Marshal(tc.action)
		assert.NoError(t, err)

		test.AssertEqualJSON(t, json.RawMessage(tc.json), actualJSON, "new action produced unexpected JSON")
	}
}

func TestValidation(t *testing.T) {
	session, err := test.CreateTestSession("", nil)
	require.NoError(t, err)

	errorFile, err := ioutil.ReadFile("testdata/validation.json")
	require.NoError(t, err)

	tests := []struct {
		ActionJSON json.RawMessage `json:"action"`
		ErrMsg     string          `json:"error"`
	}{}

	err = json.Unmarshal(errorFile, &tests)
	require.NoError(t, err)

	for _, tc := range tests {
		action, err := actions.ReadAction(tc.ActionJSON)
		require.NoError(t, err)

		err = action.Validate(session.Assets(), flows.NewValidationContext())
		assert.EqualError(t, err, tc.ErrMsg)
	}
}

var contactJSON = `{
	"uuid": "5d76d86b-3bb9-4d5a-b822-c9d86f5d8e4f",
	"name": "Ryan Lewis",
	"language": "eng",
	"timezone": "America/Guayaquil",
	"urns": [
		"tel:+12065551212?channel=57f1078f-88aa-46f4-a59a-948a5739c03d", 
		"twitterid:54784326227#nyaruka"
	],
	"groups": [
		{"uuid": "b7cf0d83-f1c9-411c-96fd-c511a4cfa86d", "name": "Testers"}
	],
	"fields": {
		"gender": {
			"text": "Male"
		}
	},
	"created_on": "2018-06-20T11:40:30.123456789-00:00"
}`

func TestExecution(t *testing.T) {
	assetsJSON, err := ioutil.ReadFile("testdata/assets.json")
	require.NoError(t, err)

	errorFile, err := ioutil.ReadFile("testdata/execution.json")
	require.NoError(t, err)

	tests := []struct {
		Description string            `json:"description"`
		ActionJSON  json.RawMessage   `json:"action"`
		Events      []json.RawMessage `json:"events"`
	}{}

	err = json.Unmarshal(errorFile, &tests)
	require.NoError(t, err)

	utils.SetTimeSource(utils.NewSequentialTimeSource(time.Date(2018, 10, 18, 14, 20, 30, 123456, time.UTC)))
	defer utils.SetTimeSource(utils.DefaultTimeSource)

	utils.SetUUIDGenerator(utils.NewSeededUUID4Generator(12345))
	defer utils.SetUUIDGenerator(utils.DefaultUUIDGenerator)

	for _, tc := range tests {
		// create unstarted session from our assets
		session, err := test.CreateSession(assetsJSON, "")
		require.NoError(t, err)

		// load our contact
		contact, err := flows.ReadContact(session.Assets(), json.RawMessage(contactJSON), true)
		require.NoError(t, err)

		// get the main flow
		flow, err := session.Assets().Flows().Get(assets.FlowUUID("bead76f5-dac4-4c9d-996c-c62b326e8c0a"))
		require.NoError(t, err)

		// add this tests action to its first node
		action, err := actions.ReadAction(tc.ActionJSON)
		require.NoError(t, err)
		flow.Nodes()[0].AddAction(action)

		trigger := triggers.NewManualTrigger(utils.NewDefaultEnvironment(), contact, flow.Reference(), nil, utils.Now())
		err = session.Start(trigger)
		require.NoError(t, err)

		run := session.Runs()[0]
		actualEventsJSON, _ := json.Marshal(run.Events())
		expectedEventsJSON, _ := json.Marshal(tc.Events)

		test.AssertEqualJSON(t, expectedEventsJSON, actualEventsJSON, "event mismatch in test '%s'", tc.Description)
	}
}
