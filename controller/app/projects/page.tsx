import { listProjects } from '../../lib/store'
import Link from 'next/link'

export default async function Projects() {
  const projects = listProjects()
  return (
    <div>
      <h2>Projects</h2>
      <ul>
        {projects.map(p => (
          <li key={p.id}><Link href={`/projects/${p.id}`}>{p.name}</Link></li>
        ))}
      </ul>
    </div>
  )
}
