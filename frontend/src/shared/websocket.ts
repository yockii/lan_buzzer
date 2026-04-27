import { ref, Ref } from 'vue'
import type { Message } from './types'

export class WebSocketClient {
  private ws: WebSocket | null = null
  private reconnectTimer: number | null = null
  private reconnectAttempts = 0
  private maxReconnectAttempts = 5
  private openHandler: (() => void) | null = null

  public connected: Ref<boolean> = ref(false)
  public messageHandlers: Map<string, (payload: any) => void> = new Map()

  constructor(private url: string) {}

  connect() {
    if (this.ws?.readyState === WebSocket.OPEN) {
      return
    }

    this.ws = new WebSocket(this.url)

    this.ws.onopen = () => {
      console.log('WebSocket connected')
      this.connected.value = true
      this.reconnectAttempts = 0
      if (this.openHandler) {
        this.openHandler()
      }
    }

    this.ws.onclose = () => {
      console.log('WebSocket disconnected')
      this.connected.value = false
      this.scheduleReconnect()
    }

    this.ws.onerror = (error) => {
      console.error('WebSocket error:', error)
    }

    this.ws.onmessage = (event) => {
      try {
        const message: Message = JSON.parse(event.data)
        const handler = this.messageHandlers.get(message.type)
        if (handler) {
          handler(message.payload)
        }
      } catch (error) {
        console.error('Failed to parse message:', error)
      }
    }
  }

  disconnect() {
    if (this.reconnectTimer) {
      clearTimeout(this.reconnectTimer)
      this.reconnectTimer = null
    }

    if (this.ws) {
      this.ws.close()
      this.ws = null
    }

    this.connected.value = false
  }

  private scheduleReconnect() {
    if (this.reconnectAttempts >= this.maxReconnectAttempts) {
      console.error('Max reconnect attempts reached')
      return
    }

    const delay = Math.min(1000 * Math.pow(2, this.reconnectAttempts), 30000)
    this.reconnectAttempts++

    this.reconnectTimer = window.setTimeout(() => {
      console.log(`Reconnecting... (attempt ${this.reconnectAttempts})`)
      this.connect()
    }, delay)
  }

  send(type: string, payload: any) {
    if (!this.ws || this.ws.readyState !== WebSocket.OPEN) {
      console.error('WebSocket not connected')
      return
    }

    const message: Message = { type, payload }
    this.ws.send(JSON.stringify(message))
  }

  on(type: string, handler: (payload: any) => void) {
    if (type === 'open') {
      this.openHandler = handler
    } else {
      this.messageHandlers.set(type, handler)
    }
  }

  off(type: string) {
    if (type === 'open') {
      this.openHandler = null
    } else {
      this.messageHandlers.delete(type)
    }
  }
}
