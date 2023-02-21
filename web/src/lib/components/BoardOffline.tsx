import React, { MouseEventHandler, useState } from "react";
import Cell, { useCell } from "./Cell";

// Board.tsx
type BoardProps = {
    board: Value[];
    turn: Value;
    winner?: Value;
    play: (index: number) => React.MouseEventHandler<HTMLElement>;
    reset: () => void;
    disable: (value: number) => boolean;
};

export default function Board(props: BoardProps) {
    return (
        <section>
            {props.board.map((value, index) => (
                <Cell onClick={props.play(index)} disabled={props.disable(value)} key={index} {...useCell(value)} />
            ))}
        </section>
    );
}

export function useBoard() {
    const [board, setBoard] = useState<(0 | 1 | 2)[]>([0, 0, 0, 0, 0, 0, 0, 0, 0]);
    const [turn, setTurn] = useState<0 | 1 | 2>(1);

    const winner = combinations.reduce<0 | 1 | 2>(evaluate(board), 0);

    const play: (index: number) => MouseEventHandler<HTMLButtonElement> = index => {
        return _e => {
            setBoard(board.map((value, i) => (i === index ? turn : value)));
            setTurn(turn === 1 ? 2 : 1);
        };
    };

    const reset = () => { setBoard([0, 0, 0, 0, 0, 0, 0, 0, 0]); setTurn(1); };

    const disable = (value: number) => value !== 0 || winner !== 0;

    return { board, play, turn, reset, winner, disable }satisfies BoardProps;
}

const combinations = [
    [0, 1, 2],
    [3, 4, 5],
    [6, 7, 8],
    [0, 3, 6],
    [1, 4, 7],
    [2, 5, 8],
    [0, 4, 8],
    [2, 4, 6],
] as const;

type Value = 0 | 1 | 2;

type Combination = typeof combinations[number];

const evaluate = (board: Value[]) => (acc: Value, [a, b, c]: Combination) => {
    // TODO -- if acc is not 0, return acc
    if (acc !== 0) return acc;
    return (board[a] === board[b] && board[a] === board[c]) ? board[a] : acc;
};