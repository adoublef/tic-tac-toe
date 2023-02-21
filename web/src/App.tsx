import './App.css';
import Game from "./lib/components/Game";

// https://github.com/rog-golang-buddies/rapidmidiex-research/blob/main/examples/ws-noughts-crosses/assets/components/useGame.js
// https://github.com/rog-golang-buddies/rapidmidiex-research/blob/main/examples/ws-noughts-crosses/assets/components/Game.js
export default function App() {
    const id = new URLSearchParams(window.location.search).get("id");
    if (!id) {
        return <p>no game id</p>;
    }

    return (
        <div className="App" >
            <h2>Tic-Tac-Toe</h2>
            <Game id={id} />
        </div>
    );
}
