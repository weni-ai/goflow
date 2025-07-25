1.14.0
----------
 * feat: add Category field to Template and TemplateReference, update related constructors and tests

1.13.0
----------
 * feat: extend MsgWppOut structure to include ActionType and ActionExternalID fields, and update related functions

1.12.0
----------
 * feat: add regex matching for Linx API in SendMsgCatalogAction

1.11.0
----------
 * Modify service logic to include search keywords in call results

1.10.0
----------
 * Add support for Instagram replies: comments, by tag and private reply

1.9.0
----------
 * feat: add cart simulation parameter support

1.8.0
----------
 * feat: add FromBCP47 function to convert BCP47 codes to Locale

1.7.0
----------
 * Add support for catalog message in whatsapp messages

1.6.1
----------
 * Add field type getter

1.6.0
----------
 * Add extra prompt in org context

1.5.5
----------
 * fix: list message translations
 * feat: add button component to wpp message

1.5.4
----------
 * Fix release

1.5.3
----------
 * Add new broadcast_type attribute

1.5.2
----------
 * Add queue uuid on topic

1.5.1
----------
 * Support MsgTemplating in WhatsApp messages

1.5.0
----------
 * Use hideUnavailable for product search

1.4.2
----------
 * Using URN identity for brain

1.4.1
----------
 * Feat: Accept order-like items in order details messages

1.4.0
----------
 * Feat: WhatsApp Order Details

1.3.0
----------
 * Add support for attachments in whatsapp flows

1.2.1
----------
 * Fix validation for dynamic lists

1.2.0
----------
 * Refactor evaluateMessageWpp to message list
 * Manipulate function return to avoid errors

1.1.0
----------
 * Add support to whatsapp flows

1.0.0
----------
 * Adjustments for sponsored product support

0.14.3-goflow-0.144.3
----------
 * Fix: Brain only send attachments when entry is @input.text

0.14.2-goflow-0.144.3
----------
 * Fix: Allow vtex intelligent search to receive locale query param

0.14.1-goflow-0.144.3
----------
 * Fix contact context identation for editor gen

0.14.0-goflow-0.144.3
----------
 * Add entry field for call_brain

0.13.1-goflow-0.144.3
----------
 * Allow re-assign urn to another channel of the same type #67

0.13.0-goflow-0.144.3
----------
 * Change response_json type in NFMReply

0.12.0-goflow-0.144.3
----------
 * Support for CTA message for whatsapp cloud channels

0.11.3-goflow-0.144.3
----------
 * Add new apiType to send products

0.11.2-goflow-0.144.3
----------
 * Adjust tests for brain webhook calls

0.11.1-goflow-0.144.3
----------
 * Adjust brain webhook calls

0.11.0-goflow-0.144.3
----------
 * Adjust Brain Card

0.10.0-goflow-0.144.3
----------
 * Add action to brain card

0.9.1-goflow-0.144.3
----------
 * Add webhook header to custom timeout for requests

0.9.0-goflow-0.144.3
----------
 * Implementations for whatsapp message card

0.8.2-goflow-0.144.3
----------
 * Add whatsapp header token to webhooks

0.8.1-goflow-0.144.3
----------
 * Fix external service using repeated param value on sprint with many runs on same flow node.

0.8.0-goflow-0.144.3
----------
 * Remove regex that removes whitespace for zeroshot category
 * Add postal code field to catalog msgs

0.7.1-goflow-0.144.3
----------
 * Fix error handling for empty product list

0.7.0-goflow-0.144.3
----------
 * Change the product list structure to preserve insertion order

0.6.3-goflow-0.144.3
----------
 * Add only the text content in the result value for wenigpt
 * Fix languages for request zeroshot

0.6.2-goflow-0.144.3
----------
 * Fix limitation of product sections

0.6.1-goflow-0.144.3
----------
 * Add nfm_reply field to input
 * Add responseJSON to the result value for wenigpt

0.6.0-goflow-0.144.3
----------
 * Implement action for WeniGPT Call card

0.5.2-goflow-0.144.3
----------
 * Add unit tests for OrgContext

0.5.1-goflow-0.144.3
----------
 * Allow the context to be empty if it doesn't have one

0.5.0-goflow-0.144.3
----------
 * Add asset to zeroshot context

0.4.1-goflow-0.144.3
----------
 * Add vtex search support with sellerId

0.4.0-goflow-0.144.3
----------
 * Add support for vtex searches

0.3.0-goflow-0.144.3
----------
 * Add catalog msg card support

0.2.3-goflow-0.144.3
----------
 * Add language field to zeroshot request 

0.2.2-goflow-0.144.3
----------
 * Fix msg and resume order
 * Add token for zeroshot in request header

0.2.1-goflow-0.144.3
----------
 * Fixes and improvements for smart router

0.2.0-goflow-0.144.3
----------
 * Add Smart Router as new router type

0.1.2-goflow-0.144.3
----------
 * Add support for session.input.order and resume.params

0.1.1-goflow-0.144.3
----------
 * Add zendesk ticket reopen check

0.1.0-goflow-0.144.3
----------
 * Add support for external services actions and events
 * Add Support for trigger.params in Msg events
 * Add support for trigger.params in ticket events
