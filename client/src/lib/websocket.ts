/**
 * WebSocketClient handles real-time communication with the WebSocket server.
 * It provides automatic reconnection, message handling, and event-based communication.
 */
export class WebSocketClient {
  private ws: WebSocket | null = null;
  private reconnectAttempts = 0;
  private maxReconnectAttempts = 5;
  private reconnectTimeout = 1000; // Start with 1 second

  /**
   * Creates a new WebSocketClient instance.
   * @param url - The WebSocket server URL to connect to
   */
  constructor(private url: string) {}

  /**
   * Establishes a WebSocket connection to the server.
   * Handles connection setup, event listeners, and automatic reconnection.
   */
  connect() {
    try {
      if (this.ws?.readyState === WebSocket.CONNECTING || this.ws?.readyState === WebSocket.OPEN) {
        console.log('WebSocket connection already exists');
        return;
      }

      this.ws = new WebSocket(this.url);

      this.ws.onopen = () => {
        console.log('WebSocket connected');
        this.reconnectAttempts = 0;
      };

      this.ws.onclose = () => {
        console.log('WebSocket disconnected');
        this.attemptReconnect();
      };

      this.ws.onerror = (error) => {
        console.error('WebSocket error:', error);
      };

      this.ws.onmessage = (event) => {
        try {
          const data = JSON.parse(event.data);
          this.handleMessage(data);
        } catch (error) {
          console.error('Error parsing WebSocket message:', error);
        }
      };
    } catch (error) {
      console.error('Error connecting to WebSocket:', error);
      this.attemptReconnect();
    }
  }

  /**
   * Attempts to reconnect to the WebSocket server using exponential backoff.
   * Will try up to maxReconnectAttempts times with increasing delays between attempts.
   */
  private attemptReconnect() {
    if (this.reconnectAttempts < this.maxReconnectAttempts) {
      this.reconnectAttempts++;
      const timeout = this.reconnectTimeout * Math.pow(2, this.reconnectAttempts - 1);
      console.log(`Attempting to reconnect in ${timeout}ms...`);
      setTimeout(() => this.connect(), timeout);
    } else {
      console.error('Max reconnection attempts reached');
    }
  }

  /**
   * Processes incoming WebSocket messages and triggers appropriate event handlers.
   * @param data - The parsed message data containing type and payload
   */
  private handleMessage(data: any) {
    switch (data.type) {
      case 'todo_updated':
        this.onTodoUpdated?.(data.payload);
        break;
      case 'todo_created':
        this.onTodoCreated?.(data.payload);
        break;
      case 'todo_deleted':
        this.onTodoDeleted?.(data.payload);
        break;
      default:
        console.warn('Unknown message type:', data.type);
    }
  }

  // Event handlers that can be set by the consumer
  /** Callback for when a todo item is updated */
  onTodoUpdated?: (todo: any) => void;
  /** Callback for when a new todo item is created */
  onTodoCreated?: (todo: any) => void;
  /** Callback for when a todo item is deleted */
  onTodoDeleted?: (todo: any) => void;

  /**
   * Sends data to the WebSocket server if connection is open.
   * @param data - The data to send (will be JSON stringified)
   */
  send(data: any) {
    if (this.ws?.readyState === WebSocket.OPEN) {
      this.ws.send(JSON.stringify(data));
    } else {
      console.error('WebSocket is not connected');
    }
  }

  /**
   * Closes the WebSocket connection and cleans up the instance.
   */
  disconnect() {
    if (this.ws) {
      this.ws.close();
      this.ws = null;
    }
  }
}

let wsClientInstance: WebSocketClient | null = null;

/**
 * Gets the singleton WebSocket client instance.
 * Creates a new instance if one doesn't exist.
 * @returns The WebSocket client instance or null if running server-side
 */
export const getWsClient = () => {
  if (typeof window === 'undefined') return null;
  
  if (!wsClientInstance) {
    wsClientInstance = new WebSocketClient('ws://localhost:8080/ws');
  }
  return wsClientInstance;
};

export const wsClient = getWsClient();