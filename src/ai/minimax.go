package ai

import "math"

/*
 * Uses minimax adversarial tree search to find the optimal move in a game.
 *
 * Args:
 *  state: The state to start the search from.
 *  engine: The game logic engine for the tree search.
 *
 * Returns:
 *  Gives both the evaluation for the best state and the Move struct associated
 *  with it. This move struct provides the action needed to get to this state
 *  and the state it will send you to.
 */
func Minimax(state TSState, engine TSEngine) (float64, Move) {
	return minimaxHelper(state, engine, math.Inf(-1), math.Inf(1))
}

/*
 * Finds the best move and its evaluation using minimax adversarial search and
 * alpha-beta pruning. This is a helper method used privately by Minimax.
 *
 * Args:
 *  state: The state to start the search from.
 *  engine: The logic engine for the tree search.
 *  alpha: The current alpha value. This should be set to -inf when first called.
 *  beta: The current beta value. This should be set to +inf when first called.
 *
 * Returns:
 *  Gives both the evaluation for the best state and the Move struct associated
 *  with it. This move struct provides the action needed to get to this state
 *  and the state it will send you to.
 */
func minimaxHelper(state TSState, engine TSEngine, alpha float64,
	beta float64) (float64, Move) {
	if engine.IsTerminal(state) {
		return engine.Evaluation(state), Move{nil, state}
	}

	fav := engine.Favorable(state)

	var extremeMove Move
	var extremeValue float64
	if fav {
		extremeValue = math.Inf(-1)
	} else {
		extremeValue = math.Inf(1)
	}

	for _, nextMove := range engine.Successors(state) {
		nextState := nextMove.State
		nextEval, _ := minimaxHelper(nextState, engine, alpha, beta)

		if fav {
			if nextEval > extremeValue {
				extremeValue = nextEval
				extremeMove = nextMove
			}

			alpha = math.Max(alpha, nextEval)
		} else {
			if nextEval < extremeValue {
				extremeValue = nextEval
				extremeMove = nextMove
			}

			beta = math.Min(beta, nextEval)
		}

		if beta < alpha {
			break
		}
	}

	return extremeValue, extremeMove
}
