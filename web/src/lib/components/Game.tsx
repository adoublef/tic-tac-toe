import Board, { useBoard } from "./Board";

type GameProps = { id: string; };

export default function Game(props: GameProps) {
    const state = useBoard(`${import.meta.env.VITE_API_URI}/game/ws?id=${props.id}`);

    // disable reset button if no moves left or game is over
    const gameOver = state.winner === 0 && state.board.includes(0);

    return (
        <div>
            <p>current player {state.turn}</p>
            <Board {...state} />
            <p>{state.board.toString()}</p>
            <button disabled={gameOver} onClick={state.reset}>reset</button>
            <button>forfeit</button>
            <p>game status: {state.winner}</p>
        </div>
    );
}
