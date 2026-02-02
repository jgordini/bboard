import React, { useEffect, useState } from "react"
import { Header } from "@fider/components"
import { http } from "@fider/services"
import { HStack, VStack } from "@fider/components/layout"
import { Trans, useLingui } from "@lingui/react/macro"

interface LeaderboardPost {
  number: number
  title: string
  slug: string
  votesCount: number
  userId: number
  userName: string
}

interface LeaderboardUser {
  userId: number
  userName: string
  votesCount: number
}

const Leaderboard = () => {
  const { i18n } = useLingui()
  const [topIdeas, setTopIdeas] = useState<LeaderboardPost[]>([])
  const [topUsers, setTopUsers] = useState<LeaderboardUser[]>([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    const fetchLeaderboard = async () => {
      const [ideasRes, usersRes] = await Promise.all([
        http.get<LeaderboardPost[]>("/api/v1/leaderboard/ideas?limit=10"),
        http.get<LeaderboardUser[]>("/api/v1/leaderboard/users?limit=10"),
      ])
      if (ideasRes.ok) setTopIdeas(ideasRes.data)
      if (usersRes.ok) setTopUsers(usersRes.data)
      setLoading(false)
    }
    fetchLeaderboard()
  }, [])

  return (
    <>
      <Header />
      <div id="p-leaderboard" className="page container mt-8">
        <h1 className="text-display2 mb-6">
          <Trans id="leaderboard.title">Leaderboard</Trans>
        </h1>
        {loading ? (
          <p className="text-muted">
            <Trans id="showpost.loading">Loading...</Trans>
          </p>
        ) : (
          <HStack spacing={8} align="start" className="flex-wrap">
            <VStack spacing={4} className="flex-1 min-w-0 max-w-xl">
              <h2 className="text-xl font-semibold">
                <Trans id="leaderboard.topideas">Top Ideas</Trans>
              </h2>
              {topIdeas.length === 0 ? (
                <p className="text-muted">
                  <Trans id="leaderboard.empty.ideas">No ideas yet.</Trans>
                </p>
              ) : (
                <ul className="list-none p-0 m-0">
                  {topIdeas.map((post, i) => (
                    <li key={post.number} className="py-2 border-b border-gray-200 last:border-0">
                      <a href={`/posts/${post.number}/${post.slug}`} className="font-medium text-green-700 hover:underline">
                        #{i + 1} {post.title}
                      </a>
                      <span className="text-sm text-gray-500 ml-2">
                        {(i18n as any)._({ id: "leaderboard.votes", message: "{count} votes" }, { count: post.votesCount })} · {post.userName}
                      </span>
                    </li>
                  ))}
                </ul>
              )}
            </VStack>
            <VStack spacing={4} className="flex-1 min-w-0 max-w-xl">
              <h2 className="text-xl font-semibold">
                <Trans id="leaderboard.topusers">Top Users</Trans>
              </h2>
              <p className="text-sm text-muted">
                <Trans id="leaderboard.topusers.help">By votes received on their ideas</Trans>
              </p>
              {topUsers.length === 0 ? (
                <p className="text-muted">
                  <Trans id="leaderboard.empty.users">No users yet.</Trans>
                </p>
              ) : (
                <ul className="list-none p-0 m-0">
                  {topUsers.map((user, i) => (
                    <li key={user.userId} className="py-2 border-b border-gray-200 last:border-0">
                      <span className="font-medium">#{i + 1} {user.userName}</span>
                      <span className="text-sm text-gray-500 ml-2">
                        {(i18n as any)._({ id: "leaderboard.votes", message: "{count} votes" }, { count: user.votesCount })}
                      </span>
                    </li>
                  ))}
                </ul>
              )}
            </VStack>
          </HStack>
        )}
      </div>
    </>
  )
}

export default Leaderboard
