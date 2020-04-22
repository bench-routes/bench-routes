import configurations from '../Services/constance';

export class PersistentConnection {
  private readonly socket: WebSocket;

  constructor() {
    this.socket = new WebSocket(configurations.brServiceConnection);
  }

  public getSocketInstance(): WebSocket {
    return this.socket;
  }

  public forceClose(): void {
    this.socket.close();
  }

  public send(data: string): void {
    this.socket.send(data);
  }
}