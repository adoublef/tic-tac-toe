import React, { MouseEventHandler, useEffect, useRef, useState } from "react";
import Cell, { useCell } from "./Cell";

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
        <section className="board">
            {props.board.map((value, index) => (
                <Cell onClick={props.play(index)} disabled={props.disable(value)} key={index} {...useCell(value)} />
            ))}
            <p>current player: {props.turn}</p>
        </section>
    );
}

export function useBoard(url: string) {
    const [board, setBoard] = useState<(0 | 1 | 2)[]>([0, 0, 0, 0, 0, 0, 0, 0, 0]);
    const [turn, setTurn] = useState<0 | 1 | 2>(1);

    const winner = combinations.reduce<0 | 1 | 2>(evaluate(board), 0);

    const play: (index: number) => MouseEventHandler<HTMLButtonElement> = index => {
        return _e => {
            setBoard(board.map((value, i) => (i === index ? turn : value)));
            setTurn(turn === 1 ? 2 : 1);
            ws.current?.send(JSON.stringify({ type: "move", index, turn }));
        };
    };

    const reset = () => { ws.current?.send(JSON.stringify({ type: "reset" })); };

    const disable = (value: number) => value !== 0 || winner !== 0;

    const [message, setMessage] = useState<Message>({ type: "init" });
    const ws = useRef<WebSocket>();

    useEffect(() => {
        const current = new WebSocket(url.replace(/^http/, "ws"));
        current.onopen = () => {
            console.log("websocket opened");
        };

        current.onmessage = (e) => {
            const message = JSON.parse(e.data);
            // console.log(message);
            setMessage(message);
        };

        current.onclose = () => {
            console.log("websocket closed");
        };

        ws.current = current;

        return () => {
            current.close();
        };
    }, [url]);

    // NOTE -- handle websocket messages
    useEffect(() => {
        switch (message.type) {
            case "move":
                console.log("move", message);
                // const { index, turn } = message;
                setBoard(board.map((value, i) => (i === message.payload.index ? message.payload.turn : value)));
                setTurn(message.payload.turn === 1 ? 2 : 1);
                break;
            case "reset":
                console.log("reset", message);
                setBoard([0, 0, 0, 0, 0, 0, 0, 0, 0]);
                setTurn(1);
                break;
            default:
                console.log("default", message);
        }
    }, [message]);

    return { board, play, turn, reset, winner, disable }satisfies BoardProps;
};

// NOTE
type Message =
    | { type: "init"; }
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