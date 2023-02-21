import Game from "../lib/components/Game";

export default function GamesPage() {
    const id = new URLSearchParams(window.location.search).get("id");
    if (!id) {
        return <p>no game id</p>;
    }


    // copy to clipboard
    const copy = () => {
        navigator.clipboard.writeText(location.href);
    };


    return (
        <div className="App" >
            <button onClick={copy}>share</button>
            <h2>Tic-Tac-Toe</h2>
            <Game id={id} />
        </div>
    );
}