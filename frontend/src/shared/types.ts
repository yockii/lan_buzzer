export type GameState = 'waiting' | 'ready' | 'locked'

export interface Player {
  id: string
  name: string
  color: string
  deviceType: 'desktop' | 'mobile'
  connectedAt: string
}

export interface Message {
  type: string
  payload: any
}

export interface StateChangedPayload {
  state: GameState
  startTime?: number
}

export interface PlayerListPayload {
  players: Player[]
}

export interface BuzzResultPayload {
  winner?: Player
  isEarly: boolean
}

export interface ErrorPayload {
  message: string
}
