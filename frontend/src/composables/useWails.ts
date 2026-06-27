// STUB: Full implementation lands in Task 25 (useWails composable).
// This stub exists so stores (Task 24) can compile against the same export shape.

export const IsInitialized: () => Promise<boolean> = async () => false
export const Verify: (pw: string) => Promise<void> = async () => {}
export const LockoutState: () => Promise<[number, string]> = async () => [0, '']
export const SetPassword: (pw: string) => Promise<void> = async () => {}
export const ClearPassword: () => Promise<void> = async () => {}
export const ListSessions: () => Promise<any[]> = async () => []
export const CreateSession: (workdir: string, prompt: string) => Promise<string> = async () => ''
export const SendMessage: (id: string, prompt: string) => Promise<void> = async () => {}
export const RespondPermission: (id: string, reqId: string, allow: boolean) => Promise<void> = async () => {}
export const CloseSession: (id: string) => Promise<void> = async () => {}
export const GetSettings: () => Promise<any> = async () => null
export const UpdateSettings: (cfg: any) => Promise<void> = async () => {}
export const GetHooksConfig: () => Promise<any> = async () => null
export const SaveHooksConfig: (cfg: any) => Promise<void> = async () => {}
export const OpenInTerminal: (workdir: string, sessionId: string, binPath: string) => Promise<void> = async () => {}
