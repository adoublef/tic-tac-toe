import React, { MouseEventHandler, useEffect, useRef, useState } from "react";
import { useWebsocket } from "../websocket";
import Cell, { useCell } from "./Cell";

type BoardProps = {
    board: Value[];
    turn: Value;
    winner?: Value;
    play: (index: number) => React.MouseEventHandler<HTMLElement>;
    reset: () => void;
    disable: (value: number) => boolean;
    forfeit: () => Promise<void>;
};

export default function Board(props: BoardProps) {
    return (
        <section className="board">
            {props.board.map((value, index) => (
                <Cell onClick={props.play(index)} disabled={props.disable(value)} key={index} {...useCell(value)} />
            ))}
            <p>current player: {props.turn}</p>
        </section>
    );
}

export function useBoard(boardId: string) {
    // `${import.meta.env.VITE_API_URI}/games/ws?id=${_id}`
    const { send, message } = useWebsocket<Message>(`${import.meta.env.VITE_API_URI}/games/ws?id=${boardId}`, { type: "reset" });

    const [board, setBoard] = useState<(0 | 1 | 2)[]>([0, 0, 0, 0, 0, 0, 0, 0, 0]);
    const [turn, setTurn] = useState<0 | 1 | 2>(1);

    const winner = combinations.reduce<0 | 1 | 2>(evaluate(board), 0);

    const play: (index: number) => MouseEventHandler<HTMLButtonElement> = index => {
        return _ => send({ type: "move", payload: { index, turn } });
    };

    const reset = () => send({ type: "reset" });

    const disable = (value: number) => value !== 0 || winner !== 0;

    const forfeit = async () => {
        send({ type: "bye", payload: "board no longer active" });

        const url = `${import.meta.env.VITE_API_URI}/games/${boardId}`;
        const response = await fetch(url, {
            method: "DELETE",
        });

        if (!response.ok) {
            console.log("error");
            return;
        }

        // This wont run as socket already closed
        // send({ type: "bye", payload: "board no longer active" });
    };

    // TODO -- cannot get rid of this `useEffect`
    useEffect(() => {
        switch (message?.type) {
            case "hi":
            case "bye":
                alert(message.payload);
                break;
            case "move":
                setBoard(board.map((value, i) => (i === message.payload.index ? message.payload.turn : value)));
                setTurn(message.payload.turn === 1 ? 2 : 1);
                break;
            case "reset":
                setBoard([0, 0, 0, 0, 0, 0, 0, 0, 0]);
                setTurn(1);
                break;
            default:
            // TODO -- handle error
        }
    }, [message]);

    return { board, play, turn, reset, winner, disable, forfeit }satisfies BoardProps;
};

type Message =
    | { type: "hi" | "bye"; payload: string; }
    | { type: "move", payload: { index: number, turn: 0 | 1 | 2; }; }
    | { type: "reset"; };

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