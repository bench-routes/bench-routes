interface Message {
    message: string;
    data?: string;
}

export class Sockets {
    Socket: WebSocket;
    SocketURL: string;

    constructor() {
        this.SocketURL = 'ws://localhost:9090/websocket';
        this.Socket = new WebSocket(this.SocketURL);
    }

    /**
     * Sends message synchronously to the service
     * @param message => message text to be sennt to the backend service
     */
    sendMessage(message: Message): void {
        this.Socket.send(message.toString());
    }

    /**
     * Closes the Socket connection with the service
     */
    close(): void {
        this.Socket.close();
    }

    /**
     * Closes the Socket connection with the service
     */
    connect(): void {
        this.Socket.OPEN;
    }

    /**
     * Returns a socket object.
     * @returns a websocket connection instance.
     */
    get(): WebSocket {
        return this.Socket;
    }


}