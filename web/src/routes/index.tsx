import { Link } from "@tanstack/react-router";
import { useState } from "react";

export default function IndexPage() {
    const create = async () => {
        const url = `${import.meta.env.VITE_API_URI}/games`;
        console.log(url);
        const response = await fetch(url, {
            method: "POST",
        });

        if (!response.ok) {
            console.log("error");
            return;
        }

        const { id } = await response.json();

        const res = (`${location.href}games?id=${id}`);

        window.location.href = res;
    };


    return (
        <div className="App">
            <button onClick={create}>create</button>
            {/* <Link to="/games" disabled>start</Link> */}
            <h2>Tic-Tac-Toe</h2>
            <p>Home</p>
        </div>
    );
}