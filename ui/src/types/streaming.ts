export type StreamKind = "source" | "target"

export interface StreamTarget {
  id: string
  name: string
  streamLink: string
  streamKey: string
}

export interface StreamTargetInput {
  name: string
  streamLink: string
  streamKey: string
}
