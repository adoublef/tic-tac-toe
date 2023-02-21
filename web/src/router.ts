import { RootRoute, Route, Router } from "@tanstack/router";
import IndexPage from "./routes";
import GamesPage from "./routes/games";

const rootRoute = new RootRoute();

const homeRoute = new Route({
    getParentRoute: () => rootRoute,
    path: "/",
    component: IndexPage,
});

const gameRoute = new Route({
    getParentRoute: () => rootRoute,
    path: "/games",
    component: GamesPage,
});

const routeTree = rootRoute.addChildren([homeRoute, gameRoute]);

export const router = new Router({ routeTree });

declare module '@tanstack/router' {
    interface Register {
        router: typeof router;
    }
}

// https://github.com/rog-golang-buddies/rapidmidiex-research/blob/main/examples/ws-noughts-crosses/assets/components/useGame.js
// https://github.com/rog-golang-buddies/rapidmidiex-research/blob/main/examples/ws-noughts-crosses/assets/components/Game.js
