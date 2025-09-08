import fs from 'node:fs'
import path from 'node:path'

export type Project = { id: string, name: string, createdAt: number, state: string, kv: Record<string,string> }
type State = { projects: Project[] }

const dataDir = path.join(process.cwd(), 'data')
const stateFile = path.join(dataDir, 'state.json')
let state: State = load()

function load(): State {
  try { return JSON.parse(fs.readFileSync(stateFile,'utf8')) } catch { return { projects: [] } }
}
function save() {
  fs.mkdirSync(dataDir, { recursive: true })
  fs.writeFileSync(stateFile, JSON.stringify(state, null, 2))
}

export function createProject(name: string): Project {
  const p: Project = { id: Math.random().toString(36).slice(2), name, createdAt: Date.now(), state: 'created', kv: {} }
  state.projects.push(p); save(); return p
}
export function listProjects() { return state.projects }
export function getProject(id: string) { return state.projects.find(p => p.id === id) }
export function setKV(id: string, kv: Record<string,string>) { const p = getProject(id); if (!p) return; p.kv = { ...p.kv, ...kv }; save() }
