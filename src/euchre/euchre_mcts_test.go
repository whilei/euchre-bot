package euchre

import (
    "ai"
    "deck"
    "testing"
)

/*
 * Tests the ai package in conjunction with the euchre package. Some of these
 * tests are not real tests, but are simple sanity checks whose output must be
 * manually checked.
 */


/*
 * Tests the output for a run playout for the first card played by the computer.
 *
 * Args:
 *  t - The testing context.
 */
func TestRunPlayout(t *testing.T) {
    setup := Setup {
        1,
        1,
        true,
        deck.Card{ deck.D, deck.Nine },
        deck.D,
        deck.Card{ },
    }

    hand := []deck.Card {
        deck.Card {
            deck.H,
            deck.Nine,
        },
        deck.Card {
            deck.H,
            deck.Ten,
        },
        deck.Card {
            deck.S,
            deck.A,
        },
        deck.Card {
            deck.D,
            deck.Q,
        },
        deck.Card {
            deck.C,
            deck.Q,
        },
    }

    played := []deck.Card {
        deck.Card {
            deck.C,
            deck.J,
        },
        deck.Card {
            deck.C,
            deck.A,
        },
    }

    var prior []Trick

    s := NewState(setup, 0, hand, played, prior, deck.Card{ })
    n := ai.NewNode()
    n.Value(s)
    e := Engine{ }

    ai.RunPlayoutDebug(n, e)
}