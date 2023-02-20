import Board, { useBoard } from "./BoardWS";


export default function Game() {
    const state = useBoard(`${import.meta.env.VITE_API_URI}/game/ws`);

    // disable reset button if no moves left or game is over
    const gameOver = state.winner === 0 && state.board.includes(0);

    return (
        <div>
            <p>current player {state.turn}</p>
            <Board {...state} />
            <p>{state.board.toString()}</p>
            <p>moves left: {`${!state.board.includes(0)}`}</p>
            <button disabled={gameOver} onClick={state.reset}>reset</button>
            <button>forfeit</button>
            <p>game status: {state.winner}</p>
        </div>
    );
}
