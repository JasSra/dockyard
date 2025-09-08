import Link from 'next/link'

export default async function Home() {
  // simple server component
  return (
    <div>
      <p>Welcome. Use the API to create projects.</p>
      <ul>
        <li><Link href="/projects">Projects</Link></li>
      </ul>
    </div>
  )
}
