import { useState, useRef, useEffect } from "react";

export function useWebsocket<T>(url: string, init?: T) {
    const [message, setMessage] = useState<T | undefined>(init);
    const ws = useRef<WebSocket>();

    useEffect(() => {
        const current = new WebSocket(url.replace(/^http/, "ws"));
        current.onopen = () => {
            current.send(JSON.stringify({ type: "hi", payload: "websocket opened" }));
        };

        current.onmessage = (e) => {
            const message: T = JSON.parse(e.data);
            setMessage(message);
        };

        ws.current = current;

        return () => {
            current.close();
        };
    }, [url]);

    const send = (message: T) => ws.current?.send(JSON.stringify(message));

    return { message, send };
};