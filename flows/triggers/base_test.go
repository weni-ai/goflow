package triggers_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sort"
	"testing"
	"time"

	"github.com/nyaruka/gocommon/dates"
	"github.com/nyaruka/gocommon/jsonx"
	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/gocommon/uuids"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/assets/static"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/engine"
	"github.com/nyaruka/goflow/flows/triggers"
	"github.com/nyaruka/goflow/test"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTriggerTypes(t *testing.T) {
	assetsJSON, err := ioutil.ReadFile("testdata/_assets.json")
	require.NoError(t, err)

	typeNames := make([]string, 0)
	for typeName := range triggers.RegisteredTypes() {
		typeNames = append(typeNames, typeName)
	}

	sort.Strings(typeNames)

	for _, typeName := range typeNames {
		testTriggerType(t, assetsJSON, typeName)
	}
}

func testTriggerType(t *testing.T, assetsJSON json.RawMessage, typeName string) {
	testPath := fmt.Sprintf("testdata/%s.json", typeName)
	testFile, err := ioutil.ReadFile(testPath)
	require.NoError(t, err)

	tests := []struct {
		Description string          `json:"description"`
		Trigger     json.RawMessage `json:"trigger"`
		ReadError   string          `json:"read_error,omitempty"`
		Events      json.RawMessage `json:"events,omitempty"`
		Context     json.RawMessage `json:"context,omitempty"`
	}{}

	err = jsonx.Unmarshal(testFile, &tests)
	require.NoError(t, err)

	defer dates.SetNowSource(dates.DefaultNowSource)
	defer uuids.SetGenerator(uuids.DefaultGenerator)

	for i, tc := range tests {
		dates.SetNowSource(dates.NewFixedNowSource(time.Date(2018, 10, 18, 14, 20, 30, 123456, time.UTC)))
		uuids.SetGenerator(uuids.NewSeededGenerator(12345))

		testName := fmt.Sprintf("test '%s' for trigger type '%s'", tc.Description, typeName)

		// create session assets
		sa, err := test.CreateSessionAssets(assetsJSON, "")
		require.NoError(t, err, "unable to create session assets in %s", testName)

		trigger, err := triggers.ReadTrigger(sa, tc.Trigger, assets.PanicOnMissing)

		if tc.ReadError != "" {
			rootErr := errors.Cause(err)
			assert.EqualError(t, rootErr, tc.ReadError, "read error mismatch in %s", testName)
			continue
		} else {
			assert.NoError(t, err, "unexpected read error in %s", testName)
		}

		// start a session with this trigger
		eng := engine.NewBuilder().Build()
		session, sprint, err := eng.NewSession(sa, trigger)
		assert.NoError(t, err)

		assert.Equal(t, flows.FlowTypeMessaging, session.Type())
		assert.NotNil(t, session.Environment())

		// clone test case and populate with actual values
		actual := tc
		actual.Events, _ = jsonx.Marshal(sprint.Events())

		actualContextJSON, err := session.Runs()[0].EvaluateTemplate(`@(json(trigger))`)
		assert.NoError(t, err)
		actual.Context = []byte(actualContextJSON)

		// re-marshal the trigger
		actual.Trigger, err = jsonx.Marshal(trigger)
		require.NoError(t, err)

		if !test.UpdateSnapshots {
			// check events generated by trigger
			test.AssertEqualJSON(t, tc.Events, actual.Events, "events mismatch in %s", testName)

			// check context representation
			test.AssertEqualJSON(t, tc.Context, actual.Context, "context mismatch in %s", testName)

			// check marshalled
			test.AssertEqualJSON(t, tc.Trigger, actual.Trigger, "marshal mismatch in %s", testName)
		} else {
			tests[i] = actual
		}
	}

	if test.UpdateSnapshots {
		actualJSON, err := jsonx.MarshalPretty(tests)
		require.NoError(t, err)

		err = ioutil.WriteFile(testPath, actualJSON, 0666)
		require.NoError(t, err)
	}
}

var assetsJSON = `{
	"flows": [
		{
			"uuid": "7c37d7e5-6468-4b31-8109-ced2ef8b5ddc",
			"name": "Registration",
            "spec_version": "13.0",
            "language": "eng",
            "type": "messaging",
            "nodes": []
        }
	],
	"channels": [
		{
			"uuid": "8cd472c4-bb85-459a-8c9a-c04708af799e",
			"name": "Facebook",
			"address": "23532562626",
			"schemes": ["facebook"],
			"roles": ["send", "receive"]
		},
		{
            "uuid": "3a05eaf5-cb1b-4246-bef1-f277419c83a7",
            "name": "Nexmo",
            "address": "+16055742523",
            "schemes": ["tel"],
            "roles": ["send", "receive"]
        }
	],
    "ticketers": [
        {
            "uuid": "19dc6346-9623-4fe4-be80-538d493ecdf5",
            "name": "Support Tickets",
            "type": "mailgun"
        }
    ],
    "users": [
        {
            "email": "bob@nyaruka.com",
            "name": "Bob McTickets"
        }
    ]
}`

func TestTriggerMarshaling(t *testing.T) {
	defer dates.SetNowSource(dates.DefaultNowSource)
	dates.SetNowSource(dates.NewFixedNowSource(time.Date(2018, 10, 20, 9, 49, 30, 1234567890, time.UTC)))

	uuids.SetGenerator(uuids.NewSeededGenerator(1234))
	defer uuids.SetGenerator(uuids.DefaultGenerator)

	env := envs.NewBuilder().Build()

	source, err := static.NewSource([]byte(assetsJSON))
	require.NoError(t, err)

	sa, err := engine.NewSessionAssets(env, source, nil)
	require.NoError(t, err)

	flow := assets.NewFlowReference("7c37d7e5-6468-4b31-8109-ced2ef8b5ddc", "Registration")
	channel := assets.NewChannelReference("3a05eaf5-cb1b-4246-bef1-f277419c83a7", "Nexmo")
	ticketer := sa.Ticketers().Get("19dc6346-9623-4fe4-be80-538d493ecdf5")
	user := sa.Users().Get("bob@nyaruka.com")
	ticket := flows.NewTicket("276c2e43-d6f9-4c36-8e54-b5af5039acf6", ticketer, "Problem", "Where are my shoes?", "123456", user)

	contact := flows.NewEmptyContact(sa, "Bob", envs.Language("eng"), nil)
	contact.AddURN(urns.URN("tel:+12065551212"), nil)

	eng := engine.NewBuilder().Build()
	session, _, err := eng.NewSession(sa, triggers.NewBuilder(env, flow, contact).Manual().Build())
	require.NoError(t, err)

	history := flows.NewChildHistory(session)

	// can't create a trigger with invalid JSON
	assert.Panics(t, func() {
		triggers.NewBuilder(env, flow, contact).FlowAction(history, json.RawMessage(`{"uuid"}`)).Build()
	})
	assert.Panics(t, func() {
		triggers.NewBuilder(env, flow, contact).FlowAction(history, nil).Build()
	})

	triggerTests := []struct {
		trigger  flows.Trigger
		snapshot string
	}{
		{
			triggers.NewBuilder(env, flow, contact).
				Campaign(triggers.NewCampaignReference("8cd472c4-bb85-459a-8c9a-c04708af799e", "Reminders"), "8d339613-f0be-48b7-92ee-155f4c7576f8").
				Build(),
			"campaign",
		},
		{
			triggers.NewBuilder(env, flow, contact).
				Channel(channel, triggers.ChannelEventTypeIncomingCall).
				WithConnection(urns.URN("tel:+12065551212")).
				Build(),
			"channel_incoming_call",
		},
		{
			triggers.NewBuilder(env, flow, contact).
				Channel(channel, triggers.ChannelEventTypeNewConversation).
				WithParams(types.NewXObject(map[string]types.XValue{"foo": types.NewXText("bar")})).
				Build(),
			"channel_new_conversation",
		},
		{
			triggers.NewBuilder(env, flow, contact).
				FlowAction(history, json.RawMessage(`{"uuid": "084e4bed-667c-425e-82f7-bdb625e6ec9e"}`)).
				Build(),
			"flow_action",
		},
		{
			triggers.NewBuilder(env, flow, contact).
				FlowAction(history, json.RawMessage(`{"uuid": "084e4bed-667c-425e-82f7-bdb625e6ec9e"}`)).
				WithConnection(channel, "tel:+12065551212").
				AsBatch().
				Build(),
			"flow_action_ivr",
		},
		{
			triggers.NewBuilder(env, flow, contact).
				Manual().
				WithParams(types.NewXObject(map[string]types.XValue{"foo": types.NewXText("bar")})).
				WithUser(user).
				WithOrigin("api").
				AsBatch().
				Build(),
			"manual",
		},
		{
			triggers.NewBuilder(env, flow, contact).
				Manual().
				WithConnection(channel, "tel:+12065551212").
				WithParams(types.NewXObject(map[string]types.XValue{"foo": types.NewXText("bar")})).
				AsBatch().
				Build(),
			"manual_ivr",
		},
		{
			triggers.NewBuilder(env, flow, contact).
				Manual().
				Build(),
			"manual_minimal",
		},
		{
			triggers.NewBuilder(env, flow, contact).
				Msg(flows.NewMsgIn(flows.MsgUUID("c8005ee3-4628-4d76-be66-906352cb1935"), urns.URN("tel:+1234567890"), channel, "Hi there", nil)).
				WithMatch(triggers.NewKeywordMatch(triggers.KeywordMatchTypeFirstWord, "hi")).
				Build(),
			"msg",
		},
		{
			triggers.NewBuilder(env, flow, contact).
				Ticket(ticket, triggers.TicketEventTypeClosed).
				Build(),
			"ticket_closed",
		},
	}

	for _, tc := range triggerTests {
		triggerJSON, err := jsonx.MarshalPretty(tc.trigger)
		assert.NoError(t, err)

		test.AssertSnapshot(t, tc.snapshot, string(triggerJSON))

		// then try to read from the JSON
		_, err = triggers.ReadTrigger(sa, triggerJSON, assets.PanicOnMissing)
		assert.NoError(t, err, "error reading trigger: %s", string(triggerJSON))
	}
}

func TestReadTrigger(t *testing.T) {
	env := envs.NewBuilder().Build()

	missingAssets := make([]assets.Reference, 0)
	missing := func(a assets.Reference, err error) { missingAssets = append(missingAssets, a) }

	sessionAssets, err := engine.NewSessionAssets(env, static.NewEmptySource(), nil)
	require.NoError(t, err)

	// error if no type field
	_, err = triggers.ReadTrigger(sessionAssets, []byte(`{"foo": "bar"}`), missing)
	assert.EqualError(t, err, "field 'type' is required")

	// error if we don't recognize action type
	_, err = triggers.ReadTrigger(sessionAssets, []byte(`{"type": "do_the_foo", "foo": "bar"}`), missing)
	assert.EqualError(t, err, "unknown type: 'do_the_foo'")
}

func TestTriggerSessionInitialization(t *testing.T) {
	env := envs.NewBuilder().WithDateFormat(envs.DateFormatMonthDayYear).Build()

	source, err := static.NewSource([]byte(assetsJSON))
	require.NoError(t, err)

	sa, err := engine.NewSessionAssets(env, source, nil)
	require.NoError(t, err)

	defaultEnv := envs.NewBuilder().Build()

	flow := assets.NewFlowReference(assets.FlowUUID("7c37d7e5-6468-4b31-8109-ced2ef8b5ddc"), "Registration")

	contact := flows.NewEmptyContact(sa, "Bob", envs.Language("eng"), nil)
	contact.AddURN(urns.URN("tel:+12065551212"), nil)

	params := types.NewXObject(map[string]types.XValue{"foo": types.NewXText("bar")})

	trigger := triggers.NewBuilder(env, flow, contact).Manual().WithParams(params).Build()

	assert.Equal(t, triggers.TypeManual, trigger.Type())
	assert.Equal(t, env, trigger.Environment())
	assert.Equal(t, contact, trigger.Contact())
	assert.Nil(t, trigger.Connection())
	assert.Equal(t, params, trigger.Params())

	eng := engine.NewBuilder().Build()
	session, _, err := eng.NewSession(sa, trigger)
	require.NoError(t, err)

	assert.Equal(t, flows.FlowTypeMessaging, session.Type())
	assert.Equal(t, contact, session.Contact())
	assert.Equal(t, env, session.Environment())
	assert.Equal(t, flow, session.Runs()[0].FlowReference())

	// contact, environment and params are optional
	trigger = triggers.NewBuilder(nil, flow, nil).Manual().Build()

	assert.Equal(t, triggers.TypeManual, trigger.Type())
	assert.Nil(t, trigger.Environment())
	assert.Nil(t, trigger.Contact())
	assert.Nil(t, trigger.Params())

	session, _, err = eng.NewSession(sa, trigger)
	require.NoError(t, err)

	assert.Equal(t, flows.FlowTypeMessaging, session.Type())
	assert.Nil(t, session.Contact())
	assert.Equal(t, defaultEnv, session.Environment()) // uses defaults
}

func TestTriggerContext(t *testing.T) {
	env := envs.NewBuilder().Build()

	source, err := static.NewSource([]byte(assetsJSON))
	require.NoError(t, err)

	sa, err := engine.NewSessionAssets(env, source, nil)
	require.NoError(t, err)

	flow := assets.NewFlowReference(assets.FlowUUID("7c37d7e5-6468-4b31-8109-ced2ef8b5ddc"), "Registration")
	user := sa.Users().Get("bob@nyaruka.com")

	contact := flows.NewEmptyContact(sa, "Jim", envs.Language("eng"), nil)
	contact.AddURN(urns.URN("tel:+12065551212"), nil)

	params := types.NewXObject(map[string]types.XValue{"foo": types.NewXText("bar")})
	trigger := triggers.NewBuilder(env, flow, contact).
		Manual().
		WithParams(params).
		WithUser(user).
		WithOrigin("api").
		AsBatch().
		Build()

	test.AssertXEqual(t, types.NewXObject(map[string]types.XValue{
		"type":    types.NewXText("manual"),
		"params":  params,
		"keyword": types.XTextEmpty,
		"user": types.NewXObject(map[string]types.XValue{
			"__default__": types.NewXText("Bob McTickets"),
			"email":       types.NewXText("bob@nyaruka.com"),
			"name":        types.NewXText("Bob McTickets"),
			"first_name":  types.NewXText("Bob"),
		}),
		"origin": types.NewXText("api"),
		"ticket": nil,
	}), flows.Context(env, trigger))
}
