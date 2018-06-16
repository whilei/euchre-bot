package player

import (
	"deck"
	"testing"
)

const (
	PICKUP_CONF             = 0.6
	CALL_CONF               = 0.6
	ALONE_CONF              = 1.2
	PICKUP_RUNS             = 5000
	PICKUP_DETERMINIZATIONS = 50
	CALL_RUNS               = 5000
	CALL_DETERMINIZATIONS   = 50
	PLAY_RUNS               = 5000
	PLAY_DETERMINIZATIONS   = 50
	ALONE_RUNS              = 5000
	ALONE_DETERMINIZATIONS  = 50
)

/*
 * Tests the different types of players functionality.
 */

/*
 * A test that defines the inputs and the expected output for a given test case.
 */
type discardTest struct {
	hand     []deck.Card
	top      deck.Card
	expected deck.Card
}

var discardTests = []discardTest{
	/*
	 * Test that the lowest trump card is discarded if all cards are trump.
	 */
	discardTest{
		[]deck.Card{
			deck.Card{deck.H, deck.J},
			deck.Card{deck.H, deck.Ten},
			deck.Card{deck.H, deck.Nine},
			deck.Card{deck.H, deck.A},
			deck.Card{deck.D, deck.J},
		},
		deck.Card{deck.H, deck.Q},
		deck.Card{deck.H, deck.Nine},
	},

	/*
	 * Whitebox testing where trumps are in ascending order. The test is to see
	 * if the lowest of a trump suit will be chosen independent of order.
	 */
	discardTest{
		[]deck.Card{
			deck.Card{deck.C, deck.Nine},
			deck.Card{deck.C, deck.Ten},
			deck.Card{deck.C, deck.Q},
			deck.Card{deck.C, deck.K},
			deck.Card{deck.C, deck.J},
		},
		deck.Card{deck.C, deck.A},
		deck.Card{deck.C, deck.Nine},
	},

	/*
	 * Another whitebox test to assure that discarding properly handles bowers.
	 */
	discardTest{
		[]deck.Card{
			deck.Card{deck.C, deck.J},
			deck.Card{deck.S, deck.J},
			deck.Card{deck.C, deck.Q},
			deck.Card{deck.C, deck.K},
			deck.Card{deck.C, deck.A},
		},
		deck.Card{deck.C, deck.Ten},
		deck.Card{deck.C, deck.Ten},
	},

	/*
	 * When there is a suit with only one card in it, but it is an A, do not discard
	 * it since it is valuable. Discard the lowest card of a non-trump suit.
	 */
	discardTest{
		[]deck.Card{
			deck.Card{deck.C, deck.A},
			deck.Card{deck.H, deck.Nine},
			deck.Card{deck.H, deck.J},
			deck.Card{deck.S, deck.K},
			deck.Card{deck.S, deck.Q},
		},
		deck.Card{deck.H, deck.Q},
		deck.Card{deck.S, deck.Q},
	},

	/*
	 * Test that a suit with only one card is discarded. This makes sense if you
	 * have trumps. Otherwise, the lowest overall card should be dropped.
	 */
	discardTest{
		[]deck.Card{
			deck.Card{deck.C, deck.Q},
			deck.Card{deck.H, deck.Nine},
			deck.Card{deck.H, deck.J},
			deck.Card{deck.S, deck.K},
			deck.Card{deck.S, deck.Q},
		},
		deck.Card{deck.H, deck.Q},
		deck.Card{deck.C, deck.Q},
	},

	/*
	 * When there is only one card for a certain suit but no trumps, just
	 * discard the lowest valued card.
	 */
	discardTest{
		[]deck.Card{
			deck.Card{deck.C, deck.Q},
			deck.Card{deck.H, deck.Nine},
			deck.Card{deck.H, deck.A},
			deck.Card{deck.S, deck.K},
			deck.Card{deck.S, deck.Q},
		},
		deck.Card{deck.D, deck.A},
		deck.Card{deck.H, deck.Nine},
	},

	/*
	 * When there is more than one card that is non-ace, non-trump and solitary
	 * in its suit, then choose the smallest one.
	 */
	discardTest{
		[]deck.Card{
			deck.Card{deck.D, deck.J},
			deck.Card{deck.S, deck.A},
			deck.Card{deck.H, deck.Nine},
			deck.Card{deck.S, deck.Q},
			deck.Card{deck.C, deck.J},
		},
		deck.Card{deck.S, deck.J},
		deck.Card{deck.H, deck.Nine},
	},
}

/*
 * The main driver for discard tests that runs all discardTests outlined above
 * globally on all testable players.
 */
func TestDiscard(t *testing.T) {
	players := getTestablePlayers()

	for i, fixture := range discardTests {
		for j, player := range players {
			copyHand := make([]deck.Card, len(fixture.hand))
			copy(copyHand, fixture.hand)

			_, discarded := player.Discard(copyHand, fixture.top)
			if discarded != fixture.expected {
				t.Logf("Fixture %d, implementation %d failed.\n", i+1, j+1)
				t.Errorf("Gave %s instead of %s.\n", discarded, fixture.expected)
			}
		}
	}
}

/*
 * Returns a list of all the different Player implementations to test.
 *
 * Returns:
 *  A list of the different player implementations to test in this file. The
 *  order of the players is [rule, smart].
 */
func getTestablePlayers() []Player {
	rule := NewRule("")
	smart := NewSmart(PICKUP_CONF, CALL_CONF, ALONE_CONF,
		PICKUP_RUNS, PICKUP_DETERMINIZATIONS,
		CALL_RUNS, CALL_DETERMINIZATIONS,
		PLAY_RUNS, PLAY_DETERMINIZATIONS,
		ALONE_RUNS, ALONE_DETERMINIZATIONS)

	players := []Player{rule, smart}

	return players
}
