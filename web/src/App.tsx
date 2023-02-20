import './App.css';
import Game from "./lib/components/Game";

// https://github.com/rog-golang-buddies/rapidmidiex-research/blob/main/examples/ws-noughts-crosses/assets/components/useGame.js
// https://github.com/rog-golang-buddies/rapidmidiex-research/blob/main/examples/ws-noughts-crosses/assets/components/Game.js
export default function App() {
    return (
        <div className="App" >
            <h2>Tic-Tac-Toe</h2>
            <Game />
        </div>
    );
}
