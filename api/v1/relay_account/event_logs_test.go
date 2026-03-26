package relayaccount

import "testing"

func TestGetEventPayloadDeposit(t *testing.T) {
	payload := getEventPayload(3, "3-0xabc123-dymension")
	if payload != "{\"network\":\"dymension\",\"tx_hash\":\"0xabc123\"}" {
		t.Fatalf("unexpected payload: %s", payload)
	}
}

func TestGetEventPayloadLegacyDepositReturnsEmpty(t *testing.T) {
	cases := []string{
		"",
		"3-0xabc123",
		"4-0xabc123-dymension",
		"invalid",
	}

	for _, c := range cases {
		payload := getEventPayload(3, c)
		if payload != "{}" {
			t.Fatalf("expected empty payload for reason %q, got %s", c, payload)
		}
	}
}

func TestGetEventPayloadNonDepositReturnsEmpty(t *testing.T) {
	payload := getEventPayload(4, "4-task-id")
	if payload != "{}" {
		t.Fatalf("expected empty payload for non-deposit event, got %s", payload)
	}
}
