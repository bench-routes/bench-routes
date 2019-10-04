interface Routing {
    homePath: string;
    pingPath: string;
    jitterPath: string;
}

export class Routes {
    routes: Routing;

    constructor() {
        this.routes = {
            homePath: 'home.html',
            pingPath: 'ping.html',
            jitterPath: 'jitter.html'
        };
    }
}

