package euchre

import (
	"deck"
)

/*
 * Returns whether a beats b given the current trump suit. a and b are assumed
 * to be different cards. Also it is assumed a leads before b, such that if a
 * and b are two different non-trump suits, a wins automatically.
 *
 * Args:
 *  a: The card that we are asking if it is greater.
 *  b: The card that we are asking if it beats a if it is led.
 *  trump: The current trump suit.
 *
 * Returns:
 *  True if a beats b, if a is led and we are given the trump suit.
 */
func Beat(a deck.Card, b deck.Card, trump deck.Suit) bool {
	var res bool
	// If a is a trump card but b is not, then a wins.
	if (a.AdjSuit(trump) == trump && b.AdjSuit(trump) != trump) ||
		(a.AdjSuit(trump) != trump && b.AdjSuit(trump) == trump) {
		res = a.AdjSuit(trump) == trump
	} else if a.AdjSuit(trump) == trump && b.AdjSuit(trump) == trump {
		// If a is a trump and so is b, then we must compare their values knowing
		// that right and left bower are a rule.
		if a.Value == deck.J || b.Value == deck.J {
			// If a is right bower, then it must win.
			if a.Value == deck.J && a.Suit == trump {
				res = true
			} else if a.Value == deck.J && a.Suit == trump.Left() {
				// If a is left bower, then it wins as long as b is not the right
				// bower.
				res = b.Value != deck.J
			} else {
				// Otherwise, a is not a J, so it is b so b must win.
				res = false
			}
		} else {
			// If neither are one of the bowers, then the values of the cards are
			// compared as normal.
			res = a.Value.Compare(b.Value) > 0
		}
	} else if a.Suit == b.Suit {
		// Otherwise, if they are both the same and they are not both trump, then
		// whoever has the higher value will win.
		res = a.Value.Compare(b.Value) > 0
	} else {
		// And lastly if they have different suits, then a wins automatically since
		// b did not lead.
		res = true
	}

	return res
}

/*
 * Given a player's current hand and the cards that have been played, the
 * possible cards for a player to play are returned. In other words, all cards
 * in the player's hand that match the suit of the led card are returned or all
 * cards otherwise. Also, the actual cards are not returned, rather their
 * position in the hand is returned. This is to make deletion easier.
 *
 * Args:
 *  hand: The player's current cards.
 *  played: The cards that have already been played.
 *  trump: The suit that is currently trump.
 *
 * Returns:
 *  The index of cards in hand that can be played according to euchre rules.
 */
func Possible(hand, played []deck.Card, trump deck.Suit) []int {
	possible := make([]int, 0, len(hand))
	if len(played) > 0 {
		for i := range hand {
			if hand[i].AdjSuit(trump) == played[0].AdjSuit(trump) {
				possible = append(possible, i)
			}
		}
	}

	if len(possible) == 0 {
		for i := range hand {
			possible = append(possible, i)
		}
	}

	return possible
}

/*
 * A function that returns the winning player (using the same number designation
 * as before) based on the trump suit, the cards that have been played, what the
 * player number is for the first player, and if anybody went alone.
 *
 * Args:
 *  played: The cards that were played.
 *  trump: The trump suit.
 *  led: The number designation of the person who played the first card.
 *  alone: The alone player if there is any. If there is not then put in an
 *         invalid player number.
 *
 * Returns:
 *  The number designation of the person who won the trick.
 */
func Winner(played []deck.Card, trump deck.Suit, led int, alone int) int {
	highPlayer := led

	if len(played) >= 2 {
		highest := played[0]
		for i, card := range played[1:] {
			if !Beat(highest, card, trump) {
				highest = card
				highPlayer = (led + i + 1) % 4
			}
		}
	}

	nextAfterHigh := (highPlayer + 1) % 4
	if alone >= 0 && alone < 4 {
		for player := led; player != nextAfterHigh; player = (player + 1) % 4 {
			// If somewhere between the leading player, and the highPlayer the
			// player who is cucked by their partner calling going alone is
			// found, then the highPlayer should be moved up one, since one
			// player is not actually playing.
			if player == (alone+2)%4 {
				highPlayer = (highPlayer + 1) % 4
				break
			}
		}
	}

	return highPlayer
}

/*
 * Provides the index of the card that wins out of the provided played cards and
 * current trump.
 *
 * Args:
 *  played: The cards that have already been played.
 *  trump: The current trump suit the cards are played in the context of.
 *
 * Returns:
 *  The index of the card that is the current winner out of the played cards. If
 *  no cards were played, -1 is returned instead.
 */
func WinnerIdx(played []deck.Card, trump deck.Suit) int {
	highIdx := -1

	if len(played) > 0 {
		highIdx = 0
		highest := played[0]
		for i, card := range played[1:] {
			if !Beat(highest, card, trump) {
				highest = card
				highIdx = i + 1
			}
		}
	}

	return highIdx
}

/*
 * Given the current played cards, the current player's turn and who is going
 * alone, find what player led the hand.
 *
 * Args:
 *  played: The cards currently played.
 *  player: The current player to go.
 *  alone: The player, if any, who is going alone.
 *
 * Returns:
 *  The player number designation that corresponds to the player who led off the
 *  current list of played cards.
 */
func Leader(played []deck.Card, player, alone int) int {
	leader := (player + 4 - len(played)) % 4

	if alone >= 0 {
		empty := (alone + 2) % 4
		if (player > leader && empty >= leader && empty < player) ||
			(player < leader && ((player+4 > empty && empty >= leader) ||
				(player > empty && empty+4 >= leader))) {
			leader--
		}
	}

	return leader
}

/*
 * Given the current player cards and the last player's turn, and who is going
 * alone, find out what player led the hand.
 *
 * Args:
 *  played: The cards that were played.
 *  player: The last player to play, should be in the range [0, 3].
 *  alone: The player who is going alone, if any.
 *
 * Returns:
 *  The player number designation for the player who played the first card in
 *  played.
 */
func LeaderInclusive(played []deck.Card, player, alone int) int {
	return Leader(played[:len(played)-1], player, alone)
}

/*
 * Creates a new copy of the hands in memory from the given state.
 *
 * Args:
 *  state: A euchre state whose player hands are to be copied.
 *
 * Returns:
 *  A copy of the players hands. A [][]deck.Card copy.
 */
func copyAllHands(state State) [][]deck.Card {
	copyHands := make([][]deck.Card, len(state.Hands))
	for i, hand := range state.Hands {
		copyHand := make([]deck.Card, len(hand))
		copy(copyHand, hand)

		copyHands[i] = copyHand
	}

	return copyHands
}

/*
 * Extracts all cards that exist in the given set and have a value of true, into
 * a list of cards.
 *
 * Args:
 *  cardsSet: A set of cards. A card exists if it is in the set and its value is
 *            true.
 *
 * Returns:
 *  A slice of the existing cards in the given set.
 */
func extractAvailableCards(cardsSet map[deck.Card]bool) []deck.Card {
	cards := make([]deck.Card, 0, len(cardsSet))
	for card, exists := range cardsSet {
		if exists {
			cards = append(cards, card)
		}
	}

	return cards
}

/*
 * A private helper method that returns what suits a given player cannot have.
 *
 * Args:
 *  prior: The list of prior tricks.
 *  trump: The current trump suit.
 *
 * Returns:
 *  A list of the suits that a player cannot have indexed by player numbers
 *  in a map.
 */
func noSuits(prior []Trick, trump deck.Suit) map[int][]deck.Suit {
	noSuits := make(map[int][]deck.Suit)

	for i := 0; i < len(prior); i++ {
		// For each trick, find out if a user did not follow suit and therefore
		// does not have this suit.
		trick := prior[i]
		alonePlayer := trick.Alone >= 0 && trick.Alone < 4
		first := trick.Cards[0]

		for player := 0; player < 4; player++ {
			if alonePlayer && player == (trick.Alone+2)%4 {
				continue
			}

			// If between the leading player and the current player, there
			// exists the player who is not playing because his partner is
			// playing alone.
			adjust := 0
			for j := trick.Led; j != player; j = (j + 1) % 4 {
				if alonePlayer && j == (trick.Alone+2)%4 {
					adjust = -1
					break
				}
			}

			// Add 4 to player, so that it is guaranteed to be after trick.Led,
			// but it does not change the final result mod 4.
			playedCard := trick.Cards[(player+4-trick.Led+adjust)%4]
			if first.AdjSuit(trump) != playedCard.AdjSuit(trump) {
				noSuits[player] = append(noSuits[player], first.AdjSuit(trump))
			}
		}
	}

	return noSuits
}
