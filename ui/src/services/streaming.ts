import { api } from "@/utils/api"
import type { StreamTarget, StreamTargetInput } from "@/types"

// Boilerplate for future API wiring. NOT yet invoked — the UI runs on local
// mock state (see useStreamTargets). Route paths are placeholders.
export const streamingService = {
  listSources: () => api.get<StreamTarget[]>("/sources"),
  createSource: (body: StreamTargetInput) => api.post<StreamTarget>("/sources", body),
  updateSource: (id: string, body: StreamTargetInput) => api.put<StreamTarget>(`/sources/${id}`, body),
  deleteSource: (id: string) => api.delete<void>(`/sources/${id}`),

  listTargets: () => api.get<StreamTarget[]>("/targets"),
  createTarget: (body: StreamTargetInput) => api.post<StreamTarget>("/targets", body),
  updateTarget: (id: string, body: StreamTargetInput) => api.put<StreamTarget>(`/targets/${id}`, body),
  deleteTarget: (id: string) => api.delete<void>(`/targets/${id}`),
}
