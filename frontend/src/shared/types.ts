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

export type QuestionType = 'single_choice' | 'true_false' | 'open_ended'

export interface Question {
  id: string
  type: QuestionType
  question: string
  options: string[]
}

export interface QuizAnswerUpdate {
  playerId: string
  playerName: string
  playerColor: string
  answer: string
  status: 'pending' | 'correct' | 'incorrect'
  timestamp: number
}

export interface QuizQuestionMessage {
  id: string
  type: QuestionType
  question: string
  options: string[]
}
