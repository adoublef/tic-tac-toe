import React, { createContext, useContext, useState, useEffect, useRef, useCallback, useMemo } from "react";

export function useWebsocket<T>(url: string, callback: (message: T) => void, deps = []) {
    const ws = useRef<WebSocket>();

    const onmessage = useMemo(() => callback, [callback].concat(deps));

    const send = (message: T) => {
        ws.current?.send(JSON.stringify(message));
    };

    useEffect(() => {
        const current = new WebSocket(url);
        {
            current.addEventListener("open", (e) => {
                console.log("websocket opened");
            });

            current.addEventListener("message", (e) =>
                onmessage(JSON.parse(e.data)));

            current.addEventListener("error", (e) => {
                console.log("websocket error");
            });

            current.addEventListener("close", (e) => {
                console.log("websocket error");
            });
        }
        ws.current = current;

        return () => { current.close(); };
    }, [url, onmessage]);

    return [send] as const;
};
