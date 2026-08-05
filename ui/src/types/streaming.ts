export type StreamKind = "source" | "target"
export type StreamType = "rtmp"

export interface StreamTarget {
  id: string
  name: string
  type: StreamType // only "rtmp" today
  streamLink: string // maps to rtmp url
  streamKey: string
  description?: string // exists in DB; not surfaced in the form yet
  createdAt: string
}

export interface StreamTargetInput {
  name: string
  streamLink: string
  streamKey: string
}
