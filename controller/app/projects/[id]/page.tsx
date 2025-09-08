import { getProject } from '../../../lib/store'

export default function ProjectPage({ params }: { params: { id: string }}) {
  const p = getProject(params.id)
  if (!p) return <div>Not found</div>
  return (
    <div>
      <h2>{p.name}</h2>
      <pre>{JSON.stringify(p, null, 2)}</pre>
    </div>
  )
}
